package apis

import (
	"encoding/json"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"sort"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func SetGlobalConverTags(ctx *gin.Context) {
	var req models.GlobalConverTagsReq
	if err := ctx.BindJSON(&req); err != nil || req.ConverId == "" || req.ChannelType == int(pbobjs.ChannelType_Unknown) || req.GlobalConverTags == nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	globalConverTags := make(map[string]bool, len(*req.GlobalConverTags))
	for _, tag := range *req.GlobalConverTags {
		if tag != "" {
			globalConverTags[tag] = true
		}
	}
	itemValue, err := json.Marshal(globalConverTags)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "set_g_conver_conf", req.ConverId, &pbobjs.SetConverConfReq{
		ConverId:    req.ConverId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
		SubChannel:  req.SubChannel,
		ItemType:    int32(commonservices.AttItemType_Setting),
		ItemKey:     string(commonservices.AttItemKey_GlobalConverConf_GlobalConverTags),
		ItemValue:   string(itemValue),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func GetGlobalConverTags(ctx *gin.Context) {
	var req models.GlobalConverTagsReq
	if err := ctx.BindJSON(&req); err != nil || req.ConverId == "" || req.ChannelType == int(pbobjs.ChannelType_Unknown) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_g_conver_conf", req.ConverId, &pbobjs.ConverConfReq{
		ConverId:    req.ConverId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
		SubChannel:  req.SubChannel,
	}, func() proto.Message {
		return &pbobjs.ConverConf{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	tags := make([]string, 0)
	if converConf, ok := resp.(*pbobjs.ConverConf); ok {
		for tag, enabled := range converConf.GlobalConverTags {
			if enabled {
				tags = append(tags, tag)
			}
		}
	}
	sort.Strings(tags)
	tools.SuccessHttpResp(ctx, models.GlobalConverTagsResp{GlobalConverTags: tags})
}

func AddConversation(ctx *gin.Context) {
	var req models.Conversation
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	bases.SyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "add_conver", req.UserId, &pbobjs.Conversation{
		UserId:      req.UserId,
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func DelConversation(ctx *gin.Context) {
	var req models.Conversations
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(req.Items) <= 0 || req.UserId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.Conversation{}
	for _, c := range req.Items {
		convers = append(convers, &pbobjs.Conversation{
			UserId:      c.UserId,
			TargetId:    c.TargetId,
			ChannelType: pbobjs.ChannelType(c.ChannelType),
		})
	}

	bases.SyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "del_convers", req.UserId, &pbobjs.ConversationsReq{
		Conversations: convers,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func ClearConverUnread(ctx *gin.Context) {
	var req models.Conversations
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(req.Items) <= 0 || req.UserId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.Conversation{}
	for _, c := range req.Items {
		convers = append(convers, &pbobjs.Conversation{
			UserId:      req.UserId,
			TargetId:    c.TargetId,
			ChannelType: pbobjs.ChannelType(c.ChannelType),
		})
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "clear_unread", req.UserId, &pbobjs.ClearUnreadReq{
		Conversations: convers,
		NoCmdMsg:      true,
	})
	tools.SuccessHttpResp(ctx, nil)
}

// undisturb_convers
func UndisturbConvers(ctx *gin.Context) {
	var undisturbConversReq models.UndisturbConversReq
	if err := ctx.BindJSON(&undisturbConversReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(undisturbConversReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	items := []*pbobjs.UndisturbConverItem{}
	for _, reqItem := range undisturbConversReq.Items {
		items = append(items, &pbobjs.UndisturbConverItem{
			TargetId:      reqItem.TargetId,
			ChannelType:   pbobjs.ChannelType(reqItem.ChannelType),
			UndisturbType: reqItem.UndisturbType,
		})
	}
	bases.SyncRpcCall(services.ToRpcCtx(ctx, undisturbConversReq.UserId), "undisturb_convers", undisturbConversReq.UserId, &pbobjs.UndisturbConversReq{
		Items: items,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

// top convers
func TopConversations(ctx *gin.Context) {
	var topConversReq models.TopConversReq
	if err := ctx.BindJSON(&topConversReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if topConversReq.UserId == "" || len(topConversReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	items := []*pbobjs.Conversation{}
	for _, reqItem := range topConversReq.Items {
		isTopInt := 0
		if reqItem.IsTop {
			isTopInt = 1
		}
		items = append(items, &pbobjs.Conversation{
			TargetId:    reqItem.TargetId,
			ChannelType: pbobjs.ChannelType(reqItem.ChannelType),
			IsTop:       int32(isTopInt),
		})
	}
	bases.SyncRpcCall(services.ToRpcCtx(ctx, topConversReq.UserId), "top_convers", topConversReq.UserId, &pbobjs.ConversationsReq{
		Conversations: items,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func QryConvers(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	start := ctx.Query("start")
	var startTime int64 = 0
	if start != "" {
		intVal, err := tools.String2Int64(start)
		if err == nil {
			startTime = intVal
		}
	}

	limitStr := ctx.Query("count")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	orderStr := ctx.Query("order")
	var order int32 = 0
	if orderStr != "" {
		intVal, err := tools.String2Int64(orderStr)
		if err == nil && intVal > 0 {
			order = 1
		}
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, userId), "qry_convers", userId, &pbobjs.QryConversationsReq{
		StartTime:  startTime,
		Count:      int32(limit),
		Order:      order,
		OnlyConver: true,
	}, func() proto.Message {
		return &pbobjs.QryConversationsResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	conversResp := resp.(*pbobjs.QryConversationsResp)
	ret := models.Conversations{
		Items:      []*models.Conversation{},
		IsFinished: conversResp.IsFinished,
	}
	for _, conver := range conversResp.Conversations {
		item := &models.Conversation{
			UserId:      conver.UserId,
			TargetId:    conver.TargetId,
			ChannelType: int(conver.ChannelType),
			SubChannel:  conver.SubChannel,
			Time:        conver.SortTime,
		}
		ret.Items = append(ret.Items, item)
	}
	tools.SuccessHttpResp(ctx, ret)
}

func QryGlobalConvers(ctx *gin.Context) {
	start := ctx.Query("start")
	var startTime int64 = 0
	if start != "" {
		intVal, err := tools.String2Int64(start)
		if err == nil {
			startTime = intVal
		}
	}

	limitStr := ctx.Query("count")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	rpcTargetId := fmt.Sprintf("random%d", tools.RandInt(1000))
	//targetId
	targetId := ctx.Query("target_id")
	//channelType
	channelTypeStr := ctx.Query("channel_type")
	channelTypeInt, err := tools.String2Int64(channelTypeStr)
	channelType := pbobjs.ChannelType_Unknown
	if err == nil {
		channelType = pbobjs.ChannelType(channelTypeInt)
	}
	//exclude user_ids
	excludeUserIds := ctx.QueryArray("exclude_user_id")
	if len(excludeUserIds) == 0 {
		excludeUserIds = ctx.QueryArray("exclude_user_ids")
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_global_convers", rpcTargetId, &pbobjs.QryGlobalConversReq{
		Start:          startTime,
		Order:          0,
		Count:          int32(limit),
		TargetId:       targetId,
		ChannelType:    channelType,
		ExcludeUserIds: excludeUserIds,
	}, func() proto.Message {
		return &pbobjs.QryGlobalConversResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	conversResp := resp.(*pbobjs.QryGlobalConversResp)
	ret := models.Conversations{
		Items:      []*models.Conversation{},
		IsFinished: conversResp.IsFinished,
	}
	for _, conver := range conversResp.Convers {
		item := &models.Conversation{
			Id:          conver.Id,
			UserId:      conver.SenderId,
			TargetId:    conver.TargetId,
			ChannelType: int(conver.ChannelType),
			Time:        conver.UpdatedTime,
		}
		ret.Items = append(ret.Items, item)
	}

	tools.SuccessHttpResp(ctx, ret)
}

func CreateConverTag(ctx *gin.Context) {
	var req models.CreateConverTagReq
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" || req.Tag == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, req.UserId), "create_user_conver_tags", req.UserId, &pbobjs.UserConverTags{
		Tags: []*pbobjs.ConverTag{
			{
				Tag:      req.Tag,
				TagName:  req.TagName,
				TagType:  pbobjs.ConverTagType_UserConverTag,
				TagOrder: req.TagOrder,
			},
		},
	}, func() proto.Message {
		return &pbobjs.UserConverTags{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func TagConvers(ctx *gin.Context) {
	var tagConvers models.TagConversReq
	if err := ctx.BindJSON(&tagConvers); err != nil || len(tagConvers.Convers) <= 0 || tagConvers.UserId == "" || tagConvers.Tag == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.SimpleConversation{}
	for _, conver := range tagConvers.Convers {
		convers = append(convers, &pbobjs.SimpleConversation{
			TargetId:    conver.TargetId,
			ChannelType: pbobjs.ChannelType(conver.ChannelType),
			SubChannel:  conver.SubChannel,
		})
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, tagConvers.UserId), "tag_add_convers", tagConvers.UserId, &pbobjs.TagConvers{
		Tag:     tagConvers.Tag,
		TagName: tagConvers.TagName,
		Convers: convers,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func UnTagConvers(ctx *gin.Context) {
	var tagConvers models.TagConversReq
	if err := ctx.BindJSON(&tagConvers); err != nil || len(tagConvers.Convers) <= 0 || tagConvers.UserId == "" || tagConvers.Tag == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.SimpleConversation{}
	for _, conver := range tagConvers.Convers {
		convers = append(convers, &pbobjs.SimpleConversation{
			TargetId:    conver.TargetId,
			ChannelType: pbobjs.ChannelType(conver.ChannelType),
			SubChannel:  conver.SubChannel,
		})
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, tagConvers.UserId), "tag_del_convers", tagConvers.UserId, &pbobjs.TagConvers{
		Tag:     tagConvers.Tag,
		TagName: tagConvers.TagName,
		Convers: convers,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}
