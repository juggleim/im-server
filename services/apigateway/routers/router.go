package routers

import (
	"im-server/services/apigateway/apis"

	"github.com/gin-gonic/gin"
)

func Route(eng *gin.Engine, prefix string) {
	group := eng.Group("/" + prefix)
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
	group.POST("/messages/private/stream/send", apis.SendPrivateStreamMsg)

	group.GET("/hismsgs/query", apis.QryHisMsgs)
	group.POST("/hismsgs/clean", apis.CleanHisMsgs)
	group.POST("/hismsgs/recall", apis.RecallHisMsgs)
	group.POST("/hismsgs/del", apis.DelHisMsgs)
	group.POST("/hismsgs/modify", apis.ModifyHisMsg)
	group.POST("/hismsgs/import", apis.ImportHisMsg)

	group.POST("/private/globalmutemembers/add", apis.AddPrivateGlobalMuteMembers)
	group.POST("/private/globalmutemembers/del", apis.DelPrivateGlobalMuteMembers)
	group.GET("/private/globalmutemembers/query", apis.QryPrivateGlobalMuteMembers)
	group.POST("/group/globalmutemembers/add", apis.AddGroupGlobalMuteMembers)
	group.POST("/group/globalmutemembers/del", apis.DelGroupGlobalMuteMembers)
	group.GET("/group/globalmutemembers/query", apis.QryGroupGlobalMuteMembers)

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
	group.GET("/convers/query", apis.QryConvers)

	group.GET("/globalconvers/query", apis.QryGlobalConvers)

	group.GET("/usertags/query", apis.QryUserTags)
	group.POST("/usertags/add", apis.AddUserTags)
	group.POST("/usertags/del", apis.DelUserTags)
	group.POST("/usertags/clear", apis.ClearUserTags)

	group.POST("/push", apis.PushWithTags)

	group.POST("/friends/add", apis.AddFriends)
	group.POST("/friends/del", apis.DelFriends)
	group.GET("/friends/query", apis.QryFriends)
	group.POST("/friends/setdisplayname", apis.SetFriendDisplayName)
}
