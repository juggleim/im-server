package services

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/services/commonservices/logs"
)

const (
	inboxPartitionTable       = "inboxmsgs_split"
	inboxPartitionFutureName  = "p_future"
	inboxPartitionLockName    = "imserver:inboxmsgs-split-partitions"
	inboxPartitionInterval    = 10 * time.Minute
	inboxPartitionTick        = 10 * time.Minute
	inboxPartitionRunTimeout  = 30 * time.Minute
	inboxPartitionLockTimeout = 5
	inboxRetainedDataCount    = 3
	inboxPreparedFutureCount  = 1
)

var (
	inboxPartitionNamePattern = regexp.MustCompile(`^p\d{12}$`)
	inboxPartitionWorkerMu    sync.Mutex
	inboxPartitionCancel      context.CancelFunc
	inboxPartitionDone        chan struct{}
	inboxPartitionRunActive   int32
)

type inboxPartition struct {
	Name        string
	Description string
	StartMillis int64
	UpperMillis int64
	Future      bool
}

type desiredInboxPartition struct {
	Name        string
	StartMillis int64
	UpperMillis int64
}

type inboxPartitionPlan struct {
	Create []desiredInboxPartition
	Drop   []string
}

type messagePartitionPolicy struct {
	TableName           string
	Interval            time.Duration
	RetainedDataCount   int
	PreparedFutureCount int
}

// StartInboxPartitionMaintenance starts one idempotent maintenance loop per
// process. Every upgraded node may attempt maintenance; a MySQL advisory lock
// ensures that only one node performs partition DDL at a time. This also works
// during a rolling deployment where the cluster hash may select an older node
// that does not have the partition-maintenance code yet.
func StartInboxPartitionMaintenance() {
	if configures.Config.MsgStoreEngine != "" && configures.Config.MsgStoreEngine != configures.MsgStoreEngine_MySQL {
		return
	}

	inboxPartitionWorkerMu.Lock()
	defer inboxPartitionWorkerMu.Unlock()
	if inboxPartitionCancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	inboxPartitionCancel = cancel
	inboxPartitionDone = make(chan struct{})
	go func(done chan struct{}) {
		defer close(done)
		retry := time.NewTicker(5 * time.Second)
		defer retry.Stop()
		ticker := time.NewTicker(inboxPartitionTick)
		defer ticker.Stop()

		initialized := runInboxPartitionMaintenance(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-retry.C:
				if !initialized {
					initialized = runInboxPartitionMaintenance(ctx)
				}
			case <-ticker.C:
				// A failed periodic run must enable the five-second retry path;
				// otherwise a missed boundary could send writes to p_future.
				initialized = runInboxPartitionMaintenance(ctx)
			}
		}
	}(inboxPartitionDone)
}

func StopInboxPartitionMaintenance() {
	inboxPartitionWorkerMu.Lock()
	cancel := inboxPartitionCancel
	done := inboxPartitionDone
	inboxPartitionCancel = nil
	inboxPartitionDone = nil
	inboxPartitionWorkerMu.Unlock()

	if cancel == nil {
		return
	}
	cancel()
	if done != nil {
		<-done
	}
}

func runInboxPartitionMaintenance(ctx context.Context) bool {
	if !atomic.CompareAndSwapInt32(&inboxPartitionRunActive, 0, 1) {
		return false
	}
	defer atomic.StoreInt32(&inboxPartitionRunActive, 0)

	db := dbcommons.GetDb()
	if db == nil {
		logs.NewLogEntity().Warnf("[INBOX_PARTITION] mysql is not initialized")
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		logs.NewLogEntity().Errorf("[INBOX_PARTITION] get sql db failed err:%v", err)
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, inboxPartitionRunTimeout)
	defer cancel()
	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		logs.NewLogEntity().Errorf("[INBOX_PARTITION] get connection failed err:%v", err)
		return false
	}
	defer conn.Close()

	locked, err := acquireInboxPartitionLock(ctx, conn)
	if err != nil {
		logs.NewLogEntity().Errorf("[INBOX_PARTITION] acquire lock failed err:%v", err)
		return false
	}
	if !locked {
		// Another upgraded node owns the maintenance lock and will reconcile the
		// same shared table. Treat this as initialized so this node waits for the
		// normal ten-minute tick instead of retrying every five seconds.
		return true
	}
	defer releaseInboxPartitionLock(conn)

	if _, err = conn.ExecContext(ctx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", inboxPartitionLockTimeout)); err != nil {
		logs.NewLogEntity().Errorf("[INBOX_PARTITION] set lock timeout failed err:%v", err)
		return false
	}

	startedAt := time.Now()
	partitions, err := reconcileInboxPartitions(ctx, conn, startedAt)
	if err != nil {
		logs.NewLogEntity().Errorf("[INBOX_PARTITION] maintenance failed duration_ms:%d err:%v", time.Since(startedAt).Milliseconds(), err)
		return false
	}
	logs.NewLogEntity().Infof("[INBOX_PARTITION] maintenance completed duration_ms:%d partitions:%d retained_data_partitions:%d", time.Since(startedAt).Milliseconds(), len(partitions), inboxRetainedDataCount)
	return true
}

func acquireInboxPartitionLock(ctx context.Context, conn *sql.Conn) (bool, error) {
	return acquireMessagePartitionLock(ctx, conn, inboxPartitionLockName)
}

func acquireMessagePartitionLock(ctx context.Context, conn *sql.Conn, lockName string) (bool, error) {
	var locked sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT GET_LOCK(?, 0)", lockName).Scan(&locked); err != nil {
		return false, err
	}
	return locked.Valid && locked.Int64 == 1, nil
}

func releaseInboxPartitionLock(conn *sql.Conn) {
	releaseMessagePartitionLock(conn, inboxPartitionLockName, "INBOX_PARTITION")
}

func releaseMessagePartitionLock(conn *sql.Conn, lockName, logTag string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var released sql.NullInt64
	if err := conn.QueryRowContext(ctx, "SELECT RELEASE_LOCK(?)", lockName).Scan(&released); err != nil {
		logs.NewLogEntity().Warnf("[%s] release lock failed err:%v", logTag, err)
	}
}

func reconcileInboxPartitions(ctx context.Context, conn *sql.Conn, now time.Time) ([]inboxPartition, error) {
	return reconcileMessagePartitions(ctx, conn, now, inboxPartitionTable)
}

func reconcileMessagePartitions(ctx context.Context, conn *sql.Conn, now time.Time, tableName string) ([]inboxPartition, error) {
	policy, err := messagePartitionPolicyForTable(tableName)
	if err != nil {
		return nil, err
	}
	partitions, err := loadMessagePartitions(ctx, conn, tableName)
	if err != nil {
		return nil, err
	}
	plan, err := planMessagePartitions(now, partitions, policy)
	if err != nil {
		return nil, err
	}
	if len(plan.Create) > 0 {
		statement, buildErr := buildCreateMessagePartitionsSQL(tableName, plan.Create)
		if buildErr != nil {
			return nil, buildErr
		}
		if _, err = conn.ExecContext(ctx, statement); err != nil {
			return nil, fmt.Errorf("create partitions: %w", err)
		}
	}
	if len(plan.Drop) > 0 {
		statement, buildErr := buildDropMessagePartitionsSQL(tableName, plan.Drop)
		if buildErr != nil {
			return nil, buildErr
		}
		if _, err = conn.ExecContext(ctx, statement); err != nil {
			return nil, fmt.Errorf("drop partition: %w", err)
		}
	}

	partitions, err = loadMessagePartitions(ctx, conn, tableName)
	if err != nil {
		return nil, err
	}
	if err = validateMessagePartitionCoverage(now, partitions, policy); err != nil {
		return nil, err
	}
	return partitions, nil
}

func loadInboxPartitions(ctx context.Context, conn *sql.Conn) ([]inboxPartition, error) {
	return loadMessagePartitions(ctx, conn, inboxPartitionTable)
}

func loadMessagePartitions(ctx context.Context, conn *sql.Conn, tableName string) ([]inboxPartition, error) {
	policy, err := messagePartitionPolicyForTable(tableName)
	if err != nil {
		return nil, err
	}
	rows, err := conn.QueryContext(ctx, `
SELECT PARTITION_NAME, PARTITION_DESCRIPTION
FROM INFORMATION_SCHEMA.PARTITIONS
WHERE TABLE_SCHEMA=DATABASE()
  AND TABLE_NAME=?
  AND PARTITION_NAME IS NOT NULL
ORDER BY PARTITION_ORDINAL_POSITION`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	partitions := make([]inboxPartition, 0)
	for rows.Next() {
		var item inboxPartition
		if err = rows.Scan(&item.Name, &item.Description); err != nil {
			return nil, err
		}
		if item.Name == inboxPartitionFutureName {
			if item.Description != "MAXVALUE" {
				return nil, fmt.Errorf("invalid p_future boundary %q", item.Description)
			}
			item.Future = true
		} else {
			start, parseErr := parseInboxPartitionName(item.Name)
			if parseErr != nil {
				return nil, parseErr
			}
			upper, parseErr := strconv.ParseInt(item.Description, 10, 64)
			if parseErr != nil {
				return nil, fmt.Errorf("partition %s has invalid boundary: %w", item.Name, parseErr)
			}
			if upper != start+policy.Interval.Milliseconds() {
				return nil, fmt.Errorf("partition %s boundary does not match its name", item.Name)
			}
			item.StartMillis = start
			item.UpperMillis = upper
		}
		partitions = append(partitions, item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(partitions) == 0 {
		return nil, fmt.Errorf("table %s does not exist or is not partitioned", tableName)
	}
	return partitions, nil
}

func planInboxPartitions(now time.Time, partitions []inboxPartition) (inboxPartitionPlan, error) {
	policy, _ := messagePartitionPolicyForTable(inboxPartitionTable)
	return planMessagePartitions(now, partitions, policy)
}

func planMessagePartitions(now time.Time, partitions []inboxPartition, policy messagePartitionPolicy) (inboxPartitionPlan, error) {
	plan := inboxPartitionPlan{}
	if policy.Interval <= 0 || policy.RetainedDataCount <= 0 || policy.PreparedFutureCount < 0 {
		return plan, fmt.Errorf("invalid partition policy for table %s", policy.TableName)
	}
	hasFuture := false
	finiteNames := make(map[string]struct{})
	lastUpper := int64(0)
	for index, partition := range partitions {
		if partition.Future {
			hasFuture = true
			if index != len(partitions)-1 {
				return plan, fmt.Errorf("p_future must be the final partition")
			}
			continue
		}
		finiteNames[partition.Name] = struct{}{}
		if partition.UpperMillis > lastUpper {
			lastUpper = partition.UpperMillis
		}
	}
	if !hasFuture {
		return plan, fmt.Errorf("p_future partition is missing")
	}

	currentStart := alignMessagePartitionStart(now.UnixMilli(), policy.Interval)
	oldestRetainedStart := currentStart - int64(policy.RetainedDataCount-1)*policy.Interval.Milliseconds()
	createStart := oldestRetainedStart
	if lastUpper > createStart {
		createStart = lastUpper
	}
	preparedThrough := currentStart + int64(policy.PreparedFutureCount)*policy.Interval.Milliseconds()
	for start := createStart; start <= preparedThrough; start += policy.Interval.Milliseconds() {
		name := formatInboxPartitionName(start)
		if _, exists := finiteNames[name]; exists {
			continue
		}
		plan.Create = append(plan.Create, desiredInboxPartition{
			Name:        name,
			StartMillis: start,
			UpperMillis: start + policy.Interval.Milliseconds(),
		})
	}

	for _, partition := range partitions {
		if !partition.Future && partition.StartMillis < oldestRetainedStart {
			plan.Drop = append(plan.Drop, partition.Name)
		}
	}
	sort.Slice(plan.Drop, func(i, j int) bool {
		left, _ := parseInboxPartitionName(plan.Drop[i])
		right, _ := parseInboxPartitionName(plan.Drop[j])
		return left < right
	})
	// A normal run drops one expired partition. If a previous policy retained
	// more partitions, dropping all excess partitions in one ALTER converges
	// immediately to the requested retention window.
	return plan, nil
}

func buildCreateInboxPartitionsSQL(partitions []desiredInboxPartition) (string, error) {
	return buildCreateMessagePartitionsSQL(inboxPartitionTable, partitions)
}

func buildCreateMessagePartitionsSQL(tableName string, partitions []desiredInboxPartition) (string, error) {
	policy, err := messagePartitionPolicyForTable(tableName)
	if err != nil {
		return "", err
	}
	if len(partitions) == 0 {
		return "", fmt.Errorf("no partitions to create")
	}
	definitions := make([]string, 0, len(partitions)+1)
	previousUpper := int64(0)
	for _, partition := range partitions {
		if !inboxPartitionNamePattern.MatchString(partition.Name) || partition.Name != formatInboxPartitionName(partition.StartMillis) {
			return "", fmt.Errorf("invalid partition name %q", partition.Name)
		}
		if partition.UpperMillis != partition.StartMillis+policy.Interval.Milliseconds() || (previousUpper > 0 && partition.UpperMillis <= previousUpper) {
			return "", fmt.Errorf("invalid partition boundary for %s", partition.Name)
		}
		definitions = append(definitions, fmt.Sprintf("PARTITION `%s` VALUES LESS THAN (%d)", partition.Name, partition.UpperMillis))
		previousUpper = partition.UpperMillis
	}
	definitions = append(definitions, "PARTITION `p_future` VALUES LESS THAN MAXVALUE")
	return fmt.Sprintf("ALTER TABLE `%s` REORGANIZE PARTITION `p_future` INTO (%s)", tableName, strings.Join(definitions, ", ")), nil
}

func buildDropInboxPartitionsSQL(names []string) (string, error) {
	return buildDropMessagePartitionsSQL(inboxPartitionTable, names)
}

func buildDropMessagePartitionsSQL(tableName string, names []string) (string, error) {
	if _, err := messagePartitionPolicyForTable(tableName); err != nil {
		return "", err
	}
	if len(names) == 0 {
		return "", fmt.Errorf("no partitions to drop")
	}
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		if name == inboxPartitionFutureName || !inboxPartitionNamePattern.MatchString(name) {
			return "", fmt.Errorf("invalid partition name %q", name)
		}
		quoted = append(quoted, "`"+name+"`")
	}
	return fmt.Sprintf("ALTER TABLE `%s` DROP PARTITION %s", tableName, strings.Join(quoted, ", ")), nil
}

func messagePartitionPolicyForTable(tableName string) (messagePartitionPolicy, error) {
	switch tableName {
	case inboxPartitionTable, sendboxPartitionTable:
		return messagePartitionPolicy{
			TableName:           tableName,
			Interval:            inboxPartitionInterval,
			RetainedDataCount:   inboxRetainedDataCount,
			PreparedFutureCount: inboxPreparedFutureCount,
		}, nil
	case cmdInboxPartitionTable, cmdSendboxPartitionTable:
		return messagePartitionPolicy{
			TableName:           tableName,
			Interval:            cmdPartitionInterval,
			RetainedDataCount:   cmdRetainedDataCount,
			PreparedFutureCount: cmdPreparedFutureCount,
		}, nil
	default:
		return messagePartitionPolicy{}, fmt.Errorf("invalid partition table %q", tableName)
	}
}

func alignInboxPartitionStart(unixMillis int64) int64 {
	return alignMessagePartitionStart(unixMillis, inboxPartitionInterval)
}

func alignMessagePartitionStart(unixMillis int64, partitionInterval time.Duration) int64 {
	interval := partitionInterval.Milliseconds()
	return unixMillis - unixMillis%interval
}

func formatInboxPartitionName(startMillis int64) string {
	return "p" + time.UnixMilli(startMillis).UTC().Format("200601021504")
}

func parseInboxPartitionName(name string) (int64, error) {
	if !inboxPartitionNamePattern.MatchString(name) {
		return 0, fmt.Errorf("invalid partition name %q", name)
	}
	parsed, err := time.ParseInLocation("200601021504", strings.TrimPrefix(name, "p"), time.UTC)
	if err != nil {
		return 0, fmt.Errorf("invalid partition name %q: %w", name, err)
	}
	return parsed.UnixMilli(), nil
}

func validateInboxPartitionCoverage(now time.Time, partitions []inboxPartition) error {
	policy, _ := messagePartitionPolicyForTable(inboxPartitionTable)
	return validateMessagePartitionCoverage(now, partitions, policy)
}

func validateMessagePartitionCoverage(now time.Time, partitions []inboxPartition, policy messagePartitionPolicy) error {
	currentStart := alignMessagePartitionStart(now.UnixMilli(), policy.Interval)
	currentName := formatInboxPartitionName(currentStart)
	nextName := formatInboxPartitionName(currentStart + policy.Interval.Milliseconds())
	hasCurrent := false
	hasNext := false
	hasFuture := false
	for _, partition := range partitions {
		switch partition.Name {
		case currentName:
			hasCurrent = true
		case nextName:
			hasNext = true
		case inboxPartitionFutureName:
			hasFuture = true
		}
	}
	if !hasCurrent || !hasNext || !hasFuture {
		return fmt.Errorf("partition coverage missing current:%t next:%t future:%t", hasCurrent, hasNext, hasFuture)
	}
	return nil
}
