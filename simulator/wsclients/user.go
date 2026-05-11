package wsclients

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"
)

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

// SubUsers 调用 Query「sub_users」，成功时解析 Ack 中的 UserStatusList（含在线状态）。
func (client *WsImClient) SubUsers(req *pbobjs.SubUsersReq) (utils.ClientErrorCode, *pbobjs.UserStatusList) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("sub_user_status", client.UserId, data)
	if code != utils.ClientErrorCode_Success || qryAck == nil {
		return code, nil
	}
	if qryAck.Code != 0 {
		return utils.ClientErrorCode(qryAck.Code), nil
	}
	resp := &pbobjs.UserStatusList{}
	tools.PbUnMarshal(qryAck.Data, resp)
	return utils.ClientErrorCode_Success, resp
}

func (client *WsImClient) UnSubUsers(req *pbobjs.SubUsersReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("unsub_users", client.UserId, data)
	if code != utils.ClientErrorCode_Success || qryAck == nil {
		return code
	}
	if qryAck.Code != 0 {
		return utils.ClientErrorCode(qryAck.Code)
	}
	return utils.ClientErrorCode_Success
}

// PubUserStatus 对应「pub_user_status」上游发布（statussubscriptions.PubUserStatusActor）。
func (client *WsImClient) PubUserStatus(upMsg *pbobjs.UpMsg) (utils.ClientErrorCode, *codec.PublishAckMsgBody) {
	data, _ := tools.PbMarshal(upMsg)
	code, pubAck := client.Publish("pub_user_status", client.UserId, data)
	return code, pubAck
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
