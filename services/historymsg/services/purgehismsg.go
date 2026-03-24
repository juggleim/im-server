package services

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/tasks"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages/dbs"
)

func PurgePrivateHisMsg(appkey string, msgTime int64) {
	taskKey := "purge_private_msg"
	if canPurgeMsg(taskKey) {
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil && appinfo.HistoryMsgSaveDay > 0 {
			expiredTime := msgTime - int64(appinfo.HistoryMsgSaveDay)*24*60*60*1000
			tasks.TaskExecute(taskKey, 60*1000, func() {
				//clean private msg
				priDao := dbs.PrivateHisMsgDao{}
				err := priDao.DelMsgsBaseTime(appkey, expiredTime)
				if err != nil {
					fmt.Println("Del expired private msgs err:", err)
				}
				//clean private del msg
				priDelDao := dbs.PrivateDelHisMsgDao{}
				err = priDelDao.DelMsgsBaseTime(appkey, expiredTime)
				if err != nil {
					fmt.Println("Del expired private del msgs err:", err)
				}
			})
		}
	}
}

func PurgeGroupHisMsg(appkey string, msgTime int64) {
	taskKey := "purge_group_msg"
	if canPurgeMsg(taskKey) {
		appinfo, exist := commonservices.GetAppInfo(appkey)
		if exist && appinfo != nil && appinfo.HistoryMsgSaveDay > 0 {
			expiredTime := msgTime - int64(appinfo.HistoryMsgSaveDay)*24*60*60*1000
			tasks.TaskExecute(taskKey, 60*1000, func() {
				//clean group msg
				grpDao := dbs.GroupHisMsgDao{}
				err := grpDao.DelMsgsBaseTime(appkey, expiredTime)
				if err != nil {
					fmt.Println("Del expired group msgs err:", err)
				}
				//clean group del msg
				grpDelDao := dbs.GroupDelHisMsgDao{}
				err = grpDelDao.DelMsgsBaseTime(appkey, expiredTime)
				if err != nil {
					fmt.Println("Del expired group del msgs err:", err)
				}
				//clean group portion msg
				grpPorDao := dbs.GroupPortionRelDao{}
				grpPorDao.DelRelsBaseTime(appkey, expiredTime)
				if err != nil {
					fmt.Println("Del expired group portion msgs err:", err)
				}
			})
		}
	}
}

func canPurgeMsg(taskKey string) bool {
	if configures.Config.MsgStoreEngine == "" || configures.Config.MsgStoreEngine == configures.MsgStoreEngine_MySQL {
		nodeName := bases.GetCluster().GetCurrentNode().Name
		node := bases.GetCluster().GetTargetNode("add_hismsg", taskKey)
		return node.Name == nodeName
	}
	return false
}
