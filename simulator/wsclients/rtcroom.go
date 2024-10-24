package wsclients

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
)

func (client *WsImClient) CreateRtcRoom(room *pbobjs.RtcInviteReq) (utils.ClientErrorCode, *pbobjs.RtcRoom) {
	data, _ := tools.PbMarshal(room)
	code, qryAck := client.Query("rtc_create", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		resp := &pbobjs.RtcRoom{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return utils.ClientErrorCode_Success, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) DestroyRtcRoom(roomId string) utils.ClientErrorCode {
	code, _ := client.Query("rtc_destroy", roomId, []byte{})
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) JoinRtcRoom(room *pbobjs.RtcInviteReq) (utils.ClientErrorCode, *pbobjs.RtcRoom) {
	data, _ := tools.PbMarshal(room)
	code, qryAck := client.Query("rtc_join", room.RoomId, data)
	if code == utils.ClientErrorCode_Success || code == 16002 {
		resp := &pbobjs.RtcRoom{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return code, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) QuitRtcRoom(roomId string) utils.ClientErrorCode {
	code, _ := client.Query("rtc_quit", roomId, []byte{})
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) QryRtcRoom(roomId string) (utils.ClientErrorCode, *pbobjs.RtcRoom) {
	code, qryAck := client.Query("rtc_qry", roomId, []byte{})
	if code == utils.ClientErrorCode_Success {
		resp := &pbobjs.RtcRoom{}
		tools.PbUnMarshal(qryAck.Data, resp)
		return code, resp
	} else {
		return utils.ClientErrorCode(code), nil
	}
}

func (client *WsImClient) RtcRoomPing(roomId string) utils.ClientErrorCode {
	code, _ := client.Query("rtc_ping", roomId, []byte{})
	return utils.ClientErrorCode(code)
}

func (client *WsImClient) RtcInvite(req *pbobjs.RtcInviteReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("rtc_invite", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(code)
	}
}

func (client *WsImClient) RtcDecline(req *pbobjs.RtcAnswerReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("rtc_decline", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(code)
	}
}

func (client *WsImClient) RtcAccept(req *pbobjs.RtcAnswerReq) utils.ClientErrorCode {
	data, _ := tools.PbMarshal(req)
	code, qryAck := client.Query("rtc_accept", client.UserId, data)
	if code == utils.ClientErrorCode_Success && qryAck.Code == 0 {
		return utils.ClientErrorCode_Success
	} else {
		return utils.ClientErrorCode(code)
	}
}
