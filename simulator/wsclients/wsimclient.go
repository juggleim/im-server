package wsclients

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"

	"github.com/gorilla/websocket"
)

type WsImClient struct {
	Address         string
	Token           string
	Appkey          string
	Platform        string
	DeviceId        string
	DeviceCompany   string
	DeviceModel     string
	DeviceOsVersion string
	PushToken       string

	DisconnectCallback func(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody)
	OnMessageCallBack  func(msg *pbobjs.DownMsg)

	UserId string

	conn            *websocket.Conn
	state           utils.ConnectState
	accssorCache    sync.Map
	myIndex         uint16
	connAckAccessor *tools.DataAccessor
	pongAccessor    *tools.DataAccessor

	inboxTime   int64
	sendboxTime int64

	obfCode   [8]byte
	isEncrypt bool
}

func NewWsImClient(address, appkey, token string, onMessage func(msg *pbobjs.DownMsg), onDisconnect func(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody)) *WsImClient {
	return &WsImClient{
		Address:            address,
		Appkey:             appkey,
		Token:              token,
		accssorCache:       sync.Map{},
		connAckAccessor:    tools.NewDataAccessor(),
		pongAccessor:       tools.NewDataAccessorWithSize(100),
		OnMessageCallBack:  onMessage,
		DisconnectCallback: onDisconnect,
		Platform:           "Web", // "Android",
		DeviceId:           "testDevice",
		obfCode:            [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		isEncrypt:          true,
	}
}

func (client *WsImClient) GetState() utils.ConnectState {
	return client.state
}

func (client *WsImClient) Connect(network, ispNum string) (utils.ClientErrorCode, *codec.ConnectAckMsgBody) {
	if client.state == utils.State_Disconnect {
		//init ws client
		var u url.URL
		if strings.HasPrefix(client.Address, "wss://") {
			addr := client.Address[6:]
			u = url.URL{Scheme: "wss", Host: addr, Path: "/im"}
		} else if strings.HasPrefix(client.Address, "ws://") {
			addr := client.Address[5:]
			u = url.URL{Scheme: "ws", Host: addr, Path: "/im"}
		}
		header := http.Header{}
		// header.Add("aaa", "test123")
		// header.Add("token", "aabbccaaaa")
		// query := u.Query()
		// query.Set("param1", "value1")
		// u.RawQuery = query.Encode()

		c, _, err := websocket.DefaultDialer.Dial(u.String(), header)
		if err != nil {
			fmt.Println("addr:", u.String(), "err:", err)
			return utils.ClientErrorCode_SocketFailed, nil
		}
		client.conn = c

		connectMsg := codec.NewConnectMessage(&codec.ConnectMsgBody{
			ProtoId:         codec.ProtoId,
			SdkVersion:      "1.0.1",
			Appkey:          client.Appkey,
			Token:           client.Token,
			Platform:        client.Platform,
			DeviceId:        client.DeviceId,
			DeviceCompany:   client.DeviceCompany,
			DeviceModel:     client.DeviceModel,
			DeviceOsVersion: client.DeviceOsVersion,
			PushToken:       client.PushToken,
			PushChannel:     "1",
			NetworkId:       network,
			IspNum:          ispNum,
		})
		wsMsg := connectMsg.ToImWebsocketMsg()

		//encrypt
		Encrypt(wsMsg, client)

		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		err = client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		if err != nil {
			fmt.Println(err)
			return utils.ClientErrorCode_ConnectTimeout, nil
		}
		client.state = utils.State_Connecting
		go client.startListener()

		connAckObj, err := client.connAckAccessor.GetWithTimeout(10 * time.Second)
		if err != nil {
			fmt.Println(err)
			return utils.ClientErrorCode_ConnectTimeout, nil
		}
		connAck := connAckObj.(*codec.ConnectAckMsgBody)
		clientCode := utils.Trans2ClientErrorCoce(connAck.Code)
		if connAck.Code == int32(errs.IMErrorCode_SUCCESS) { //链接成功
			client.UserId = connAck.UserId
			client.state = utils.State_Connected
			return clientCode, connAck
		} else {
			return clientCode, connAck
		}
	} else {
		return utils.ClientErrorCode_ConnectExisted, nil
	}
}

func (client *WsImClient) startListener() {
	for client.state != utils.State_Disconnect {
		_, msgBs, err := client.conn.ReadMessage()
		if err == nil {
			wsImMsg := &codec.ImWebsocketMsg{}
			err = tools.PbUnMarshal(msgBs, wsImMsg)
			if err == nil {
				Decrypt(wsImMsg, client)
				go func() {
					switch wsImMsg.Cmd {
					case int32(codec.Cmd_ConnectAck):
						client.OnConnectAck(wsImMsg.GetConnectAckMsgBody())
					case int32(codec.Cmd_Disconnect):
						client.OnDisconnect(wsImMsg.GetDisconnectMsgBody())
					case int32(codec.Cmd_Publish):
						client.OnPublish(wsImMsg.GetPublishMsgBody(), int(wsImMsg.Qos))
					case int32(codec.Cmd_PublishAck):
						client.OnPublishAck(wsImMsg.GetPubAckMsgBody())
					case int32(codec.Cmd_QueryAck):
						client.OnQueryAck(wsImMsg.GetQryAckMsgBody())
					}
				}()
			}
		} else {
			client.state = utils.State_Disconnect
		}
	}
	fmt.Println("Stop client listener.")
}

func (client *WsImClient) Reconnect(network, ispNum string) (utils.ClientErrorCode, *codec.ConnectAckMsgBody) {
	if client.state != utils.State_Connecting {
		if client.conn != nil {
			client.conn.Close()
		}
		return client.Connect(network, ispNum)
	} else {
		return utils.ClientErrorCode_ConnectExisted, nil
	}
}

func (client *WsImClient) Disconnect() {
	if client.conn != nil {
		disMsg := codec.NewDisconnectMessage(&codec.DisconnectMsgBody{
			Code: int32(errs.IMErrorCode_SUCCESS),
		})
		wsMsg := disMsg.ToImWebsocketMsg()
		Encrypt(wsMsg, client)
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		client.conn.Close()
		client.conn = nil
	}
	client.state = utils.State_Disconnect
}

func (client *WsImClient) Logout() {
	if client.conn != nil {
		disMsg := codec.NewDisconnectMessage(&codec.DisconnectMsgBody{
			Code: int32(errs.IMErrorCode_CONNECT_LOGOUT),
		})
		wsMsg := disMsg.ToImWebsocketMsg()
		Encrypt(wsMsg, client)
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		client.conn.Close()
		client.conn = nil
	}
	client.state = utils.State_Disconnect
}

func (client *WsImClient) OnConnectAck(msg *codec.ConnectAckMsgBody) {
	client.connAckAccessor.Put(msg)
	if msg.Code > 0 {
		client.state = utils.State_Disconnect
	}
}

func (client *WsImClient) OnPublishAck(msg *codec.PublishAckMsgBody) {
	dataAccessor, ok := client.accssorCache.LoadAndDelete(msg.Index)
	if ok {
		dataAccessor.(*tools.DataAccessor).Put(msg)
	}
}
func (client *WsImClient) OnQueryAck(msg *codec.QueryAckMsgBody) {
	dataAccessor, ok := client.accssorCache.LoadAndDelete(msg.Index)
	if ok {
		dataAccessor.(*tools.DataAccessor).Put(msg)
	}
}
func (client *WsImClient) OnDisconnect(msg *codec.DisconnectMsgBody) {
	if client.DisconnectCallback != nil {
		client.DisconnectCallback(utils.Trans2ClientErrorCoce(msg.Code), msg)
	}
}
func (client *WsImClient) OnPong(msg *codec.ImWebsocketMsg) {
	client.pongAccessor.Put(msg)
}
func (client *WsImClient) OnPublish(msg *codec.PublishMsgBody, needAck int) {
	if needAck > 0 {
		ackMsg := codec.NewServerPublishAckMessage(&codec.PublishAckMsgBody{
			Index: msg.Index,
		})
		wsMsg := ackMsg.ToImWebsocketMsg()
		Encrypt(wsMsg, client)
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
	}
	if msg.Topic == "msg" {
		downMsg := pbobjs.DownMsg{}
		err := tools.PbUnMarshal(msg.Data, &downMsg)
		if err == nil && client.OnMessageCallBack != nil {
			client.OnMessageCallBack(&downMsg)
			if downMsg.IsSend {
				client.sendboxTime = downMsg.MsgTime
			} else {
				client.inboxTime = downMsg.MsgTime
			}
		}
	} else if msg.Topic == "ntf" {
		ntf := pbobjs.Notify{}
		err := tools.PbUnMarshal(msg.Data, &ntf)
		if err == nil {
			if ntf.Type == pbobjs.NotifyType_Msg {
				isContinue := true
				for isContinue {
					code, downSet := client.SyncMsgs(&pbobjs.SyncMsgReq{
						SyncTime:        client.inboxTime,
						SendBoxSyncTime: client.sendboxTime,
					})
					if code == utils.ClientErrorCode(errs.IMErrorCode_SUCCESS) {
						for _, downMsg := range downSet.Msgs {
							client.OnMessageCallBack(downMsg)
						}
						if downSet.IsFinished {
							isContinue = false
						}
					} else {
						fmt.Println("ntf pull msg error, code:", code)
						isContinue = false
					}
				}
			} else if ntf.Type == pbobjs.NotifyType_ChatroomMsg {
				code, downSet := client.SyncChatroomMsgs(&pbobjs.SyncChatroomReq{
					SyncTime:   0,
					ChatroomId: ntf.ChatroomId,
				})
				if code == utils.ClientErrorCode(errs.IMErrorCode_SUCCESS) {
					for _, downMsg := range downSet.Msgs {
						client.OnMessageCallBack(downMsg)
					}

				} else {
					fmt.Println("ntf pull msg error, code:", code)
				}
			}
		} else {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println(msg.Topic, msg.Data)
	}
}

func (client *WsImClient) Publish(method, targetId string, data []byte) (code utils.ClientErrorCode, pubAck *codec.PublishAckMsgBody) {
	if client.state == utils.State_Connected {
		index := client.getMyIndex()
		protoMsg := codec.NewUserPublishMessage(&codec.PublishMsgBody{
			Index:    int32(index),
			Topic:    method,
			TargetId: targetId,
			Data:     data,
		})
		dataAccessor := tools.NewDataAccessor()
		client.accssorCache.Store(int32(index), dataAccessor)

		wsMsg := protoMsg.ToImWebsocketMsg()
		Encrypt(wsMsg, client)
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		obj, err := dataAccessor.GetWithTimeout(10 * time.Second)
		if err == nil {
			pubAck := obj.(*codec.PublishAckMsgBody)
			return utils.Trans2ClientErrorCoce(pubAck.Code), pubAck
		} else {
			return utils.ClientErrorCode_SendTimeout, nil
		}
	} else {
		return utils.ClientErrorCode_ConnectClosed, nil
	}
}

func (client *WsImClient) Ping() utils.ClientErrorCode {
	if client.state == utils.State_Connected {
		pingMsg := codec.NewPingMessage()
		wsMsg := pingMsg.ToImWebsocketMsg()
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		_, err := client.pongAccessor.GetWithTimeout(15 * time.Second)
		if err == nil {
			return utils.ClientErrorCode_Success
		} else {
			return utils.ClientErrorCode_PingTimeout
		}
	} else {
		return utils.ClientErrorCode_ConnectClosed
	}
}

func (client *WsImClient) Query(method, targetId string, data []byte) (utils.ClientErrorCode, *codec.QueryAckMsgBody) {
	if client.state == utils.State_Connected {
		index := int32(client.getMyIndex())
		protoMsg := codec.NewQueryMessage(&codec.QueryMsgBody{
			Index:    index,
			Topic:    method,
			TargetId: targetId,
			Data:     data,
		})
		dataAccessor := tools.NewDataAccessor()
		client.accssorCache.Store(index, dataAccessor)

		wsMsg := protoMsg.ToImWebsocketMsg()
		Encrypt(wsMsg, client)
		wsMsgBs, _ := tools.PbMarshal(wsMsg)
		client.conn.WriteMessage(websocket.BinaryMessage, wsMsgBs)
		obj, err := dataAccessor.GetWithTimeout(10 * time.Second)
		if err == nil {
			queryAck := obj.(*codec.QueryAckMsgBody)
			return utils.Trans2ClientErrorCoce(queryAck.Code), queryAck
		} else {
			return utils.ClientErrorCode_QueryTimeout, nil
		}
	} else {
		return utils.ClientErrorCode_ConnectClosed, nil
	}
}

func (client *WsImClient) getMyIndex() uint16 {
	client.myIndex = client.myIndex + 1
	return client.myIndex
}

func (client *WsImClient) SendPrivateMsg(targetId string, upMsg *pbobjs.UpMsg) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(upMsg)
	code, pubAck := client.Publish("p_msg", targetId, data)
	return code, pubAck
}

func (client *WsImClient) SendGroupMsg(targetId string, upMsg *pbobjs.UpMsg) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(upMsg)
	code, pubAck := client.Publish("g_msg", targetId, data)
	return code, pubAck
}

func (client *WsImClient) SendChatMsg(targetId string, upMsg *pbobjs.UpMsg) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(upMsg)
	code, pubAck := client.Publish("c_msg", targetId, data)
	return code, pubAck
}

func (client *WsImClient) AddChatAtt(targetId string, att *pbobjs.ChatAttReq) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(att)
	code, pubAck := client.Publish("c_add_att", targetId, data)
	return code, pubAck
}

func (client *WsImClient) DelChatAtt(targetId string, att *pbobjs.ChatAttReq) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(att)
	code, pubAck := client.Publish("c_del_att", targetId, data)
	return code, pubAck
}

func (client *WsImClient) QryHistoryMsgs(qryHisMsgReq *pbobjs.QryHisMsgsReq) (utils.ClientErrorCode, *pbobjs.DownMsgSet) {
	data, _ := tools.PbMarshal(qryHisMsgReq)
	code, qryAck := client.Query("qry_hismsgs", qryHisMsgReq.TargetId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		downMsgSet := &pbobjs.DownMsgSet{}
		tools.PbUnMarshal(qryAck.Data, downMsgSet)
		return utils.ClientErrorCode_Success, downMsgSet
	} else {
		return utils.ClientErrorCode_Unknown, nil
	}
}

func (client *WsImClient) QryFirstUnreadMsg(req *pbobjs.QryFirstUnreadMsgReq) (utils.ClientErrorCode, *pbobjs.DownMsg) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_first_unread_msg", req.TargetId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		msg := &pbobjs.DownMsg{}
		tools.PbUnMarshal(qryAck.Data, msg)
		return utils.ClientErrorCode_Success, msg
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) DelHisMsgs(delHisMsgs *pbobjs.DelHisMsgsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(delHisMsgs)
	code, qryAck := client.Query("del_hismsg", delHisMsgs.TargetId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(qryAck.Code)
	}
}

func (client *WsImClient) QryConversation(req *pbobjs.QryConverReq) (utils.ClientErrorCode, *pbobjs.QryConverResp) {
	data, _ := tools.PbMarshal(req)
	code, ack := client.Query("qry_conver", client.UserId, data)
	if code == utils.ClientErrorCode_Success && ack.Code == 0 {
		resp := &pbobjs.QryConverResp{}
		tools.PbUnMarshal(ack.Data, resp)
		return utils.ClientErrorCode_Success, resp
	}
	return code, nil
}

func (client *WsImClient) QryConversations(req *pbobjs.QryConversationsReq) (utils.ClientErrorCode, *pbobjs.QryConversationsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_convers", client.UserId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryConversationsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) ClearUnread(req *pbobjs.ClearUnreadReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("clear_unread", client.UserId, data)
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) QryMentionMsgs(req *pbobjs.QryMentionMsgsReq) (utils.ClientErrorCode, *pbobjs.QryMentionMsgsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_mention_msgs", client.UserId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryMentionMsgsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) SyncConversations(req *pbobjs.QryConversationsReq) (utils.ClientErrorCode, *pbobjs.QryConversationsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("sync_convers", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryConversationsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) SyncMsgs(req *pbobjs.SyncMsgReq) (utils.ClientErrorCode, *pbobjs.DownMsgSet) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("sync_msgs", client.UserId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.DownMsgSet{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode_Unknown, nil
	}
}

func (client *WsImClient) RecallMsg(req *pbobjs.RecallMsgReq) (utils.ClientErrorCode, *codec.QueryAckMsgBody) {
	data, _ := tools.PbMarshal(req)
	code, pubAck := client.Query("recall_msg", req.TargetId, data)
	return code, pubAck
}

func (client *WsImClient) ModifyMsg(req *pbobjs.ModifyMsgReq) (utils.ClientErrorCode, *codec.QueryAckMsgBody) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("modify_msg", req.TargetId, data)
	return code, qryAck
}

func (client *WsImClient) MarkReadMsg(req *pbobjs.MarkReadReq) (utils.ClientErrorCode, *codec.QueryAckMsgBody) {
	data, _ := tools.PbMarshal(req)
	code, ack := client.Query("mark_read", client.UserId, data)
	return code, ack
}

func (client *WsImClient) QryHisMsgsByIds(targetId string, req *pbobjs.QryHisMsgByIdsReq) (utils.ClientErrorCode, *codec.QueryAckMsgBody) {
	data, _ := tools.PbMarshal(req)
	return client.Query("qry_hismsg_by_ids", targetId, data)
}

// chatroom
func (client *WsImClient) JoinChatroom(chatroomId string) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(&pbobjs.ChatroomInfo{
		ChatId: chatroomId,
	})
	code, _ := client.Query("c_join", chatroomId, data)
	if code == utils.ClientErrorCode_Success {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(code)
	}
}
func (client *WsImClient) QuitChatroom(chatroomId string) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(&pbobjs.ChatroomInfo{
		ChatId: chatroomId,
	})
	code, _ := client.Query("c_quit", chatroomId, data)
	if code == utils.ClientErrorCode_Success {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(code)
	}
}

func (client *WsImClient) SendChatroomMsg(chatroomId string, upMsg *pbobjs.UpMsg) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(upMsg)
	code, pubAck := client.Publish("c_msg", chatroomId, data)
	return code, pubAck
}

func (client *WsImClient) SyncChatroomMsgs(req *pbobjs.SyncChatroomReq) (utils.ClientErrorCode, *pbobjs.SyncChatroomMsgResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("c_sync_msgs", req.ChatroomId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.SyncChatroomMsgResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		fmt.Println("server code:", code, qryAck)
		return utils.ClientErrorCode_Unknown, nil
	}
}

func (client *WsImClient) GetFileCred(req *pbobjs.QryFileCredReq) (utils.ClientErrorCode, *pbobjs.QryFileCredResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("file_cred", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryFileCredResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) QryUserInfo(req *pbobjs.UserIdReq) (utils.ClientErrorCode, *pbobjs.UserInfo) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_userinfo", client.UserId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.UserInfo{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode(code), resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) QryReadMsgDetail(req *pbobjs.QryReadDetailReq) (utils.ClientErrorCode, *pbobjs.QryReadDetailResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_read_detail", req.TargetId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryReadDetailResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode(code), resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) UndisturbConvers(req *pbobjs.UndisturbConversReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("undisturb_convers", client.UserId, data)
	return code
}

func (client *WsImClient) SetTopConvers(req *pbobjs.ConversationsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("top_convers", client.UserId, data)
	return code
}

func (client *WsImClient) DelConvers(req *pbobjs.ConversationsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("del_convers", client.UserId, data)
	return code
}

func (client *WsImClient) CleanHisMsgs(req *pbobjs.CleanHisMsgReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("clean_hismsg", client.UserId, data)
	return code
}

func (client *WsImClient) QryMergedMsgs(msgId string, req *pbobjs.QryMergedMsgsReq) (utils.ClientErrorCode, *pbobjs.DownMsgSet) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_merged_msgs", msgId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.DownMsgSet{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode(code), resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) SubUsers(req *pbobjs.UserIdsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("sub_users", client.UserId, data)
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) UnSubUsers(req *pbobjs.UserIdsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("unsub_users", client.UserId, data)
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) SetUserUndisturb(req *pbobjs.UserUndisturb) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("set_user_undisturb", client.UserId, data)
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) GetUserUndisturb() (utils.ClientErrorCode, *pbobjs.UserUndisturb) {
	data, _ := tools.PbMarshal(&pbobjs.Nil{})
	code, qryAck := client.Query("get_user_undisturb", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.UserUndisturb{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode(code), resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}
