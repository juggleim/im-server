package main

import (
	"fmt"
	"im-server/services/botmsg"
	"im-server/services/broadcast"
	"im-server/services/logmanager"
	"im-server/services/rtcroom"
	sensitivemanager "im-server/services/sensitivemanager"
	"im-server/services/subscriptions"
	"im-server/services/userstatussub"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/commons/imstarters"
	"im-server/commons/kvdbcommons"
	"im-server/commons/logs"
	"im-server/commons/mongocommons"
	"im-server/commons/tools"
	"im-server/services/admingateway"
	"im-server/services/apigateway"
	"im-server/services/connectmanager"
	"im-server/services/conversation"
	"im-server/services/fileplugin"
	"im-server/services/group"
	"im-server/services/historymsg"
	"im-server/services/message"
	push "im-server/services/pushmanager"
	"im-server/services/usermanager"

	hisMsgMongo "im-server/services/historymsg/storages/mongodbs"
	msgMongo "im-server/services/message/storages/mongodbs"
	pushMongo "im-server/services/pushmanager/storages/mongodbs"
)

func main() {
	start := time.Now()
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()

	//init configures
	if err := configures.InitConfigures(); err != nil {
		//logs.Error("Init Configures failed.", err)
		fmt.Println("Init Configures failed.", err)
		return
	}
	//init logs
	logs.InitLogs()
	//init mysql
	if err := dbcommons.InitMysql(); err != nil {
		logs.Error("Init Mysql failed.", err)
		return
	}
	//upgrade mysql
	dbcommons.Upgrade()
	//init tsdb
	if err := kvdbcommons.InitKvdb(); err != nil {
		logs.Error("Init KvDB failed.", err)
	}
	//init mongodb
	if configures.Config.MsgStoreEngine == configures.MsgStoreEngine_Mongo {
		if err := mongocommons.InitMongodb(); err != nil {
			logs.Error("Init MongoDB failed.", err)
			return
		} else {
			hisMsgMongo.RegistCollections()
			msgMongo.RegistCollections()
			pushMongo.RegistCollections()
			mongocommons.InitMongoCollections()
		}
	}

	//init cluster
	exts := map[string]string{}
	exts[bases.NodeTag_Nav] = tools.ToJson(bases.HttpNodeExt{
		Port: configures.Config.NavGateway.HttpPort,
	})
	exts[bases.NodeTag_Api] = tools.ToJson(bases.HttpNodeExt{
		Port: configures.Config.ApiGateway.HttpPort,
	})
	exts[bases.NodeTag_Connect] = tools.ToJson(bases.ConnectNodeExt{
		WsPort: configures.Config.ConnectManager.WsPort,
	})
	exts[bases.NodeTag_Admin] = tools.ToJson(bases.HttpNodeExt{
		Port: configures.Config.ConnectManager.WsPort,
	})
	if err := bases.InitImServer(exts); err != nil {
		logs.Error("Init Cluster failed.", err)
		return
	}

	imstarters.Loaded(&admingateway.AdminGateway{})
	imstarters.Loaded(&apigateway.ApiGateway{})
	imstarters.Loaded(&connectmanager.ConnectManager{})
	imstarters.Loaded(&message.MessageManager{})
	imstarters.Loaded(&conversation.ConversationManager{})
	imstarters.Loaded(&usermanager.UserManager{})
	imstarters.Loaded(&historymsg.HistoryMsgManager{})
	imstarters.Loaded(&group.GroupManager{})
	imstarters.Loaded(&push.PushManager{})
	imstarters.Loaded(&fileplugin.FilePlugin{})
	imstarters.Loaded(&subscriptions.SubscriptionManager{})
	imstarters.Loaded(&broadcast.BroadcastManager{})
	imstarters.Loaded(&logmanager.LogManager{})
	imstarters.Loaded(&sensitivemanager.SensitiveManager{})
	imstarters.Loaded(&userstatussub.UserStatusSubManager{})
	imstarters.Loaded(&botmsg.BotMsgManager{})
	imstarters.Loaded(&rtcroom.RtcRoomManager{})

	imstarters.Startup()
	fmt.Println("expand:", time.Since(start))

	closeChan := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigChan
		imstarters.Shutdown(true)
		signal.Stop(sigChan)
		close(closeChan)
	}()

	<-closeChan
}
