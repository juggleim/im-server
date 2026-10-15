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
	cmdInboxPartitionTable    = "cmdinboxmsgs_split"
	cmdSendboxPartitionTable  = "cmdsendboxmsgs_split"
	cmdPartitionInterval      = 24 * time.Hour
	cmdPartitionTick          = cmdPartitionInterval
	cmdPartitionRunTimeout    = inboxPartitionRunTimeout
	cmdPartitionLockTimeout   = inboxPartitionLockTimeout
	cmdRetainedDataCount      = 7
	cmdPreparedFutureCount    = 1
	cmdPartitionRetryInterval = 5 * time.Second
)

type cmdPartitionMaintenance struct {
	tableName string
	lockName  string
	logTag    string

	mu        sync.Mutex
	cancel    context.CancelFunc
	done      chan struct{}
	runActive int32
}

var (
	cmdInboxPartitionMaintenance = &cmdPartitionMaintenance{
		tableName: cmdInboxPartitionTable,
		lockName:  "imserver:cmdinboxmsgs-split-partitions",
		logTag:    "CMD_INBOX_PARTITION",
	}
	cmdSendboxPartitionMaintenance = &cmdPartitionMaintenance{
		tableName: cmdSendboxPartitionTable,
		lockName:  "imserver:cmdsendboxmsgs-split-partitions",
		logTag:    "CMD_SENDBOX_PARTITION",
	}
)

func StartCmdInboxPartitionMaintenance() {
	cmdInboxPartitionMaintenance.start()
}

func StopCmdInboxPartitionMaintenance() {
	cmdInboxPartitionMaintenance.stop()
}

func StartCmdSendboxPartitionMaintenance() {
	cmdSendboxPartitionMaintenance.start()
}

func StopCmdSendboxPartitionMaintenance() {
	cmdSendboxPartitionMaintenance.stop()
}

func (worker *cmdPartitionMaintenance) start() {
	if configures.Config.MsgStoreEngine != "" && configures.Config.MsgStoreEngine != configures.MsgStoreEngine_MySQL {
		return
	}

	worker.mu.Lock()
	defer worker.mu.Unlock()
	if worker.cancel != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	worker.cancel = cancel
	worker.done = make(chan struct{})
	go func(done chan struct{}) {
		defer close(done)
		retry := time.NewTicker(cmdPartitionRetryInterval)
		defer retry.Stop()
		ticker := time.NewTicker(cmdPartitionTick)
		defer ticker.Stop()

		initialized := worker.run(ctx)
		for {
			select {
			case <-ctx.Done():
				return
			case <-retry.C:
				if !initialized {
					initialized = worker.run(ctx)
				}
			case <-ticker.C:
				initialized = worker.run(ctx)
			}
		}
	}(worker.done)
}

func (worker *cmdPartitionMaintenance) stop() {
	worker.mu.Lock()
	cancel := worker.cancel
	done := worker.done
	worker.cancel = nil
	worker.done = nil
	worker.mu.Unlock()

	if cancel == nil {
		return
	}
	cancel()
	if done != nil {
		<-done
	}
}

func (worker *cmdPartitionMaintenance) run(ctx context.Context) bool {
	if !atomic.CompareAndSwapInt32(&worker.runActive, 0, 1) {
		return false
	}
	defer atomic.StoreInt32(&worker.runActive, 0)

	db := dbcommons.GetDb()
	if db == nil {
		logs.NewLogEntity().Warnf("[%s] mysql is not initialized", worker.logTag)
		return false
	}
	sqlDB, err := db.DB()
	if err != nil {
		logs.NewLogEntity().Errorf("[%s] get sql db failed err:%v", worker.logTag, err)
		return false
	}

	runCtx, cancel := context.WithTimeout(ctx, cmdPartitionRunTimeout)
	defer cancel()
	conn, err := sqlDB.Conn(runCtx)
	if err != nil {
		logs.NewLogEntity().Errorf("[%s] get connection failed err:%v", worker.logTag, err)
		return false
	}
	defer conn.Close()

	locked, err := acquireMessagePartitionLock(runCtx, conn, worker.lockName)
	if err != nil {
		logs.NewLogEntity().Errorf("[%s] acquire lock failed err:%v", worker.logTag, err)
		return false
	}
	if !locked {
		return true
	}
	defer releaseMessagePartitionLock(conn, worker.lockName, worker.logTag)

	if _, err = conn.ExecContext(runCtx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", cmdPartitionLockTimeout)); err != nil {
		logs.NewLogEntity().Errorf("[%s] set lock timeout failed err:%v", worker.logTag, err)
		return false
	}

	startedAt := time.Now()
	partitions, err := reconcileMessagePartitions(runCtx, conn, startedAt, worker.tableName)
	if err != nil {
		logs.NewLogEntity().Errorf("[%s] maintenance failed duration_ms:%d err:%v", worker.logTag, time.Since(startedAt).Milliseconds(), err)
		return false
	}
	logs.NewLogEntity().Infof("[%s] maintenance completed duration_ms:%d partitions:%d retained_data_partitions:%d", worker.logTag, time.Since(startedAt).Milliseconds(), len(partitions), cmdRetainedDataCount)
	return true
}
