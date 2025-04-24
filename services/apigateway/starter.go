package apigateway

import (
	"fmt"
	"net/http"

	"im-server/commons/configures"
	"im-server/commons/gmicro"
	"im-server/services/apigateway/apis"

	"github.com/gin-gonic/gin"
)

type ApiGateway struct {
	httpServer *gin.Engine
}

func (ser *ApiGateway) RegisterActors(register gmicro.IActorRegister) {

}

func (ser *ApiGateway) Startup(args map[string]interface{}) {
	ser.httpServer = gin.Default()
	ser.httpServer.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "ok")
	})
	ser.httpServer.HEAD("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, nil)
	})
	group := ser.httpServer.Group("/apigateway")
	group.Use(apis.Signature)
	group.POST("/users/register", apis.Register)
	group.POST("/users/update", apis.UpdateUser)
	group.POST("/users/settings/set", apis.SetUserSettings)
	group.GET("/users/settings/get", apis.GetUserSettings)
	group.GET("/users/info", apis.QryUserInfo)
	group.POST("/users/kick", apis.KickUsers)
	group.POST("/users/onlinestatus/query", apis.QryUserOnlineStatus)
	group.POST("/users/banusers/ban", apis.UserBan)
	group.POST("/users/banusers/unban", apis.UserUnBan)
	group.GET("/users/banusers/query", apis.QryBanUsers)
	group.POST("/users/blockusers/block", apis.BlockUser)
	group.POST("/users/blockusers/unblock", apis.UnBlockUser)
	group.GET("/users/blockusers/query", apis.QryBlockUsers)

	group.POST("/bots/add", apis.AddBot)

	group.POST("/messages/private/send", apis.SendPrivateMsg)
	group.POST("/messages/system/send", apis.SendSystemMsg)
	group.POST("/messages/group/send", apis.SendGroupMsg)
	group.POST("/messages/groupcast/send", apis.SendGroupCastMsg)
	group.POST("/messages/broadcast/send", apis.SendBroadCastMsg)
	group.POST("/messages/markread", apis.MarkRead)
	group.POST("/messages/private/stream/create", apis.CreatePrivateStreamMsg)
	group.POST("/messages/private/stream/append", apis.AppendPrivateStreamMsg)
	group.POST("/messages/private/stream/complete", apis.CompletePrivateStreamMsg)

	group.GET("/hismsgs/query", apis.QryHisMsgs)
	group.POST("/hismsgs/clean", apis.CleanHisMsgs)
	group.POST("/hismsgs/recall", apis.RecallHisMsgs)
	group.POST("/hismsgs/del", apis.DelHisMsgs)
	group.POST("/hismsgs/modify", apis.ModifyHisMsg)
	group.POST("/hismsgs/import", apis.ImportHisMsg)

	group.POST("/private/globalmutemembers/add", apis.AddPrivateGlobalMuteMembers)
	group.POST("/private/globalmutemembers/del", apis.DelPrivateGlobalMuteMembers)
	group.GET("/private/globalmutemembers/query", apis.QryPrivateGlobalMuteMembers)

	group.POST("/groups/add", apis.GroupAddMembers)
	group.GET("/groups/info", apis.QryGroupInfo)
	group.POST("/groups/update", apis.UpdateGroup)
	group.POST("/groups/del", apis.GroupDissolve)
	group.POST("/groups/settings/set", apis.SetGroupSettings)
	group.GET("/groups/settings/get", apis.GetGroupSettings)
	group.POST("/groups/members/add", apis.GroupAddMembers)
	group.POST("/groups/members/del", apis.GroupDelMembers)
	group.GET("/groups/members/query", apis.GroupMembers)
	group.POST("/groups/members/update", apis.GroupMemberUpdate)
	group.POST("/groups/members/querybyids", apis.GroupMembersByIds)
	group.POST("/groups/groupmute/set", apis.GroupMute)
	group.POST("/groups/groupmembermute/set", apis.GroupMemberMute)
	group.POST("/groups/groupmemberallow/set", apis.GroupMemberAllow)

	group.GET("/sensitivewords/list", apis.QrySensitiveWords)
	group.POST("/sensitivewords/import", apis.ImportSensitiveWords)
	group.POST("/sensitivewords/add", apis.AddSensitiveWords)
	group.POST("/sensitivewords/del", apis.DeleteSensitiveWords)

	group.POST("/convers/undisturb", apis.UndisturbConvers)
	group.POST("/convers/del", apis.DelConversation)
	group.POST("/convers/add", apis.AddConversation)
	group.POST("/convers/clearunread", apis.ClearConverUnread)
	group.POST("/convers/top", apis.TopConversations)

	group.GET("/globalconvers/query", apis.QryGlobalConvers)

	group.GET("/usertags/query", apis.QryUserTags)
	group.POST("/usertags/add", apis.AddUserTags)
	group.POST("/usertags/del", apis.DelUserTags)
	group.POST("/usertags/clear", apis.ClearUserTags)

	group.POST("/push", apis.PushWithTags)

	group.POST("/friends/add", apis.AddFriends)
	group.POST("/friends/del", apis.DelFriends)

	httpPort := configures.Config.ApiGateway.HttpPort
	go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
	fmt.Println("Start apigateway with port:", httpPort)
}

func (ser *ApiGateway) Shutdown() {

}
