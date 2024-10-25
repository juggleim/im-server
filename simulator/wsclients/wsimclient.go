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

	DisconnectCallback  func(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody)
	OnMessageCallBack   func(msg *pbobjs.DownMsg)
	OnStreamMsgCallBack func(msg *pbobjs.StreamDownMsg)

	UserId string

	conn            *websocket.Conn
	state           utils.ConnectState
	accssorCache    sync.Map
	myIndex         uint16
	connAckAccessor *tools.DataAccessor
	pongAccessor    *tools.DataAccessor
	lock            *sync.RWMutex

	inboxTime   int64
	sendboxTime int64

	obfCode   [8]byte
	isEncrypt bool
}

func NewWsImClient(address, appkey, token string, onMessage func(msg *pbobjs.DownMsg), onStreamMsg func(*pbobjs.StreamDownMsg), onDisconnect func(code utils.ClientErrorCode, disMsg *codec.DisconnectMsgBody)) *WsImClient {
	return &WsImClient{
		Address:             address,
		Appkey:              appkey,
		Token:               token,
		accssorCache:        sync.Map{},
		connAckAccessor:     tools.NewDataAccessor(),
		pongAccessor:        tools.NewDataAccessorWithSize(100),
		OnMessageCallBack:   onMessage,
		OnStreamMsgCallBack: onStreamMsg,
		DisconnectCallback:  onDisconnect,
		Platform:            "Web", // "Android",
		DeviceId:            "testDevice",
		obfCode:             [8]byte{1, 2, 3, 4, 5, 6, 7, 8},
		isEncrypt:           true,
		lock:                &sync.RWMutex{},
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
			go client.startPing()
			return clientCode, connAck
		} else {
			return clientCode, connAck
		}
	} else {
		return utils.ClientErrorCode_ConnectExisted, nil
	}
}

func (client *WsImClient) WriteMessage(data []byte) error {
	if client.state == utils.State_Connected {
		client.lock.Lock()
		defer client.lock.Unlock()
		return client.conn.WriteMessage(websocket.BinaryMessage, data)
	}
	return fmt.Errorf("not connected")
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

func (client *WsImClient) startPing() {
	for client.state != utils.State_Disconnect {
		client.Ping()
		time.Sleep(30 * time.Second)
	}
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
		client.WriteMessage(wsMsgBs)
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
		client.WriteMessage(wsMsgBs)
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
		if needAck > 0 {
			ackMsg := codec.NewServerPublishAckMessage(&codec.PublishAckMsgBody{
				Index: msg.Index,
			})
			wsMsg := ackMsg.ToImWebsocketMsg()
			Encrypt(wsMsg, client)
			wsMsgBs, _ := tools.PbMarshal(wsMsg)
			client.WriteMessage(wsMsgBs)
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
			}
		} else {
			fmt.Println("error:", err)
		}
	} else if msg.Topic == "stream_msg" {
		streamMsg := pbobjs.StreamDownMsg{}
		err := tools.PbUnMarshal(msg.Data, &streamMsg)
		if err == nil && client.OnStreamMsgCallBack != nil {
			client.OnStreamMsgCallBack(&streamMsg)
		}
	} else if msg.Topic == "rtc_room_event" {
		var event pbobjs.RtcRoomEvent
		err := tools.PbUnMarshal(msg.Data, &event)
		if err != nil {
			fmt.Println("err:", err)
		} else {
			fmt.Println(tools.ToJson(&event))
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
		client.WriteMessage(wsMsgBs)
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
		client.WriteMessage(wsMsgBs)
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
		client.WriteMessage(wsMsgBs)
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
