package storages

import (
	"im-server/services/chatroom/storages/dbs"
	"im-server/services/chatroom/storages/models"
)

func NewChatroomStorage() models.IChatroomStorage {
	return &dbs.ChatroomDao{}
}

func NewChatroomMemberStorage() models.IChatroomMemberStorage {
	return &dbs.ChatroomMemberDao{}
}

func NewChatroomExtStorage() models.IChatroomExtStorage {
	return &dbs.ChatroomExtDao{}
}

func NewChatroomBanUserStorage() models.IChatroomBanUserStorage {
	return &dbs.ChatroomBanUserDao{}
}
