package storages

import (
	"im-server/services/conversation/storages/dbs"
	"im-server/services/conversation/storages/models"
)

func NewConversationStorage() models.IConversationStorage {
	return &dbs.ConversationDao{}
}

func NewMentionMsgStorage() models.IMentionMsgStorage {
	return &dbs.MentionMsgDao{}
}

func NewGlobalConversationStorage() models.IGlobalConverStorage {
	return &dbs.GlobalConverDao{}
}

func NewUserConverTagStorage() models.IUserConverTagStorage {
	return &dbs.UserConverTagDao{}
}

func NewConverTagRelStorage() models.IConverTagRelStorage {
	return &dbs.ConverTagRelDao{}
}
