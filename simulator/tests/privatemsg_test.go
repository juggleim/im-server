package tests

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices/msgdefines"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPrivateMsg(t *testing.T) {
	wsClient1 := wsclients.NewWsImClient(WsAddr, Appkey, Token1, nil, nil, nil)
	code, connectAck := wsClient1.Connect("", "")
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.Equal(t, User1, connectAck.UserId)
		Print("%s is connected.", User1)
	}

	msgContainer := map[string]*pbobjs.DownMsg{}
	wsClient2 := wsclients.NewWsImClient(WsAddr, Appkey, Token2, func(msg *pbobjs.DownMsg) {
		msgContainer[msg.MsgId] = msg
	}, nil, nil)
	code, connectAck = wsClient2.Connect("", "")
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.Equal(t, User2, connectAck.UserId)
		Print("%s is connected.", User2)
	}

	//send private msg
	flag := msgdefines.SetStoreMsg(0)
	flag = msgdefines.SetCountMsg(flag)
	code, pubAck := wsClient1.SendPrivateMsg(User2, &pbobjs.UpMsg{
		MsgType:    "t:txtmsg",
		MsgContent: []byte("{\"content\":\"hello im.\"}"),
		Flags:      flag,
	})
	var msgId string = ""
	var msgTime int64 = 0
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		msgId = pubAck.MsgId
		msgTime = pubAck.Timestamp
	}
	assert.NotEmpty(t, msgId)
	assert.NotEqual(t, 0, msgTime)
	Print("%s send msg to %s. [msg_id:%s,msg_time:%d]", User1, User2, msgId, msgTime)

	//check user2 received msg
	time.Sleep(1 * time.Second)
	receiveMsg, isReceived := msgContainer[msgId]
	assert.Equal(t, true, isReceived)
	assert.Equal(t, msgId, receiveMsg.MsgId)
	Print("%s online receive msg.[msg_id:%s,msg_time:%d]", User2, receiveMsg.MsgId, receiveMsg.MsgTime)

	//get offline msg from sendbox
	code, msgs := wsClient1.SyncMsgs(&pbobjs.SyncMsgReq{
		SyncTime:        msgTime - 1,
		ContainsSendBox: true,
		SendBoxSyncTime: msgTime - 1,
	})
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, msgs)
		assert.NotEqual(t, 0, len(msgs.Msgs))
		Print("%s success to sync msg from sendbox. msg_len:%d", User1, len(msgs.Msgs))
		msg := msgs.Msgs[0]
		assert.Equal(t, msgId, msg.MsgId)
		Print("%s sync msg from sendbox. msg_id:%s", User1, msg.MsgId)
	}

	//get offline msg from inbox
	code, msgs = wsClient2.SyncMsgs(&pbobjs.SyncMsgReq{
		SyncTime:        msgTime - 1,
		ContainsSendBox: false,
		SendBoxSyncTime: msgTime - 1,
	})
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, msgs)
		assert.NotEqual(t, 0, len(msgs.Msgs))
		Print("%s success to sync msg from inbox. msg_len:%d", User2, len(msgs.Msgs))
		msg := msgs.Msgs[0]
		Print("%s sync msg from inbox. msg_id:%s", User2, msg.MsgId)
		assert.Equal(t, msgId, msg.MsgId)
	}

	//sender get msg from history msg
	code, hisMsgs := wsClient1.QryHistoryMsgs(&pbobjs.QryHisMsgsReq{
		TargetId:    User2,
		ChannelType: pbobjs.ChannelType_Private,
		StartTime:   msgTime + 1,
		Count:       1,
	})
	Print("%s query history msg. code:%d", User1, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, hisMsgs)
		assert.NotEqual(t, 0, len(hisMsgs.Msgs))
		msg := hisMsgs.Msgs[0]
		Print("%s qry history msg. msg_id:%s, is_send:%v", User1, msg.MsgId, msg.IsSend)
		assert.Equal(t, msgId, msg.MsgId)
		assert.Equal(t, true, msg.IsSend)
	}
	//receiver get msg from history msg
	code, hisMsgs = wsClient2.QryHistoryMsgs(&pbobjs.QryHisMsgsReq{
		TargetId:    User1,
		ChannelType: pbobjs.ChannelType_Private,
		StartTime:   msgTime + 1,
		Count:       1,
	})
	Print("%s query history msg. code:%d", User2, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, hisMsgs)
		assert.NotEqual(t, 0, len(hisMsgs.Msgs))
		msg := hisMsgs.Msgs[0]
		Print("%s qry history msg. msg_id:%s, is_send:%v", User2, msg.MsgId, msg.IsSend)
		assert.Equal(t, msgId, msg.MsgId)
		assert.Equal(t, false, msg.IsSend)
	}

	wsClient1.Disconnect()
	wsClient2.Disconnect()
}
