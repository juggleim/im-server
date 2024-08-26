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
	group := ser.httpServer.Group("/apigateway")
	group.Use(apis.Signature)
	group.POST("/users/register", apis.Register)
	group.POST("/users/update", apis.UpdateUser)
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

	group.GET("/hismsgs/query", apis.QryHisMsgs)
	group.POST("/hismsgs/clean", apis.CleanHisMsgs)
	group.POST("/hismsgs/recall", apis.RecallHisMsgs)
	group.POST("/hismsgs/del", apis.DelHisMsgs)

	group.POST("/groups/add", apis.GroupAddMembers)
	group.GET("/groups/info", apis.QryGroupInfo)
	group.POST("/groups/update", apis.UpdateGroup)
	group.POST("/groups/del", apis.GroupDissolve)
	group.POST("/groups/settings/set", apis.SetGroupSettings)
	group.POST("/groups/settings/get", apis.GetGroupSettings)
	group.POST("/groups/members/add", apis.GroupAddMembers)
	group.POST("/groups/members/del", apis.GroupDelMembers)
	group.GET("/groups/members/query", apis.GroupMembers)
	group.POST("/groups/members/querybyids", apis.GroupMembersByIds)
	group.POST("/groups/groupmute/set", apis.GroupMute)
	group.POST("/groups/groupmembermute/set", apis.GroupMemberMute)
	group.POST("/groups/groupmemberallow/set", apis.GroupMemberAllow)

	group.POST("/chatrooms/create", apis.CreateChatroom)
	group.POST("/chatrooms/destroy", apis.DestroyChatroom)
	group.GET("/chatrooms/info", apis.QryChatroomInfo)
	group.POST("/chatrooms/chrmmute/set", apis.ChrmMute)
	group.POST("/chatrooms/mutemembers/add", apis.AddChrmMuteMembers)
	group.POST("/chatrooms/mutemembers/del", apis.DelChrmMuteMembers)
	group.GET("/chatrooms/mutemembers/query", apis.QryChrmMuteMembers)
	group.POST("/chatrooms/banmembers/add", apis.AddChrmBanMembers)
	group.POST("/chatrooms/banmembers/del", apis.DelChrmBanMembers)
	group.GET("/chatrooms/banmembers/query", apis.QryChrmBanMembers)
	group.POST("/chatrooms/allowmembers/add", apis.AddChrmAllowMembers)
	group.POST("/chatrooms/allowmembers/del", apis.DelChrmAllowMembers)
	group.GET("/chatrooms/allowmembers/query", apis.QryChrmAllowMembers)

	group.GET("/sensitivewords/list", apis.QrySensitiveWords)
	group.POST("/sensitivewords/import", apis.ImportSensitiveWords)
	group.POST("/sensitivewords/add", apis.AddSensitiveWords)
	group.POST("/sensitivewords/del", apis.DeleteSensitiveWords)

	group.POST("/convers/undisturb", apis.UndisturbConvers)
	group.POST("/convers/del", apis.DelConversation)
	group.POST("/convers/add", apis.AddConversation)

	group.GET("/globalconvers/query", apis.QryGlobalConvers)

	group.GET("/usertags/query", apis.QryUserTags)
	group.POST("/usertags/add", apis.AddUserTags)
	group.POST("/usertags/del", apis.DelUserTags)
	group.POST("/usertags/clear", apis.ClearUserTags)

	group.POST("/push", apis.PushWithTags)

	httpPort := configures.Config.ApiGateway.HttpPort
	go ser.httpServer.Run(fmt.Sprintf(":%d", httpPort))
	fmt.Println("start ApiGateway with port :", httpPort)
}

func (ser *ApiGateway) Shutdown() {

}
