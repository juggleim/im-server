package persistentqueues

import (
	"im-server/commons/tools"
)

var persistentPools *tools.SinglePools

func init() {
	persistentPools = tools.NewSinglePools(64, 500, true)
}

func DoPersistent(executor func()) {
	persistentPools.GetPool(tools.RandStr(10)).Submit(executor)
}

func DoPersistentWithTag(tag string, executor func()) {
	persistentPools.GetPool(tag).Submit(executor)
}

func NoBlockDoPersistentWithTag(tag string, executor func()) {
	persistentPools.GetPool(tag).NoBlockSubmit(executor)
}
