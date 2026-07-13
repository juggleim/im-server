package services

import (
	"testing"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"

	"google.golang.org/protobuf/proto"
)

func TestTotalUnreadCountFilter(t *testing.T) {
	privateConver := &models.Conversation{
		TargetId:             "u1",
		ChannelType:          pbobjs.ChannelType_Private,
		LatestUnreadMsgIndex: 5,
		LatestReadMsgIndex:   2,
	}
	groupConver := &models.Conversation{
		TargetId:             "g1",
		ChannelType:          pbobjs.ChannelType_Group,
		LatestUnreadMsgIndex: 4,
		LatestReadMsgIndex:   1,
		ConverExts: &pbobjs.ConverExts{
			ConverTags:       map[string]bool{"team": true},
			GlobalConverTags: map[string]bool{"global_team": true},
		},
	}
	chatroomConver := &models.Conversation{
		TargetId:             "c1",
		ChannelType:          pbobjs.ChannelType_Chatroom,
		LatestUnreadMsgIndex: 2,
		LatestReadMsgIndex:   2,
		UnreadTag:            1,
		ConverExts: &pbobjs.ConverExts{
			ConverTags: map[string]bool{"muted": true},
		},
	}
	deletedConver := &models.Conversation{
		TargetId:             "d1",
		ChannelType:          pbobjs.ChannelType_Private,
		LatestUnreadMsgIndex: 100,
		IsDeleted:            1,
	}
	userConvers := &UserConversations{
		Appkey:        "app",
		UserId:        "user",
		ConverItemMap: make(map[string]*models.Conversation),
	}
	for _, conver := range []*models.Conversation{privateConver, groupConver, chatroomConver, deletedConver} {
		userConvers.ConverItemMap[getConverItemKey(conver.TargetId, conver.SubChannel, conver.ChannelType)] = conver
	}

	tests := []struct {
		name   string
		filter *pbobjs.ConverFilter
		want   int64
	}{
		{
			name: "nil filter counts all undeleted conversations",
			want: 7,
		},
		{
			name: "channel types keep only listed channel types",
			filter: &pbobjs.ConverFilter{
				ChannelTypes: []pbobjs.ChannelType{pbobjs.ChannelType_Private, pbobjs.ChannelType_Chatroom},
			},
			want: 4,
		},
		{
			name: "include conversations keep only included conversations",
			filter: &pbobjs.ConverFilter{
				IncludeConvers: []*pbobjs.SimpleConversation{
					{TargetId: "g1", ChannelType: pbobjs.ChannelType_Group},
				},
			},
			want: 3,
		},
		{
			name: "include conversations override exclude conversations",
			filter: &pbobjs.ConverFilter{
				IncludeConvers: []*pbobjs.SimpleConversation{
					{TargetId: "g1", ChannelType: pbobjs.ChannelType_Group},
				},
				ExcludeConvers: []*pbobjs.SimpleConversation{
					{TargetId: "g1", ChannelType: pbobjs.ChannelType_Group},
				},
			},
			want: 3,
		},
		{
			name: "exclude conversations apply when include conversations is empty",
			filter: &pbobjs.ConverFilter{
				ExcludeConvers: []*pbobjs.SimpleConversation{
					{TargetId: "g1", ChannelType: pbobjs.ChannelType_Group},
				},
			},
			want: 4,
		},
		{
			name: "tag keeps conversations containing tag",
			filter: &pbobjs.ConverFilter{
				Tag: "team",
			},
			want: 3,
		},
		{
			name: "tag keeps conversations containing global tag",
			filter: &pbobjs.ConverFilter{
				Tag: "global_team",
			},
			want: 3,
		},
		{
			name: "tag overrides except tag",
			filter: &pbobjs.ConverFilter{
				Tag:       "team",
				ExceptTag: "team",
			},
			want: 3,
		},
		{
			name: "except tag applies when tag is empty",
			filter: &pbobjs.ConverFilter{
				ExceptTag: "muted",
			},
			want: 6,
		},
		{
			name: "except global tag applies when tag is empty",
			filter: &pbobjs.ConverFilter{
				ExceptTag: "global_team",
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userConvers.TotalUnreadCount(tt.filter); got != tt.want {
				t.Fatalf("TotalUnreadCount() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestConverFilterExceptTagProtoRoundTrip(t *testing.T) {
	bs, err := proto.Marshal(&pbobjs.ConverFilter{ExceptTag: "muted"})
	if err != nil {
		t.Fatalf("proto.Marshal() error = %v", err)
	}

	var filter pbobjs.ConverFilter
	if err := proto.Unmarshal(bs, &filter); err != nil {
		t.Fatalf("proto.Unmarshal() error = %v", err)
	}
	if filter.ExceptTag != "muted" {
		t.Fatalf("ExceptTag = %q, want muted", filter.ExceptTag)
	}
}
