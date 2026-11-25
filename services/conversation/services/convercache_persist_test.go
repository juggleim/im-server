package services

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"testing"
)

func TestMergePendingConversationPreservesSettingsOnMessageUpdate(t *testing.T) {
	setting := &models.Conversation{
		AppKey:         "app",
		UserId:         "user",
		TargetId:       "group",
		ChannelType:    pbobjs.ChannelType_Group,
		SyncTime:       100,
		LatestMsgId:    "old",
		UndisturbType:  1,
		IsTop:          1,
		TopUpdatedTime: 90,
		UnreadTag:      1,
		IsDeleted:      1,
		ConverExts:     &pbobjs.ConverExts{ConverTags: map[string]bool{"tag": true}},
	}
	pending := mergePendingConversation(nil, setting, converPersistUndisturb|converPersistTop|converPersistUnreadTag|converPersistExts, true)

	message := &models.Conversation{
		AppKey:               "app",
		UserId:               "user",
		TargetId:             "group",
		ChannelType:          pbobjs.ChannelType_Group,
		SyncTime:             200,
		SortTime:             200,
		LatestMsgId:          "new",
		LatestUnreadMsgIndex: 8,
		UndisturbType:        0,
		IsTop:                0,
		UnreadTag:            0,
		IsDeleted:            0,
	}
	pending = mergePendingConversation(pending, message, converPersistMessage|converPersistDeleted, true)
	got := pending.(*pendingConversation).Conversation

	if got.LatestMsgId != "new" || got.SyncTime != 200 || got.SortTime != 200 || got.LatestUnreadMsgIndex != 8 || got.IsDeleted != 0 {
		t.Fatalf("message fields were not merged: %+v", got)
	}
	if got.UndisturbType != 1 || got.IsTop != 1 || got.TopUpdatedTime != 90 || got.UnreadTag != 1 {
		t.Fatalf("message update overwrote setting fields: %+v", got)
	}
	if got.ConverExts == nil || !got.ConverExts.ConverTags["tag"] {
		t.Fatalf("message update overwrote conversation extensions: %+v", got.ConverExts)
	}
}

func TestMergePendingConversationKeepsLaterUndisturbUpdate(t *testing.T) {
	message := &models.Conversation{
		AppKey:        "app",
		UserId:        "user",
		TargetId:      "group",
		ChannelType:   pbobjs.ChannelType_Group,
		SyncTime:      100,
		LatestMsgId:   "first",
		UndisturbType: 0,
	}
	pending := mergePendingConversation(nil, message, converPersistMessage, true)

	setting := cloneConversation(message)
	setting.UndisturbType = 1
	pending = mergePendingConversation(pending, setting, converPersistUndisturb, true)

	staleMessage := cloneConversation(message)
	staleMessage.SyncTime = 200
	staleMessage.LatestMsgId = "second"
	staleMessage.UndisturbType = 0
	pending = mergePendingConversation(pending, staleMessage, converPersistMessage, true)
	got := pending.(*pendingConversation).Conversation

	if got.UndisturbType != 1 {
		t.Fatalf("later message lost pending undisturb update: %+v", got)
	}
	if got.LatestMsgId != "second" || got.SyncTime != 200 {
		t.Fatalf("later message was not merged: %+v", got)
	}
}

func TestMergePendingConversationClonesExtensions(t *testing.T) {
	item := &models.Conversation{
		AppKey:      "app",
		UserId:      "user",
		TargetId:    "group",
		ChannelType: pbobjs.ChannelType_Group,
		ConverExts:  &pbobjs.ConverExts{ConverTags: map[string]bool{"first": true}},
	}
	pending := mergePendingConversation(nil, item, converPersistExts, true).(*pendingConversation)
	item.ConverExts.ConverTags["after_enqueue"] = true
	if pending.Conversation.ConverExts.ConverTags["after_enqueue"] {
		t.Fatal("pending extensions share the source map")
	}

	updated := cloneConversation(item)
	updated.ConverExts = &pbobjs.ConverExts{ConverTags: map[string]bool{"updated": true}}
	mergePendingConversation(pending, updated, converPersistExts, true)
	updated.ConverExts.ConverTags["after_merge"] = true
	if pending.Conversation.ConverExts.ConverTags["after_merge"] {
		t.Fatal("merged extensions share the source map")
	}
}

func TestMergePendingConversationDoesNotCreateForCacheOnlyMerge(t *testing.T) {
	item := &models.Conversation{TargetId: "group"}
	if got := mergePendingConversation(nil, item, converPersistRead, false); got != nil {
		t.Fatalf("cache-only merge unexpectedly created pending state: %#v", got)
	}
}
