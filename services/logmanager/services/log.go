package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
	"strconv"
	"strings"
)

const maxLogCount = 10000

var (
	db        *leveldb.DB
	metaCount = NewMetaCount()
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

func tableRowsCount(appKey, database string, table string) (count int64) {
	tableRowsCountKey := tableMetaName(appKey, database, table) + "_count"
	v, _ := db.Get([]byte(tableRowsCountKey), nil)
	if v != nil {
		count = bytesToInt64(v)
	}
	return
}

func dbRowsCount(appKey, database string) (count int64) {
	dbrRowsCountKey := dbMetaName(appKey, database) + "_count"
	v, _ := db.Get([]byte(dbrRowsCountKey), nil)
	if v != nil {
		count = bytesToInt64(v)
	}
	return
}

func writeLog(appKey, database string, table string, key string, value string) (err error) {
	batch := new(leveldb.Batch)
	tableName := tableName(appKey, database, table)

	batch.Put([]byte(tableName+":"+key), []byte(value))

	tableRowsCountKey := tableMetaName(appKey, database, table) + ":count"
	tableRowsCount := metaCount.Increment(tableRowsCountKey, func(key string) int64 {
		v, _ := db.Get([]byte(key), nil)
		return bytesToInt64(v)
	})

	batch.Put([]byte(tableRowsCountKey), int64ToBytes(tableRowsCount))

	databaseRowsCountKey := dbMetaName(appKey, database) + ":count"
	databaseRowsCount := metaCount.Increment(databaseRowsCountKey, func(key string) int64 {
		v, _ := db.Get([]byte(key), nil)
		return bytesToInt64(v)
	})
	batch.Put([]byte(databaseRowsCountKey), int64ToBytes(databaseRowsCount))

	return db.Write(batch, nil)
}

func fetchLogs(appKey, database string, table string, prefix string, lastKey string, count int) (logs []LogEntity, err error) {
	tableName := tableName(appKey, database, table)
	lastTimestamp := keyTimestamp(lastKey)

	iter := db.NewIterator(util.BytesPrefix([]byte(tableName+":"+prefix)), nil)
	if lastKey == "" {
		var iterCount int
		for iter.Last(); iter.Valid(); iter.Prev() {
			iterCount++
			if iterCount > maxLogCount {
				break
			}
			key := iter.Key()
			value := iter.Value()

			logs = append(logs, LogEntity{
				Key:   string(key),
				Value: string(value),
			})
			if len(logs) >= maxLogCount || len(logs) >= count {
				break
			}
		}
	} else {
		var iterCount int
		if ok := iter.Seek([]byte(tableName + ":" + lastKey)); !ok {
			iter.Last()
		}

		for ok := iter.Valid(); ok; ok = iter.Prev() {
			iterCount++
			if iterCount > maxLogCount {
				break
			}

			key := iter.Key()
			value := iter.Value()

			timestamp := keyTimestamp(string(key))
			if timestamp < lastTimestamp {
				logs = append(logs, LogEntity{
					Key:   string(key),
					Value: string(value),
				})
				if len(logs) >= maxLogCount || len(logs) >= count {
					break
				}
			}
		}
	}

	iter.Release()
	if err := iter.Error(); err != nil {
		err = errors.Wrap(err, "Failed to fetch connection logs")
	}

	return
}

func keyTimestamp(key string) (timestamp int64) {
	strList := strings.Split(key, "_")
	if len(strList) > 0 {
		timeStr := strList[len(strList)-1]
		timestamp, _ = strconv.ParseInt(timeStr, 10, 64)
		return
	}

	return 0
}

func int64ToBytes(i int64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, i)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	return buf.Bytes()
}

func bytesToInt64(bs []byte) int64 {
	buf := bytes.NewReader(bs)
	var x int64
	_ = binary.Read(buf, binary.LittleEndian, &x)
	return x
}

type FetchLogOptions struct {
	AppKey    string
	Table     string
	UserId    string
	StartTime int64
	Session   string
	LastKey   string
	Count     int
	Prev      bool
}

func FetchAppServerLogs(options FetchLogOptions) (logs []LogEntity, err error) {
	lastTimestamp := keyTimestamp(options.LastKey)
	var upIter func(iterator.Iterator) bool
	if options.Prev {
		upIter = func(iter iterator.Iterator) bool {
			return iter.Prev()
		}
	} else {
		upIter = func(iter iterator.Iterator) bool {
			return iter.Next()
		}
	}

	prefix := options.AppKey + ".server." + options.Table
	if options.Table == connectionTable {
		prefix = prefix + ":" + options.UserId
	} else if options.Table == sessionTable {
		prefix = prefix + ":" + options.Session
	} else if options.Table == businessTable {
		prefix = prefix + ":" + options.Session
	}

	iter := db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	if options.LastKey == "" {
		var iterCount int
		for iter.Last(); iter.Valid(); upIter(iter) {
			iterCount++
			if iterCount > maxLogCount {
				break
			}
			key := iter.Key()
			value := iter.Value()
			if options.StartTime != 0 {
				timestamp := keyTimestamp(string(key))
				if timestamp < options.StartTime {
					continue
				}
			}

			logs = append(logs, LogEntity{
				Key:   string(key),
				Value: string(value),
			})
			if len(logs) >= maxLogCount || len(logs) >= options.Count {
				break
			}
		}
	} else {
		var iterCount int
		if ok := iter.Seek([]byte(options.LastKey)); !ok {
			iter.Last()
		}

		for ok := iter.Valid(); ok; ok = upIter(iter) {
			iterCount++
			if iterCount > maxLogCount {
				break
			}

			key := iter.Key()
			value := iter.Value()

			timestamp := keyTimestamp(string(key))
			if options.StartTime != 0 {
				if timestamp < options.StartTime {
					continue
				}
			}
			if (options.Prev && timestamp < lastTimestamp) || (!options.Prev && timestamp > lastTimestamp) {
				logs = append(logs, LogEntity{
					Key:   string(key),
					Value: string(value),
				})
				if len(logs) >= maxLogCount || len(logs) >= options.Count {
					break
				}
			}
		}
	}

	iter.Release()
	if err := iter.Error(); err != nil {
		err = errors.Wrap(err, "Failed to fetch connection logs")
	}
	if !options.Prev {
		// reverse the order of logs
		for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
			logs[i], logs[j] = logs[j], logs[i]
		}
	}

	return
}
