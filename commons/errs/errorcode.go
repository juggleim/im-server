package errs

/*
0 : success
10000~10999 : api
11000~11999 : connect
12000~12999 : private msg
13000~13999 : group
14000~14999 : chatroom
*/

type IMErrorCode int32

var IMErrorCode_SUCCESS IMErrorCode = 0
var IMErrorCode_PBILLEGAL IMErrorCode = 1 //pb解析失败，内部错误码
var IMErrorCode_DEFAULT IMErrorCode = 2

// api errorcode
var (
	IMErrorCode_API_DEFAULT            IMErrorCode = 10000
	IMErrorCode_API_APPKEY_REQUIRED    IMErrorCode = 10001
	IMErrorCode_API_NONCE_REQUIRED     IMErrorCode = 10002
	IMErrorCode_API_TIMESTAMP_REQUIRED IMErrorCode = 10003
	IMErrorCode_API_SIGNATURE_REQUIRED IMErrorCode = 10004
	IMErrorCode_API_APP_NOT_EXISTED    IMErrorCode = 10005
	IMErrorCode_API_SIGNATURE_FAIL     IMErrorCode = 10006
	IMErrorCode_API_REQ_BODY_ILLEGAL   IMErrorCode = 10007
	IMErrorCode_API_INTERNAL_TIMEOUT   IMErrorCode = 10008
	IMErrorCode_API_INTERNAL_RESP_FAIL IMErrorCode = 10009
	IMErrorCode_API_PARAM_REQUIRED     IMErrorCode = 10010
	IMErrorCode_API_PARAM_ILLEGAL      IMErrorCode = 10011
)

// user errorcode
var (
	IMErrorCode_USER_DEFAULT             IMErrorCode = 10100
	IMErrorCode_USER_COUNT_EXCEED        IMErrorCode = 10101
	IMErrorCode_USER_NOT_EXIST           IMErrorCode = 10102
	IMErrorCode_USER_NOT_SUPPROT_SETTING IMErrorCode = 10103
	IMErrorCode_USER_TIMEZONE_ILLGAL     IMErrorCode = 10104
	IMErrorCode_USER_EXISTED             IMErrorCode = 10106
)

// connect errorcode
var (
	IMErrorCode_CONNECT_DEFAULT                 IMErrorCode = 11000
	IMErrorCode_CONNECT_APPKEY_REQUIRED         IMErrorCode = 11001
	IMErrorCode_CONNECT_TOKEN_REQUIRED          IMErrorCode = 11002
	IMErrorCode_CONNECT_APP_NOT_EXISTED         IMErrorCode = 11003
	IMErrorCode_CONNECT_TOKEN_ILLEGAL           IMErrorCode = 11004
	IMErrorCode_CONNECT_TOKEN_AUTHFAIL          IMErrorCode = 11005
	IMErrorCode_CONNECT_TOKEN_EXPIRED           IMErrorCode = 11006
	IMErrorCode_CONNECT_NEED_REDIRECT           IMErrorCode = 11007
	IMErrorCode_CONNECT_UNSUPPROTEDPLATFORM     IMErrorCode = 11008
	IMErrorCode_CONNECT_APP_BLOCK               IMErrorCode = 11009
	IMErrorCode_CONNECT_USER_BLOCK              IMErrorCode = 11010
	IMErrorCode_CONNECT_KICKED_OFF              IMErrorCode = 11011 //被踢下线
	IMErrorCode_CONNECT_LOGOUT                  IMErrorCode = 11012 //注销
	IMErrorCode_CONNECT_UNSUPPORTEDTOPIC        IMErrorCode = 11013
	IMErrorCode_CONNECT_EXCEEDLIMITED           IMErrorCode = 11014
	IMErrorCode_CONNECT_PARAM_REQUIRED          IMErrorCode = 11015
	IMErrorCode_CONNECT_CLOSE_NET_ERR           IMErrorCode = 11016
	IMErrorCode_CONNECT_CLOSE_PB_DECODE_FAIL    IMErrorCode = 11017
	IMErrorCode_CONNECT_CLOSE_HEARTBEAT_TIMEOUT IMErrorCode = 11018
	IMErrorCode_CONNECT_CLOSE_DATA_ILLEGAL      IMErrorCode = 11019
	IMErrorCode_CONNECT_UNSECURITYDOMAIN        IMErrorCode = 11020
	IMErrorCode_CONNECT_FUNCTIONDISABLED        IMErrorCode = 11021
	IMErrorCode_CONNECT_KICKED_BY_SELF          IMErrorCode = 11022
)

// msg errorcode
var (
	IMErrorCode_MSG_DEFAULT         IMErrorCode = 12000
	IMErrorCode_MSG_ADDFAILED       IMErrorCode = 12001
	IMErrorCode_MSG_DELFAILED       IMErrorCode = 12002
	IMErrorCode_MSG_UPDFAILED       IMErrorCode = 12003
	IMErrorCode_MSG_PARAM_ILLEGAL   IMErrorCode = 12004
	IMErrorCode_MSG_BLOCK           IMErrorCode = 12005
	IMErrorCode_MSG_MSGEXTDUPLICATE IMErrorCode = 12006
	IMErrorCode_MSG_Hit_Sensitive   IMErrorCode = 12007
	IMErrorCode_MSG_MSGEXTOVERLIMIT IMErrorCode = 12008
)

// conversation
var (
	IMErrorCode_CONVER_ADDTAGFAIL       IMErrorCode = 12101
	IMErrorCode_CONVER_TAGADDCONVERFAIL IMErrorCode = 12102
)

// group errorcode
var (
	IMErrorCode_GROUP_DEFAULT                IMErrorCode = 13000
	IMErrorCode_GROUP_GROUPNOTEXIST          IMErrorCode = 13001
	IMErrorCode_GROUP_NOTGROUPMEMBER         IMErrorCode = 13002
	IMErrorCode_GROUP_GROUPMUTE              IMErrorCode = 13003
	IMErrorCode_GROUP_GROUPMEMBERMUTE        IMErrorCode = 13004
	IMErrorCode_GROUP_GROUPMEMBERCOUNTEXCEED IMErrorCode = 13005
	IMErrorCode_GROUP_NOSNAPSHOT             IMErrorCode = 13006
)

// chatroom errorcode
var (
	IMErrorCode_CHATROOM_DEFAULT     IMErrorCode = 14000
	IMErrorCode_CHATROOM_NOTMEMBER   IMErrorCode = 14001
	IMErrorCode_CHATROOM_ATTFULL     IMErrorCode = 14002
	IMErrorCode_CHATROOM_SEIZEFAILED IMErrorCode = 14003
	IMErrorCode_CHATROOM_ATTNOTEXIST IMErrorCode = 14004
	IMErrorCode_CHATROOM_NOTEXIST    IMErrorCode = 14005
	IMErrorCode_CHATROOM_HASDELETED  IMErrorCode = 14006
	IMErrorCode_CHATROOM_MUTE        IMErrorCode = 14007
	IMErrorCode_CHATROOM_BAN         IMErrorCode = 14008
)

// other
var (
	IMErrorCode_OTHER_DEFAULT IMErrorCode = 15000
	IMErrorCode_OTHER_NOOSS   IMErrorCode = 15001
	IMErrorCode_OTHER_SIGNERR IMErrorCode = 15002

	IMErrorCode_OTHER_CLIENT_LOG_PARAM_ILLEGAL IMErrorCode = 15101
)

// rtcroom errorcode
var (
	IMErrorCode_RTCROOM_DEFAULT            IMErrorCode = 16000
	IMErrorCode_RTCROOM_ROOMNOTEXIST       IMErrorCode = 16001
	IMErrorCode_RTCROOM_ROOMHASEXIST       IMErrorCode = 16002
	IMErrorCode_RTCROOM_ROOMHASDELETED     IMErrorCode = 16003
	IMErrorCode_RTCROOM_NOTMEMBER          IMErrorCode = 16004
	IMErrorCode_RTCROOM_HASMEMBER          IMErrorCode = 16005
	IMErrorCode_RTCROOM_CREATEROOMFAILED   IMErrorCode = 16006
	IMErrorCode_RTCROOM_UPDATEFAILED       IMErrorCode = 16007
	IMErrorCode_RTCROOM_RTCAUTHFAILED      IMErrorCode = 16008
	IMErrorCode_RTCROOM_INACTIVERTCCHANNEL IMErrorCode = 16009
	IMErrorCode_RTCROOM_PARAMILLIGAL       IMErrorCode = 16010

	IMErrorCode_RTCINVITE_HASACCEPT  IMErrorCode = 16100
	IMErrorCode_RTCINVITE_REJECT     IMErrorCode = 16101
	IMErrorCode_RTCINVITE_PEERHANGUP IMErrorCode = 16102
	IMErrorCode_RTCINVITE_BUSY       IMErrorCode = 16103
	IMErrorCode_RTCINVITE_CANCEL     IMErrorCode = 16104
)

// app errorcode
var (
	IMErrorCode_APP_DEFAULT          IMErrorCode = 17000
	IMErrorCode_APP_APPKEY_REQUIRED  IMErrorCode = 17001
	IMErrorCode_APP_NOT_EXISTED      IMErrorCode = 17002
	IMErrorCode_APP_REQ_BODY_ILLEGAL IMErrorCode = 17003
	IMErrorCode_APP_INTERNAL_TIMEOUT IMErrorCode = 17004
	IMErrorCode_APP_NOT_LOGIN        IMErrorCode = 17005
	IMErrorCode_APP_CONTINUE         IMErrorCode = 17006
	IMErrorCode_APP_QRCODE_EXPIRED   IMErrorCode = 17007
	IMErrorCode_APP_SMS_SEND_FAILED  IMErrorCode = 17008
	IMErrorCode_APP_SMS_CODE_EXPIRED IMErrorCode = 17009

	//friends
	IMErrorCode_APP_FRIEND_DEFAULT         IMErrorCode = 17100
	IMErrorCode_APP_FRIEND_APPLY_DECLINE   IMErrorCode = 17101
	IMErrorCode_APP_FRIEND_APPLY_REPEATED  IMErrorCode = 17102
	IMErrorCode_APP_FRIEND_CONFIRM_EXPIRED IMErrorCode = 17103

	//group
	IMErrorCode_APP_GROUP_DEFAULT       IMErrorCode = 17200
	IMErrorCode_APP_GROUP_MEMBEREXISTED IMErrorCode = 17201

	//assistant
	IMErrorCode_APP_ASSISTANT_PROMPT_DBERROR IMErrorCode = 17300
)

var imCode2ApiErrorMap map[IMErrorCode]*ApiErrorMsg = map[IMErrorCode]*ApiErrorMsg{
	//api
	IMErrorCode_SUCCESS:                newApiErrorMsg(200, IMErrorCode_SUCCESS, "success"),
	IMErrorCode_API_DEFAULT:            newApiErrorMsg(200, IMErrorCode_API_DEFAULT, "default error"),
	IMErrorCode_API_APPKEY_REQUIRED:    newApiErrorMsg(400, IMErrorCode_API_APPKEY_REQUIRED, "appkey is required"),
	IMErrorCode_API_NONCE_REQUIRED:     newApiErrorMsg(400, IMErrorCode_API_NONCE_REQUIRED, "nonce is required"),
	IMErrorCode_API_TIMESTAMP_REQUIRED: newApiErrorMsg(400, IMErrorCode_API_TIMESTAMP_REQUIRED, "timestamp is required"),
	IMErrorCode_API_SIGNATURE_REQUIRED: newApiErrorMsg(400, IMErrorCode_API_SIGNATURE_REQUIRED, "signature is required"),
	IMErrorCode_API_APP_NOT_EXISTED:    newApiErrorMsg(500, IMErrorCode_API_APP_NOT_EXISTED, "app not existed"),
	IMErrorCode_API_SIGNATURE_FAIL:     newApiErrorMsg(403, IMErrorCode_API_SIGNATURE_FAIL, "signature is fail"),
	IMErrorCode_API_REQ_BODY_ILLEGAL:   newApiErrorMsg(400, IMErrorCode_API_REQ_BODY_ILLEGAL, "request body is illegal"),
	IMErrorCode_API_INTERNAL_TIMEOUT:   newApiErrorMsg(500, IMErrorCode_API_INTERNAL_TIMEOUT, "internal service timeout"),
	IMErrorCode_API_INTERNAL_RESP_FAIL: newApiErrorMsg(500, IMErrorCode_API_INTERNAL_RESP_FAIL, "internal service error"),
	IMErrorCode_API_PARAM_REQUIRED:     newApiErrorMsg(400, IMErrorCode_API_PARAM_REQUIRED, "required params missing."),

	//conn
	IMErrorCode_CONNECT_DEFAULT:         newApiErrorMsg(200, IMErrorCode_CONNECT_DEFAULT, "default error"),
	IMErrorCode_CONNECT_APPKEY_REQUIRED: newApiErrorMsg(200, IMErrorCode_CONNECT_APPKEY_REQUIRED, "appkey is required"),
	IMErrorCode_CONNECT_TOKEN_REQUIRED:  newApiErrorMsg(200, IMErrorCode_CONNECT_TOKEN_REQUIRED, "token is required"),
	IMErrorCode_CONNECT_APP_NOT_EXISTED: newApiErrorMsg(200, IMErrorCode_CONNECT_APP_NOT_EXISTED, "app not exist"),
	IMErrorCode_CONNECT_TOKEN_ILLEGAL:   newApiErrorMsg(401, IMErrorCode_CONNECT_TOKEN_ILLEGAL, "token illegal"),
	IMErrorCode_CONNECT_TOKEN_AUTHFAIL:  newApiErrorMsg(401, IMErrorCode_CONNECT_TOKEN_AUTHFAIL, "token auth fail"),
	IMErrorCode_CONNECT_TOKEN_EXPIRED:   newApiErrorMsg(200, IMErrorCode_CONNECT_TOKEN_EXPIRED, "token expired"),
}

func GetApiErrorByCode(code IMErrorCode) *ApiErrorMsg {
	if err, ok := imCode2ApiErrorMap[code]; ok {
		return err
	}
	return newApiErrorMsg(200, code, "")
}

type ApiErrorMsg struct {
	HttpCode int         `json:"-"`
	Code     IMErrorCode `json:"code"`
	Msg      string      `json:"msg"`
}

func newApiErrorMsg(httpCode int, code IMErrorCode, msg string) *ApiErrorMsg {
	return &ApiErrorMsg{
		HttpCode: httpCode,
		Code:     code,
		Msg:      msg,
	}
}

/*
type ErrorCode int32
const (
	ERR_SUCCESS ErrorCode = 0
	//api
	ERR_API_DEFAULT ErrorCode = 10000
	//conn
	ERR_CONN_DEFAULT         ErrorCode = 11000
	ERR_CONN_APPKEY_REQUIRED ErrorCode = 11001
	ERR_CONN_TOKEN_REQUIRED  ErrorCode = 11002
	ERR_CONN_APP_NOT_EXISTED ErrorCode = 11003
	ERR_CONN_TOKEN_ILLEGAL   ErrorCode = 11004
	ERR_CONN_TOKEN_AUTHFAIL  ErrorCode = 11005
	ERR_CONN_TOKEN_EXPIRED   ErrorCode = 11006
)
*/
