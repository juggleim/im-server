package services

import (
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const maxLogCount = 10000

var (
	db *leveldb.DB
)

type LogEntity struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func InitLogDB(path string) (err error) {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	db, err = leveldb.OpenFile(path, o)
	if err != nil {
		err = errors.Wrap(err, "Failed to open log db")
		return

	}
	return nil
}

func CloseLogDB() {
	err := db.Close()
	if err != nil {
		return
	}
}

func writeLog(table string, key string, value string) (err error) {
	batch := new(leveldb.Batch)

	batch.Put([]byte(table+":"+key), []byte(value))

	return db.Write(batch, nil)
}

func qryLogs(table string, keyPrefix, startKey string, count int) ([]LogEntity, error) {
	iter := db.NewIterator(util.BytesPrefix([]byte(table+":"+keyPrefix)), nil)
	iterCount := 0
	logs := []LogEntity{}
	if startKey == "" {
		iter.First()
	} else {
		iter.Seek([]byte(table + ":" + startKey))
	}
	for ; iter.Valid(); iter.Next() {
		logs = append(logs, LogEntity{
			Key:   string(iter.Key()),
			Value: string(iter.Value()),
		})
		iterCount++
		if iterCount > maxLogCount {
			break
		}
		if iterCount >= count {
			break
		}
	}
	return logs, nil
}
