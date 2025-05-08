package imcontext

var ServiceName string = "connectmanager"

type ActionType string

const (
	Action_Query        ActionType = "qry"         //0
	Action_QueryAck     ActionType = "qry_ack"     //1
	Action_QueryConfirm ActionType = "qry_confirm" //2

	Action_UserPub    ActionType = "u_pub"     //3
	Action_UserPubAck ActionType = "u_pub_ack" //4

	Action_ServerPub    ActionType = "s_pub"     //5
	Action_ServerPubAck ActionType = "s_pub_ack" //6

	Action_Connect    ActionType = "connect"     //7
	Action_Disconnect ActionType = "disconnect"  //8
	Action_ConnectErr ActionType = "connect_err" //9
)

const (
	StateKey_Connected            string = "state.connected" //1:success; 2:failed
	StateKey_ObfuscationCode      string = "state.obfuscation_code"
	StateKey_CtxLocker            string = "state.ctx_locker"
	StateKey_ServerPubCallbackMap string = "state.pub_callback_map"
	StateKey_QueryConfirmMap      string = "state.qry_confirm_map"
	StateKey_ConnectSession       string = "state.connect_session"
	StateKey_ConnectCreateTime    string = "state.connect_timestamp"
	StateKey_ServerMsgIndex       string = "state.server_msg_index"
	StateKey_ClientMsgIndex       string = "state.client_msg_index"
	StateKey_Appkey               string = "state.appkey"
	StateKey_UserID               string = "state.userid"
	StateKey_DeviceID             string = "state.device_id"
	StateKey_Platform             string = "state.platform"
	StateKey_Version              string = "state.version"
	StateKey_ClientIp             string = "state.client_ip"
	StateKey_Limiter              string = "state.limiter"
	StateKey_Referer              string = "state.referer"
	// StateKey_Extra                string = "state.extra"
	StateKey_InstanceId string = "state.instance_id"
)

type CloseReason int

const (
	Close_Normal CloseReason = 0
)

type Attachment interface{}

type WsHandleContext interface {
	Write(message interface{})
	Close(err error)
	Attachment() Attachment
	SetAttachment(attachment Attachment)
	IsActive() bool
	RemoteAddr() string
}
