package tests

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
	"testing"
	"time"

	// serversdk "github.com/juggleim/imserver-sdk-go"

	"github.com/stretchr/testify/assert"
)

func TestGroupMsg(t *testing.T) {
	// sdk := serversdk.NewJuggleIMSdk(Appkey, AppSecret, ApiURL)

	// //create groups
	// req := serversdk.GroupMembersReq{
	// 	GroupId:   Group1,
	// 	MemberIds: []string{User1, User2, User3},
	// }
	// apiCode, _, err := sdk.CreateGroup(req)
	// Print("create group[%s],members[%v] code:[%d] err:[%v]", req.GroupId, req.MemberIds, apiCode, err)
	// assert.Equal(t, nil, err)
	// assert.Equal(t, serversdk.ApiCode_Success, apiCode)

	//qry groupmembers

	//user1 connect
	wsClient1 := wsclients.NewWsImClient(WsAddr, Appkey, Token1, nil, nil, nil)
	code, connectAck := wsClient1.Connect("", "")
	Print("%s connecting... code[%d]", User1, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.Equal(t, User1, connectAck.UserId)
		Print("%s is connected.", User1)
	}
	//user2 connect
	msgContainer2 := map[string]*pbobjs.DownMsg{}
	wsClient2 := wsclients.NewWsImClient(WsAddr, Appkey, Token2, func(msg *pbobjs.DownMsg) {
		msgContainer2[msg.MsgId] = msg
	}, nil, nil)
	code, connectAck = wsClient2.Connect("", "")
	Print("%s connecting... code[%d]", User2, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.Equal(t, User2, connectAck.UserId)
		Print("%s is connected.", User2)
	}
	//user3 connect
	msgContainer3 := map[string]*pbobjs.DownMsg{}
	wsClient3 := wsclients.NewWsImClient(WsAddr, Appkey, Token3, func(msg *pbobjs.DownMsg) {
		msgContainer3[msg.MsgId] = msg
	}, nil, nil)
	code, connectAck = wsClient3.Connect("", "")
	Print("%s connecting... code[%d]", User3, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.Equal(t, User3, connectAck.UserId)
		Print("%s is connected.", User3)
	}
	flag := commonservices.SetStoreMsg(0)
	flag = commonservices.SetCountMsg(flag)
	//send group msg
	code, pubAck := wsClient1.SendGroupMsg(Group1, &pbobjs.UpMsg{
		MsgType:    "test_msg",
		MsgContent: []byte("{\"time:\":\"" + time.Now().Format(TimeFormat) + "\"}"),
		Flags:      flag,
	})
	Print("%s send group msg to %s. code[%d]", User1, Group1, code)
	var grpMsgId string
	var grpMsgTime int64
	var grpMsgIndex int64
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		grpMsgId = pubAck.MsgId
		grpMsgTime = pubAck.Timestamp
		grpMsgIndex = pubAck.MsgSeqNo
	}
	Print("group msg. msg_id[%s] msg_time[%d] msg_index[%d]", grpMsgId, grpMsgTime, grpMsgIndex)
	assert.NotEqual(t, "", grpMsgId)
	assert.NotEqual(t, 0, grpMsgTime)
	assert.NotEqual(t, 0, grpMsgIndex)

	time.Sleep(time.Second)
	//user2 receive msg
	recvMsg2, isRecv2 := msgContainer2[grpMsgId]
	Print("%s receive msg. is_recv[%v]", User2, isRecv2)
	if assert.Equal(t, true, isRecv2) {
		Print("%s receive msg. msg_id[%s]", User2, recvMsg2.MsgId)
		assert.Equal(t, grpMsgId, recvMsg2.MsgId)
	}
	//user3 receive msg
	recvMsg3, isRecv3 := msgContainer3[grpMsgId]
	Print("%s receive msg. is_recv[%v]", User3, isRecv3)
	if assert.Equal(t, true, isRecv3) {
		Print("%s receive msg. msg_id[%s]", User3, recvMsg3.MsgId)
		assert.Equal(t, grpMsgId, recvMsg3.MsgId)
	}

	//sync offline msg from sendbox
	code, msgs := wsClient1.SyncMsgs(&pbobjs.SyncMsgReq{
		SyncTime:        grpMsgTime - 1,
		SendBoxSyncTime: grpMsgTime - 1,
		ContainsSendBox: true,
	})
	Print("%s sync msg from sendbox. code[%d]", User1, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, msgs)
		assert.NotEqual(t, 0, len(msgs.Msgs))
		msg := msgs.Msgs[0]
		Print("send box msg. msg_id[%s] is_send[%v]", msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, true, msg.IsSend)
	}
	//user2 sync offline msg from inbox
	code, msgs = wsClient2.SyncMsgs(&pbobjs.SyncMsgReq{
		SyncTime:        grpMsgTime - 1,
		ContainsSendBox: false,
		SendBoxSyncTime: grpMsgTime - 1,
	})
	Print("%s sync msg from inbox. code[%d]", User2, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, msgs)
		assert.NotEqual(t, 0, len(msgs.Msgs))
		msg := msgs.Msgs[0]
		Print("in box msg. msg_id[%s] is_send[%v]", msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, false, msg.IsSend)
	}
	//user3 sync offline msg from inbox
	code, msgs = wsClient3.SyncMsgs(&pbobjs.SyncMsgReq{
		SyncTime:        grpMsgTime - 1,
		ContainsSendBox: false,
		SendBoxSyncTime: grpMsgTime - 1,
	})
	Print("%s sync msg from inbox. code[%d]", User3, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, msgs)
		assert.NotEqual(t, 0, len(msgs.Msgs))
		msg := msgs.Msgs[0]
		Print("in box msg. msg_id[%s] is_send[%v]", msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, false, msg.IsSend)
	}

	//user1 get history msg
	code, hisMsgs := wsClient1.QryHistoryMsgs(&pbobjs.QryHisMsgsReq{
		TargetId:    Group1,
		ChannelType: pbobjs.ChannelType_Group,
		StartTime:   grpMsgTime + 1,
		Count:       1,
	})
	Print("%s get msg from history. code[%d]", User1, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, hisMsgs)
		assert.NotEqual(t, 0, len(hisMsgs.Msgs))
		msg := hisMsgs.Msgs[0]
		Print("%s qry history msg. msg_id:%s, is_send:%v", User1, msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, true, msg.IsSend)
	}
	//user2 get history msg
	code, hisMsgs = wsClient2.QryHistoryMsgs(&pbobjs.QryHisMsgsReq{
		TargetId:    Group1,
		ChannelType: pbobjs.ChannelType_Group,
		StartTime:   grpMsgTime + 1,
		Count:       1,
	})
	Print("%s get msg from history. code[%d]", User2, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, hisMsgs)
		assert.NotEqual(t, 0, len(hisMsgs.Msgs))
		msg := hisMsgs.Msgs[0]
		Print("%s qry history msg. msg_id:%s, is_send:%v", User2, msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, false, msg.IsSend)
	}
	//user3 get history msg
	code, hisMsgs = wsClient2.QryHistoryMsgs(&pbobjs.QryHisMsgsReq{
		TargetId:    Group1,
		ChannelType: pbobjs.ChannelType_Group,
		StartTime:   grpMsgTime + 1,
		Count:       1,
	})
	Print("%s get msg from history. code[%d]", User3, code)
	if assert.Equal(t, utils.ClientErrorCode_Success, code) {
		assert.NotEqual(t, nil, hisMsgs)
		assert.NotEqual(t, 0, len(hisMsgs.Msgs))
		msg := hisMsgs.Msgs[0]
		Print("%s qry history msg. msg_id:%s, is_send:%v", User3, msg.MsgId, msg.IsSend)
		assert.Equal(t, grpMsgId, msg.MsgId)
		assert.Equal(t, false, msg.IsSend)
	}
	wsClient1.Disconnect()
	wsClient2.Disconnect()
	wsClient3.Disconnect()
}
