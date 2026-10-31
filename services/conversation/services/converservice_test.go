package services

import (
	"reflect"
	"testing"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
)

func TestCollectPrivateTargetUserInfoIds(t *testing.T) {
	userId := "current-user"
	convers := []*models.Conversation{
		{TargetId: "carried", ChannelType: pbobjs.ChannelType_Private, LatestMsgId: "carried-msg"},
		{TargetId: "sender", ChannelType: pbobjs.ChannelType_Private, LatestMsgId: "sender-msg"},
		{TargetId: "nil-info", ChannelType: pbobjs.ChannelType_Private, LatestMsgId: "nil-info-msg"},
		{TargetId: "no-msg", ChannelType: pbobjs.ChannelType_Private, LatestMsgId: "missing-msg"},
		{TargetId: "deleted", ChannelType: pbobjs.ChannelType_Private, LatestMsgId: "deleted-msg", IsDeleted: 1},
		{TargetId: "group", ChannelType: pbobjs.ChannelType_Group},
		{TargetId: "sender", ChannelType: pbobjs.ChannelType_Private},
		{TargetId: "", ChannelType: pbobjs.ChannelType_Private},
		nil,
	}
	priMsgs := map[string]*pbobjs.DownMsg{
		"carried-msg": {
			SenderId:       "other-user",
			TargetUserInfo: &pbobjs.UserInfo{UserId: "carried"},
		},
		"sender-msg": {
			SenderId:       userId,
			TargetUserInfo: &pbobjs.UserInfo{UserId: "sender"},
		},
		"nil-info-msg": {
			SenderId: "other-user",
		},
	}

	got := collectPrivateTargetUserInfoIds(userId, convers, priMsgs)
	want := []string{"sender", "nil-info", "no-msg", "deleted"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("collectPrivateTargetUserInfoIds() = %v, want %v", got, want)
	}
}
