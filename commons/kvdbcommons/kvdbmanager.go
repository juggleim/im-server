package kvdbcommons

import (
	"fmt"
	"im-server/commons/configures"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
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
