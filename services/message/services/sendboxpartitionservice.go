package services

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/services/commonservices/logs"
)

const (
	sendboxPartitionTable       = "sendboxmsgs_split"
	sendboxPartitionLockName    = "imserver:sendboxmsgs-split-partitions"
	sendboxPartitionTick        = inboxPartitionTick
	sendboxPartitionRunTimeout  = inboxPartitionRunTimeout
	sendboxPartitionLockTimeout = inboxPartitionLockTimeout
)

var (
	sendboxPartitionWorkerMu  sync.Mutex
	sendboxPartitionCancel    context.CancelFunc
	sendboxPartitionDone      chan struct{}
	sendboxPartitionRunActive int32
)

// StartSendboxPartitionMaintenance starts one idempotent maintenance loop per
// process. All upgraded nodes may attempt maintenance; the table-specific MySQL
// advisory lock allows only one node to execute sendbox partition DDL at once.
func StartSendboxPartitionMaintenance() {
	if configures.Config.MsgStoreEngine != "" && configures.Config.MsgStoreEngine != configures.MsgStoreEngine_MySQL {
		return
	}

	sendboxPartitionWorkerMu.Lock()
	defer sendboxPartitionWorkerMu.Unlock()
	if sendboxPartitionCancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sendboxPartitionCancel = cancel
	sendboxPartitionDone = make(chan struct{})
	go func(done chan struct{}) {
		defer close(done)
		retry := time.NewTicker(5 * time.Second)
		defer retry.Stop()
		ticker := time.NewTicker(sendboxPartitionTick)
		defer ticker.Stop()

		initialized := runSendboxPartitionMaintenance(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-retry.C:
				if !initialized {
					initialized = runSendboxPartitionMaintenance(ctx)
				}
			case <-ticker.C:
				initialized = runSendboxPartitionMaintenance(ctx)
			}
		}
	}(sendboxPartitionDone)
}

func StopSendboxPartitionMaintenance() {
	sendboxPartitionWorkerMu.Lock()
	cancel := sendboxPartitionCancel
	done := sendboxPartitionDone
	sendboxPartitionCancel = nil
	sendboxPartitionDone = nil
	sendboxPartitionWorkerMu.Unlock()

	if cancel == nil {
		return
	}
	cancel()
	if done != nil {
		<-done
	}
}

func runSendboxPartitionMaintenance(ctx context.Context) bool {
	if !atomic.CompareAndSwapInt32(&sendboxPartitionRunActive, 0, 1) {
		return false
	}
	defer atomic.StoreInt32(&sendboxPartitionRunActive, 0)

	db := dbcommons.GetDb()
	if db == nil {
		logs.NewLogEntity().Warnf("[SENDBOX_PARTITION] mysql is not initialized")
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		logs.NewLogEntity().Errorf("[SENDBOX_PARTITION] get sql db failed err:%v", err)
		return false
	}

	runCtx, cancel := context.WithTimeout(ctx, sendboxPartitionRunTimeout)
	defer cancel()
	conn, err := sqlDB.Conn(runCtx)
	if err != nil {
		logs.NewLogEntity().Errorf("[SENDBOX_PARTITION] get connection failed err:%v", err)
		return false
	}
	defer conn.Close()

	locked, err := acquireMessagePartitionLock(runCtx, conn, sendboxPartitionLockName)
	if err != nil {
		logs.NewLogEntity().Errorf("[SENDBOX_PARTITION] acquire lock failed err:%v", err)
		return false
	}
	if !locked {
		return true
	}
	defer releaseMessagePartitionLock(conn, sendboxPartitionLockName, "SENDBOX_PARTITION")

	if _, err = conn.ExecContext(runCtx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", sendboxPartitionLockTimeout)); err != nil {
		logs.NewLogEntity().Errorf("[SENDBOX_PARTITION] set lock timeout failed err:%v", err)
		return false
	}

	startedAt := time.Now()
	partitions, err := reconcileMessagePartitions(runCtx, conn, startedAt, sendboxPartitionTable)
	if err != nil {
		logs.NewLogEntity().Errorf("[SENDBOX_PARTITION] maintenance failed duration_ms:%d err:%v", time.Since(startedAt).Milliseconds(), err)
		return false
	}
	logs.NewLogEntity().Infof("[SENDBOX_PARTITION] maintenance completed duration_ms:%d partitions:%d retained_data_partitions:%d", time.Since(startedAt).Milliseconds(), len(partitions), inboxRetainedDataCount)
	return true
}
