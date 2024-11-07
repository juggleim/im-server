package kvdbcommons

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/tools"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var kvdb *leveldb.DB

func InitKvdb() (err error) {
	if configures.Config.Kvdb.IsOpen {
		path := configures.Config.Kvdb.DataPath
		if path == "" {
			path = fmt.Sprintf("%s/tsdb_data", configures.Config.Log.LogPath)
		}
		o := &opt.Options{
			Filter:              filter.NewBloomFilter(10),
			CompactionTableSize: 10 * opt.MiB,
		}
		kvdb, err = leveldb.OpenFile(path, o)
		if err != nil {
			err = errors.Wrap(err, "failed to open ts db")
			return
		}
		startEvictTask()
	}
	return nil
}

func CloseKvdb() {
	if kvdb != nil {
		err := kvdb.Close()
		if err != nil {
			fmt.Println("failed to close ts db")
		}
	}
}

var kvdbLocks *tools.SegmentatedLocks

func init() {
	kvdbLocks = tools.NewSegmentatedLocks(256)
}

func Set(key, value []byte) error {
	if kvdb == nil {
		return errors.New("kv db not init")
	}
	return kvdb.Put(key, value, nil)
}

func BatchSet(kvPairs []KeyValPair) error {
	if kvdb == nil {
		return errors.New("kv db not init")
	}
	if len(kvPairs) <= 0 {
		return nil
	}
	batch := new(leveldb.Batch)
	for _, kv := range kvPairs {
		batch.Put(kv.Key, kv.Val)
	}
	return kvdb.Write(batch, nil)
}

func Get(key []byte) ([]byte, error) {
	if kvdb == nil {
		return []byte{}, errors.New("kv db not init")
	}
	return kvdb.Get(key, nil)
}

func Exist(key []byte) (bool, error) {
	if kvdb == nil {
		return false, errors.New("kv db not init")
	}
	return kvdb.Has(key, nil)
}

func SetNx(key, value []byte) (bool, error) {
	lock := kvdbLocks.GetLocks(base64.URLEncoding.EncodeToString(key))
	lock.Lock()
	defer lock.Unlock()

	isExist, err := Exist(key)
	if err != nil {
		return false, err
	}
	if isExist {
		return false, nil
	}
	return true, Set(key, value)
}

func SetNxWithIncrByStep(key []byte, step int64) (bool, int64, error) {
	lock := kvdbLocks.GetLocks(base64.URLEncoding.EncodeToString(key))
	lock.Lock()
	defer lock.Unlock()

	isExist, err := Exist(key)
	if err != nil {
		return false, 0, err
	}
	if isExist {
		val, err := Get(key)
		if err != nil {
			return false, 0, err
		}
		intVal := tools.BytesToInt64(val)
		intVal = intVal + step
		return false, intVal, Set(key, tools.Int64ToBytes(intVal))
	}
	return true, step, Set(key, tools.Int64ToBytes(step))
}

func SetNxWithIncr(key []byte) (bool, int64, error) {
	return SetNxWithIncrByStep(key, 1)
}

func Delete(key []byte) error {
	if kvdb == nil {
		return errors.New("kv db not init")
	}
	return kvdb.Delete(key, nil)
}

type KeyValPair struct {
	Key []byte
	Val []byte
}

func Scan(prefix, start []byte, count int) ([]KeyValPair, error) {
	if kvdb == nil {
		return []KeyValPair{}, errors.New("kv db not init")
	}
	iterCount := 0
	ret := []KeyValPair{}
	var skipKeys []byte
	iter := kvdb.NewIterator(util.BytesPrefix(prefix), nil)
	if len(start) <= 0 {
		iter.First()
	} else {
		begin := append(prefix, start...)
		iter.Seek(begin)
		skipKeys = begin
	}
	for ; iter.Valid(); iter.Next() {
		keyBs := []byte{}
		keyBs = append(keyBs, iter.Key()...)
		valBs := []byte{}
		valBs = append(valBs, iter.Value()...)
		if len(skipKeys) > 0 && bytes.Equal(keyBs, skipKeys) {
			continue
		}
		ret = append(ret, KeyValPair{
			Key: keyBs,
			Val: valBs,
		})
		iterCount++
		if iterCount >= count {
			break
		}
	}
	return ret, nil
}
