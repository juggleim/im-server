package kvdbcommons

import (
	"bytes"
	"fmt"
	"im-server/commons/configures"

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

func Put(key, value []byte) error {
	if kvdb == nil {
		return errors.New("kv db not init")
	}
	return kvdb.Put(key, value, nil)
}

func BatchPut(kvPairs []KeyValPair) error {
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
