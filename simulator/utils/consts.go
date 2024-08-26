package utils

type ConnectState uint8

const (
	State_Disconnect ConnectState = 0
	State_Connecting ConnectState = 1
	State_Connected  ConnectState = 2
)

type ClientErrorCode int

const (
	ClientErrorCode_Success        ClientErrorCode = 0
	ClientErrorCode_Unknown        ClientErrorCode = 20000
	ClientErrorCode_SocketFailed   ClientErrorCode = 20001
	ClientErrorCode_ConnectTimeout ClientErrorCode = 20002
	ClientErrorCode_NeedRedirect   ClientErrorCode = 20003
	ClientErrorCode_ConnectExisted ClientErrorCode = 20004
	ClientErrorCode_PingTimeout    ClientErrorCode = 20005
	ClientErrorCode_ConnectFailed  ClientErrorCode = 20006
	ClientErrorCode_ConnectClosed  ClientErrorCode = 20007

	ClientErrorCode_SendTimeout ClientErrorCode = 21001

	ClientErrorCode_QueryTimeout ClientErrorCode = 22001
)

func Trans2ClientErrorCoce(serverCode int32) ClientErrorCode {
	return ClientErrorCode(serverCode)
}
