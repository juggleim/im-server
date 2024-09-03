package wsclients

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
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
