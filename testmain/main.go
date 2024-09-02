package main

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/commons/logs"
	"im-server/commons/mongocommons"
	hisMsgMongo "im-server/services/historymsg/storages/mongodbs"
	msgMongo "im-server/services/message/storages/mongodbs"
)

func main() {
	//init configures
	if err := configures.InitConfigures(); err != nil {
		//logs.Error("Init Configures failed.", err)
		fmt.Println("Init Configures failed.", err)
		return
	}
	// init logs
	logs.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		logs.Error("Init Mysql failed.", err)
		return
	}
	//init mongodb
	if configures.Config.MsgStoreEngine == configures.MsgStoreEngine_Mongo {
		if err := mongocommons.InitMongodb(); err != nil {
			logs.Error("Init MongoDB failed.", err)
			return
		} else {
			hisMsgMongo.RegistCollections()
			msgMongo.RegistCollections()
			mongocommons.InitMongoCollections()
		}
	}
}
