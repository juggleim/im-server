package dbs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"im-server/commons/dbcommons"
	"strings"
	"time"
)

const (
	msgPurgeBatchSize   = 1000
	msgPurgeMaxRows     = 10000
	msgPurgeMaxDuration = 5 * time.Second
	msgPurgeBatchPause  = 20 * time.Millisecond
)

// delMsgsBaseTime deletes a bounded amount of data. Reaching the row or time
// budget is a successful partial purge; a later scheduler run continues it.
func delMsgsBaseTime(tableName, appkey string, startTime int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), msgPurgeMaxDuration)
	defer cancel()
	db, err := dbcommons.GetDb().DB()
	if err != nil {
		return err
	}

	deleted := 0
	for deleted < msgPurgeMaxRows {
		if ctx.Err() != nil {
			return nil
		}

		batchSize := msgPurgeBatchSize
		if remaining := msgPurgeMaxRows - deleted; remaining < batchSize {
			batchSize = remaining
		}

		var ids []int64
		rows, err := db.QueryContext(ctx, fmt.Sprintf(
			"SELECT id FROM `%s` WHERE app_key=? AND send_time<? ORDER BY send_time ASC, id ASC LIMIT ?",
			tableName,
		), appkey, startTime, batchSize)
		if err != nil {
			if isPurgeDeadline(err) {
				return nil
			}
			return err
		}
		ids, err = scanPurgeIDs(rows)
		if err != nil {
			if isPurgeDeadline(err) {
				return nil
			}
			return err
		}
		if len(ids) == 0 {
			return nil
		}

		args := make([]interface{}, len(ids))
		for i, id := range ids {
			args[i] = id
		}
		_, err = db.ExecContext(ctx, fmt.Sprintf(
			"DELETE FROM `%s` WHERE id IN (%s)",
			tableName,
			strings.TrimSuffix(strings.Repeat("?,", len(ids)), ","),
		), args...)
		if err != nil {
			if isPurgeDeadline(err) {
				return nil
			}
			return err
		}
		deleted += len(ids)
		if len(ids) < batchSize || deleted >= msgPurgeMaxRows {
			return nil
		}

		timer := time.NewTimer(msgPurgeBatchPause)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}
	return nil
}

func scanPurgeIDs(rows *sql.Rows) ([]int64, error) {
	defer rows.Close()
	ids := make([]int64, 0, msgPurgeBatchSize)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func isPurgeDeadline(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)
}
