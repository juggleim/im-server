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
		return &dbs.InboxMsgSplitDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.InboxMsgDao{}
	default:
		return &dbs.InboxMsgSplitDao{}
	}
}

func NewSendboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.SendboxMsgSplitDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.SendboxMsgDao{}
	default:
		return &dbs.SendboxMsgSplitDao{}
	}
}

func NewCmdInboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.CmdInboxMsgSplitDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.CmdInboxMsgDao{}
	default:
		return &dbs.CmdInboxMsgSplitDao{}
	}
}

func NewCmdSendboxMsgStorage() models.IMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.CmdSendboxMsgSplitDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.CmdSendboxMsgDao{}
	default:
		return &dbs.CmdSendboxMsgSplitDao{}
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
