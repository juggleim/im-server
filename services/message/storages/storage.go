package storages

import (
	"im-server/commons/configures"
	"im-server/services/message/storages/dbs"
	"im-server/services/message/storages/models"
	"im-server/services/message/storages/mongodbs"
)

func NewInboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.InboxMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.InboxMsgDao{}
	default:
		return &dbs.InboxMsgDao{}
	}
}

func NewSendboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.SendboxMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.SendboxMsgDao{}
	default:
		return &dbs.SendboxMsgDao{}
	}
}

func NewCmdInboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.CmdInboxMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.CmdInboxMsgDao{}
	default:
		return &dbs.CmdInboxMsgDao{}
	}
}

func NewCmdSendboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.CmdSendboxMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.CmdSendboxMsgDao{}
	default:
		return &dbs.CmdSendboxMsgDao{}
	}
}

func NewBrdInboxMsgStorage() models.IBroadcastMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.BrdInboxMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.BrdInboxMsgDao{}
	default:
		return &dbs.BrdInboxMsgDao{}
	}
}

func NewFriendRelStorage() models.IFriendRelStorage {
	return &dbs.FriendRelDao{}
}
