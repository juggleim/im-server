package apis

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func Register(ctx *gin.Context) {
	var userInfo models.UserInfo
	if err := ctx.BindJSON(&userInfo); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "reg_user", userInfo.UserId, &pbobjs.UserInfo{
		UserId:       userInfo.UserId,
		Nickname:     userInfo.Nickname,
		UserPortrait: userInfo.UserPortrait,
		ExtFields:    commonservices.Map2KvItems(userInfo.ExtFields),
	}, func() proto.Message {
		return &pbobjs.UserRegResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	rpcResp, ok := resp.(*pbobjs.UserRegResp)
	if !ok {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	tools.SuccessHttpResp(ctx, models.UserRegResp{
		UserId: rpcResp.UserId,
		Token:  rpcResp.Token,
	})
}

func UpdateUser(ctx *gin.Context) {
	var req models.UserInfo
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "upd_user_info", req.UserId, &pbobjs.UserInfo{
		UserId:       req.UserId,
		Nickname:     req.Nickname,
		UserPortrait: req.UserPortrait,
		ExtFields:    commonservices.Map2KvItems(req.ExtFields),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func SetUserSettings(ctx *gin.Context) {
	var req models.UserSettings
	if err := ctx.BindJSON(&req); err != nil || req.UserId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	kvMap := make(map[string]string)
	for k, v := range req.Settings {
		if commonservices.CheckUserSettingKey(k) {
			kvMap[k] = fmt.Sprintf("%v", v)
		}
	}
	if len(kvMap) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "set_user_settings", req.UserId, &pbobjs.UserInfo{
		UserId:   req.UserId,
		Settings: commonservices.Map2KvItems(kvMap),
	})
	tools.SuccessHttpResp(ctx, nil)
}

func GetUserSettings(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	if userId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, userInfo, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "get_user_settings", userId, &pbobjs.Nil{}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != errs.IMErrorCode_SUCCESS {
		tools.ErrorHttpResp(ctx, code)
		return
	}
	ret := &models.UserSettings{
		UserId:   userId,
		Settings: map[string]interface{}{},
	}
	if userInfo != nil {
		uInfo, ok := userInfo.(*pbobjs.UserInfo)
		if ok && uInfo != nil {
			for _, setting := range uInfo.Settings {
				ret.Settings[setting.Key] = setting.Value
			}
		}
	}
	tools.SuccessHttpResp(ctx, ret)
}

func QryUserInfo(ctx *gin.Context) {
	userid := ctx.Query("user_id")
	if userid == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	code, userInfoObj, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, userid), "qry_user_info", userid, &pbobjs.UserIdReq{
		UserId: userid,
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	userInfo := userInfoObj.(*pbobjs.UserInfo)
	tools.SuccessHttpResp(ctx, &models.UserInfo{
		UserId:       userInfo.UserId,
		Nickname:     userInfo.Nickname,
		UserPortrait: userInfo.UserPortrait,
		ExtFields:    commonservices.Kvitems2Map(userInfo.ExtFields),
		UpdatedTime:  userInfo.UpdatedTime,
	})
}

func KickUsers(ctx *gin.Context) {
	var req models.KickUserReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "kick_user", req.UserId, &pbobjs.KickUserReq{
		UserId:    req.UserId,
		Platforms: req.Platforms,
		DeviceIds: req.DeviceIds,
		Ext:       req.Ext,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryUserOnlineStatus(ctx *gin.Context) {
	var userOnlineReq models.UserOnlineStatusReq
	if err := ctx.BindJSON(&userOnlineReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	if len(userOnlineReq.UserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	ret := models.UserOnlineStatusResp{
		Items: []*models.UserOnlineStatusItem{},
	}
	tmpMap := sync.Map{}
	groups := bases.GroupTargets("qry_online_status", userOnlineReq.UserIds)
	wg := sync.WaitGroup{}
	for _, ids := range groups {
		wg.Add(1)
		userIds := ids
		go func() {
			defer wg.Done()
			_, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_online_status", userIds[0], &pbobjs.UserOnlineStatusReq{
				UserIds: userIds,
			}, func() proto.Message {
				return &pbobjs.UserOnlineStatusResp{}
			})
			if err == nil {
				onlineResp, ok := resp.(*pbobjs.UserOnlineStatusResp)
				if ok && len(onlineResp.Items) > 0 {
					for _, item := range onlineResp.Items {
						tmpMap.Store(item.UserId, item)
					}
				}
			}
		}()
	}
	wg.Wait()
	tmpMap.Range(func(key, value any) bool {
		item := value.(*pbobjs.UserOnlineItem)
		ret.Items = append(ret.Items, &models.UserOnlineStatusItem{
			UserId:   item.UserId,
			IsOnline: item.IsOnline,
		})
		return true
	})
	tools.SuccessHttpResp(ctx, ret)
}

func UserBan(ctx *gin.Context) {
	var banReq models.BanUsersReq
	if err := ctx.BindJSON(&banReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(banReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	groups := map[string][]*pbobjs.BanUser{}
	for _, user := range banReq.Items {
		node := bases.GetCluster().GetTargetNode("ban_users", user.UserId)
		if node != nil && node.Name != "" {
			var arr []*pbobjs.BanUser
			var ok bool
			endTime := user.EndTime
			if endTime <= 0 && user.EndTimeOffset > 0 {
				endTime = time.Now().UnixMilli() + user.EndTimeOffset
			}
			pbBanUser := &pbobjs.BanUser{
				UserId:     user.UserId,
				EndTime:    endTime,
				ScopeKey:   user.ScopeKey,
				ScopeValue: user.ScopeValue,
				Ext:        user.Ext,
			}
			if arr, ok = groups[node.Name]; ok {
				arr = append(arr, pbBanUser)
			} else {
				arr = []*pbobjs.BanUser{pbBanUser}
			}
			groups[node.Name] = arr
		}
	}
	wg := sync.WaitGroup{}
	for _, banUsers := range groups {
		wg.Add(1)
		users := banUsers
		go func() {
			defer wg.Done()
			bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "ban_users", users[0].UserId, &pbobjs.BanUsersReq{
				BanUsers: users,
				IsAdd:    true,
			})
		}()
	}
	wg.Wait()
	tools.SuccessHttpResp(ctx, nil)
}

func UserUnBan(ctx *gin.Context) {
	var banReq models.BanUsersReq
	if err := ctx.BindJSON(&banReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(banReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	groups := map[string][]*pbobjs.BanUser{}
	for _, user := range banReq.Items {
		node := bases.GetCluster().GetTargetNode("ban_users", user.UserId)
		if node != nil && node.Name != "" {
			var arr []*pbobjs.BanUser
			var ok bool
			pbBanUser := &pbobjs.BanUser{
				UserId:     user.UserId,
				EndTime:    user.EndTime,
				ScopeKey:   user.ScopeKey,
				ScopeValue: user.ScopeValue,
			}
			if arr, ok = groups[node.Name]; ok {
				arr = append(arr, pbBanUser)
			} else {
				arr = []*pbobjs.BanUser{pbBanUser}
			}
			groups[node.Name] = arr
		}
	}
	wg := sync.WaitGroup{}
	for _, banUsers := range groups {
		wg.Add(1)
		users := banUsers
		go func() {
			defer wg.Done()
			bases.AsyncRpcCall(services.ToRpcCtx(ctx, ""), "ban_users", users[0].UserId, &pbobjs.BanUsersReq{
				BanUsers: users,
				IsAdd:    false,
			})
		}()
	}
	wg.Wait()
	tools.SuccessHttpResp(ctx, nil)
}

func QryBanUsers(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	offsetStr := ctx.Query("offset")

	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, ""), "qry_ban_users", fmt.Sprintf("%s%d", services.GetCtxString(ctx, services.CtxKey_AppKey), tools.RandInt(100000)), &pbobjs.QryBanUsersReq{
		Limit:  limit,
		Offset: offsetStr,
	}, func() proto.Message {
		return &pbobjs.QryBanUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	ret := &models.QryBanUsersResp{
		Items: []*models.BanUser{},
	}
	banUsers := resp.(*pbobjs.QryBanUsersResp)
	ret.Offset = banUsers.Offset
	for _, u := range banUsers.Items {
		ret.Items = append(ret.Items, &models.BanUser{
			UserId:      u.UserId,
			CreatedTime: u.CreatedTime,
			EndTime:     u.EndTime,
			ScopeKey:    u.ScopeKey,
			ScopeValue:  u.ScopeValue,
			Ext:         u.Ext,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}

func BlockUser(ctx *gin.Context) {
	var blockReq models.BlockUsersReq
	if err := ctx.BindJSON(&blockReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(blockReq.BlockUserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	for _, blockUserId := range blockReq.BlockUserIds {
		targetId := commonservices.GetConversationId(blockReq.UserId, blockUserId, pbobjs.ChannelType_Private)
		bases.AsyncRpcCall(services.ToRpcCtx(ctx, blockReq.UserId), "block_users", targetId, &pbobjs.BlockUsersReq{
			UserIds: []string{blockUserId},
			IsAdd:   true,
		})
	}
	tools.SuccessHttpResp(ctx, nil)
}

func UnBlockUser(ctx *gin.Context) {
	var blockReq models.BlockUsersReq
	if err := ctx.BindJSON(&blockReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(blockReq.BlockUserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	for _, blockUserId := range blockReq.BlockUserIds {
		targetId := commonservices.GetConversationId(blockReq.UserId, blockUserId, pbobjs.ChannelType_Private)
		bases.AsyncRpcCall(services.ToRpcCtx(ctx, blockReq.UserId), "block_users", targetId, &pbobjs.BlockUsersReq{
			UserIds: []string{blockUserId},
			IsAdd:   false,
		})
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryBlockUsers(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	offsetStr := ctx.Query("offset")
	code, resp, err := bases.SyncRpcCall(services.ToRpcCtx(ctx, userId), "qry_block_users", userId, &pbobjs.QryBlockUsersReq{
		UserId: userId,
		Limit:  limit,
		Offset: offsetStr,
	}, func() proto.Message {
		return &pbobjs.QryBlockUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	ret := &models.QryBlockUsersResp{
		UserId: userId,
		Items:  []*models.BlockUser{},
	}
	blockUsers := resp.(*pbobjs.QryBlockUsersResp)
	ret.Offset = blockUsers.Offset
	for _, u := range blockUsers.Items {
		ret.Items = append(ret.Items, &models.BlockUser{
			BlockUserId: u.BlockUserId,
			CreatedTime: u.CreatedTime,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}
