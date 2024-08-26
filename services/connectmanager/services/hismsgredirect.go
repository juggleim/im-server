package services

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
)

type IGetChannelType interface {
	GetChannelType() pbobjs.ChannelType
}

func HisMsgRedirect(method string, data []byte, fromId, targetId string) (string, bool) {
	var typeGetter IGetChannelType
	switch method {
	case "clean_hismsg":
		var req pbobjs.CleanHisMsgReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "del_msg":
		var req pbobjs.DelHisMsgsReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "mark_read":
		var req pbobjs.MarkReadReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "modify_msg":
		var req pbobjs.ModifyMsgReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "qry_hismsgs":
		var req pbobjs.QryHisMsgsReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "qry_first_unread_msg":
		var req pbobjs.QryFirstUnreadMsgReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "qry_hismsg_by_ids":
		var req pbobjs.QryHisMsgByIdsReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "qry_read_infos":
		var req pbobjs.QryReadInfosReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	case "recall_msg":
		var req pbobjs.RecallMsgReq
		err := tools.PbUnMarshal(data, &req)
		if err == nil {
			typeGetter = &req
		}
	}
	if typeGetter != nil && (typeGetter.GetChannelType() == pbobjs.ChannelType_Private || typeGetter.GetChannelType() == pbobjs.ChannelType_System) {
		converId := commonservices.GetConversationId(fromId, targetId, typeGetter.GetChannelType())
		fmt.Println("converId:", converId)
		return converId, true
	}
	return "", false
}
