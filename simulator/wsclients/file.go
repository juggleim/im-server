package wsclients

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
)

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
