package wsclients

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
)

func (client *WsImClient) QryConversation(req *pbobjs.QryConverReq) (utils.ClientErrorCode, *pbobjs.Conversation) {
	data, _ := tools.PbMarshal(req)
	code, ack := client.Query("qry_conver", client.UserId, data)
	if code == utils.ClientErrorCode_Success && ack.Code == 0 {
		resp := &pbobjs.Conversation{}
		tools.PbUnMarshal(ack.Data, resp)
		return utils.ClientErrorCode_Success, resp
	}
	return code, nil
}

func (client *WsImClient) QryTotalUnreadCount(req *pbobjs.QryTotalUnreadCountReq) (utils.ClientErrorCode, *pbobjs.QryTotalUnreadCountResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_total_unread_count", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryTotalUnreadCountResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) QryConversations(req *pbobjs.QryConversationsReq) (utils.ClientErrorCode, *pbobjs.QryConversationsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_convers", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryConversationsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) QryTopConvers(req *pbobjs.QryTopConversReq) (utils.ClientErrorCode, *pbobjs.QryConversationsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_top_convers", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryConversationsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) ClearUnread(req *pbobjs.ClearUnreadReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("clear_unread", client.UserId, data)
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) QryMentionMsgs(req *pbobjs.QryMentionMsgsReq) (utils.ClientErrorCode, *pbobjs.QryMentionMsgsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("qry_mention_msgs", client.UserId, data)
	fmt.Println(code)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryMentionMsgsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) SyncConversations(req *pbobjs.QryConversationsReq) (utils.ClientErrorCode, *pbobjs.QryConversationsResp) {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("sync_convers", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.QryConversationsResp{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) UndisturbConvers(req *pbobjs.UndisturbConversReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("undisturb_convers", client.UserId, data)
	return code
}

func (client *WsImClient) SetTopConvers(req *pbobjs.ConversationsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("top_convers", client.UserId, data)
	return code
}

func (client *WsImClient) DelConvers(req *pbobjs.ConversationsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("del_convers", client.UserId, data)
	return code
}

func (client *WsImClient) MarkUnRead(req *pbobjs.ConversationsReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("mark_unread", client.UserId, data)
	return code
}

func (client *WsImClient) TagAddConvers(req *pbobjs.TagConvers) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, _ := client.Query("tag_add_convers", client.UserId, data)
	return code
}
