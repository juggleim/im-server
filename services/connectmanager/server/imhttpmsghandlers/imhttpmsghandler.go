package imhttpmsghandlers

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"im-server/services/connectmanager/server/codec"
	"io"
	"net/http"
	"time"
)

func ImHttpPubHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	method := ""
	tokenStr := ""
	appkey := ""
	deviceId := ""
	instanceId := ""
	platform := ""
	if r != nil {
		method = r.Method
		tokenStr = r.Header.Get("x-token")
		appkey = r.Header.Get("x-appkey")
		deviceId = r.Header.Get("x-deviceid")
		instanceId = r.Header.Get("x-instanceid")
		platform = r.Header.Get("x-platform")
	}

	if appkey == "" || tokenStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBs, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBs) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	wsMsg := &codec.ImWebsocketMsg{}
	err = tools.PbUnMarshal(bodyBs, wsMsg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pubMsgBody := wsMsg.GetPublishMsgBody()
	if pubMsgBody == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//check token
	userId, code := parseToken(appkey, tokenStr)
	if code != errs.IMErrorCode_SUCCESS {
		responseAck(w, codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index: pubMsgBody.Index,
			Code:  int32(code),
		}))
		return
	}

	if pubMsgBody.Topic == "" || pubMsgBody.TargetId == "" {
		responseAck(w, codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index: pubMsgBody.Index,
			Code:  int32(errs.IMErrorCode_CONNECT_PARAM_REQUIRED),
		}))
		return
	}

	resp, err := bases.SyncUnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_UserPub,
		AppKey:       appkey,
		Session:      "im_" + tools.GenerateUUIDShort11(),
		DeviceId:     deviceId,
		InstanceId:   instanceId,
		Platform:     platform,
		Method:       "upstream",
		RequesterId:  pubMsgBody.TargetId,
		ReqIndex:     pubMsgBody.Index,
		Qos:          wsMsg.Qos,
		AppDataBytes: pubMsgBody.Data,
		TargetId:     userId,
		ExtParams: map[string]string{
			commonservices.RpcExtKey_RealMethod: pubMsgBody.Topic,
		},
	}, 5*time.Second)

	if err != nil || resp == nil {
		responseAck(w, codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
			Index: pubMsgBody.Index,
			Code:  int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC),
		}))
		return
	}

	var afterModified *codec.SimplifiedDownMsg
	if resp.ModifiedMsg != nil {
		afterModified = &codec.SimplifiedDownMsg{
			MsgType:    resp.ModifiedMsg.MsgType,
			MsgContent: resp.ModifiedMsg.MsgContent,
		}
	}

	ack := codec.NewUserPublishAckMessage(&codec.PublishAckMsgBody{
		Index:       pubMsgBody.Index,
		Code:        resp.ResultCode,
		MsgId:       resp.MsgId,
		Timestamp:   resp.MsgSendTime,
		MsgSeqNo:    resp.MsgSeqNo,
		MemberCount: resp.MemberCount,
		ClientMsgId: resp.ClientMsgId,
		ModifiedMsg: afterModified,
	})
	responseAck(w, ack)
}

func responseAck(w http.ResponseWriter, ack *codec.UserPublishAckMessage) {
	w.WriteHeader(http.StatusOK)
	bs, _ := tools.PbMarshal(ack.ToImWebsocketMsg())
	w.Write(bs)
}

func parseToken(appkey, tokenStr string) (string, errs.IMErrorCode) {
	if tokenStr == "" {
		return "", errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL
	}
	tokenWrap, err := tokens.ParseTokenString(tokenStr)
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL
	}
	if tokenWrap.AppKey != appkey {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	appInfo, exist := commonservices.GetAppInfo(appkey)
	if !exist || appInfo == nil {
		return "", errs.IMErrorCode_CONNECT_APP_NOT_EXISTED
	}
	token, err := tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
	if err != nil {
		return "", errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL
	}
	if appInfo.TokenEffectiveMinute > 0 && (token.TokenTime+int64(appInfo.TokenEffectiveMinute)*60*1000) < time.Now().UnixMilli() {
		return "", errs.IMErrorCode_CONNECT_TOKEN_EXPIRED
	}
	return token.UserId, errs.IMErrorCode_SUCCESS
}
