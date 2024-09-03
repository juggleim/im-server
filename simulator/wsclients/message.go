package wsclients

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"
)

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
