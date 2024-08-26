package apis

import (
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func CreateChatroom(ctx *gin.Context) {
	var req models.ChatroomInfo
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_create", "", req.ChatId, &pbobjs.ChatroomInfo{
		ChatId:   req.ChatId,
		ChatName: req.ChatName,
		IsMute:   req.IsMute > 0,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, nil)
}

func DestroyChatroom(ctx *gin.Context) {
	var req models.ChatroomInfo
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_destroy", "", req.ChatId, &pbobjs.ChatroomInfo{
		ChatId:   req.ChatId,
		ChatName: req.ChatName,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tools.SuccessHttpResp(ctx, nil)
}

func QryChatroomInfo(ctx *gin.Context) {
	chatId := ctx.Query("chat_id")
	orderStr := ctx.Query("order")
	var order int32 = 0
	if orderStr != "" {
		intVal, err := tools.String2Int64(orderStr)
		if err == nil {
			order = int32(intVal)
		}
	}
	countStr := ctx.Query("count")
	var count int32 = 100
	if countStr != "" {
		intVal, err := tools.String2Int64(countStr)
		if err == nil && intVal > 0 {
			count = int32(intVal)
		}
	}
	code, resp, err := services.SyncApiCall(ctx, "c_qry_chrm", "", chatId, &pbobjs.ChatroomReq{
		ChatId: chatId,
		Count:  count,
		Order:  order,
	}, func() proto.Message {
		return &pbobjs.ChatroomInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	chrmInfo := resp.(*pbobjs.ChatroomInfo)
	var isMute int = 0
	if chrmInfo.IsMute {
		isMute = 1
	}
	ret := &models.ChatroomInfo{
		ChatId:      chatId,
		ChatName:    chrmInfo.ChatName,
		MemberCount: chrmInfo.MemberCount,
		IsMute:      isMute,
		Members:     []*models.ChatroomMember{},
		Atts:        []*models.ChatroomAtt{},
	}
	for _, member := range chrmInfo.Members {
		ret.Members = append(ret.Members, &models.ChatroomMember{
			MemberId:   member.MemberId,
			MemberName: member.MemberName,
			AddedTime:  member.AddedTime,
		})
	}
	for _, att := range chrmInfo.Atts {
		ret.Atts = append(ret.Atts, &models.ChatroomAtt{
			Key:     att.Key,
			Value:   att.Value,
			UserId:  att.UserId,
			AttTime: att.AttTime,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}

func ChrmMute(ctx *gin.Context) {
	var muteReq models.ChatroomInfo
	if err := ctx.BindJSON(&muteReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "chrm_mute", "", muteReq.ChatId, &pbobjs.ChatroomInfo{
		ChatId: muteReq.ChatId,
		IsMute: muteReq.IsMute > 0,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func AddChrmMuteMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Mute,
		MemberIds: req.MemberIds,
		IsDelete:  false,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelChrmMuteMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Mute,
		MemberIds: req.MemberIds,
		IsDelete:  true,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryChrmMuteMembers(ctx *gin.Context) {
	chatId := ctx.Query("chat_id")
	offset := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 {
			limit = intVal
		}
	}
	if limit > 1000 {
		limit = 1000
	}
	code, resp, err := services.SyncApiCall(ctx, "c_qry_ban_user", "", chatId, &pbobjs.QryChrmBanUsersReq{
		ChatId:  chatId,
		BanType: pbobjs.ChrmBanType_Mute,
		Offset:  offset,
		Limit:   limit,
	}, func() proto.Message {
		return &pbobjs.QryChrmBanUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	ret := &models.ChrmBanUsers{
		ChatId:  chatId,
		Members: []*models.ChatroomMember{},
	}
	if banUsersResp, ok := resp.(*pbobjs.QryChrmBanUsersResp); ok {
		ret.ChatId = banUsersResp.ChatId
		ret.Offset = banUsersResp.Offset
		for _, member := range banUsersResp.Members {
			ret.Members = append(ret.Members, &models.ChatroomMember{
				MemberId:  member.MemberId,
				AddedTime: member.CreatedTime,
			})
		}
	}

	tools.SuccessHttpResp(ctx, ret)
}

func AddChrmBanMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Ban,
		MemberIds: req.MemberIds,
		IsDelete:  false,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelChrmBanMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Ban,
		MemberIds: req.MemberIds,
		IsDelete:  true,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryChrmBanMembers(ctx *gin.Context) {
	chatId := ctx.Query("chat_id")
	offset := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 {
			limit = intVal
		}
	}
	if limit > 1000 {
		limit = 1000
	}
	code, resp, err := services.SyncApiCall(ctx, "c_qry_ban_user", "", chatId, &pbobjs.QryChrmBanUsersReq{
		ChatId:  chatId,
		BanType: pbobjs.ChrmBanType_Ban,
		Offset:  offset,
		Limit:   limit,
	}, func() proto.Message {
		return &pbobjs.QryChrmBanUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	ret := &models.ChrmBanUsers{
		ChatId:  chatId,
		Members: []*models.ChatroomMember{},
	}
	if banUsersResp, ok := resp.(*pbobjs.QryChrmBanUsersResp); ok {
		ret.ChatId = banUsersResp.ChatId
		ret.Offset = banUsersResp.Offset
		for _, member := range banUsersResp.Members {
			ret.Members = append(ret.Members, &models.ChatroomMember{
				MemberId:  member.MemberId,
				AddedTime: member.CreatedTime,
			})
		}
	}

	tools.SuccessHttpResp(ctx, ret)
}

func AddChrmAllowMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Allow,
		MemberIds: req.MemberIds,
		IsDelete:  false,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func DelChrmAllowMembers(ctx *gin.Context) {
	var req models.ChrmBanUserReq
	if err := ctx.BindJSON(&req); err != nil || req.ChatId == "" || len(req.MemberIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "c_ban_user", "", req.ChatId, &pbobjs.BatchBanUserReq{
		ChatId:    req.ChatId,
		BanType:   pbobjs.ChrmBanType_Allow,
		MemberIds: req.MemberIds,
		IsDelete:  true,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryChrmAllowMembers(ctx *gin.Context) {
	chatId := ctx.Query("chat_id")
	offset := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 {
			limit = intVal
		}
	}
	if limit > 1000 {
		limit = 1000
	}
	code, resp, err := services.SyncApiCall(ctx, "c_qry_ban_user", "", chatId, &pbobjs.QryChrmBanUsersReq{
		ChatId:  chatId,
		BanType: pbobjs.ChrmBanType_Allow,
		Offset:  offset,
		Limit:   limit,
	}, func() proto.Message {
		return &pbobjs.QryChrmBanUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	ret := &models.ChrmBanUsers{
		ChatId:  chatId,
		Members: []*models.ChatroomMember{},
	}
	if banUsersResp, ok := resp.(*pbobjs.QryChrmBanUsersResp); ok {
		ret.ChatId = banUsersResp.ChatId
		ret.Offset = banUsersResp.Offset
		for _, member := range banUsersResp.Members {
			ret.Members = append(ret.Members, &models.ChatroomMember{
				MemberId:  member.MemberId,
				AddedTime: member.CreatedTime,
			})
		}
	}

	tools.SuccessHttpResp(ctx, ret)
}
