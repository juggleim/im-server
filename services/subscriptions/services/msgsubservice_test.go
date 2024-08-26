package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/configures"
	"im-server/commons/dbcommons"
	"im-server/commons/logs"
	"im-server/commons/pbdefines/pbobjs"
	"testing"
	"time"
)

func TestMsgSubHandle(t *testing.T) {
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

	ctx := context.WithValue(context.Background(), bases.CtxKey_AppKey, "appkey")
	MsgSubHandle(ctx, &pbobjs.DownMsgSet{
		Msgs: []*pbobjs.DownMsg{
			{
				TargetId:       "1",
				ChannelType:    pbobjs.ChannelType_Private,
				MsgType:        "text",
				SenderId:       "1",
				MsgId:          "111",
				MsgSeqNo:       0,
				MsgContent:     []byte("hello, this is a test message"),
				MsgTime:        time.Now().UnixMilli(),
				Flags:          0,
				IsSend:         false,
				Platform:       "",
				ClientUid:      "",
				PushData:       nil,
				MentionInfo:    nil,
				IsRead:         false,
				ReferMsg:       nil,
				TargetUserInfo: nil,
				GroupInfo:      nil,
				MergedMsgs:     nil,
				UndisturbType:  0,
				MemberCount:    0,
				ReadCount:      0,
				UnreadIndex:    0,
			},
		},
	})
}
