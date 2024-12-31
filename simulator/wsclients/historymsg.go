package wsclients

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
	"im-server/simulator/utils"
)

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
	code, qryAck := client.Query("del_msg", delHisMsgs.TargetId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(qryAck.Code)
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

func (client *WsImClient) BatchTranslate(req *pbobjs.TransReq) (utils.ClientErrorCode, *pbobjs.TransReq) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("batch_trans", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.TransReq{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode(code), resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}
