package storages

import (
	"im-server/commons/configures"
	"im-server/services/historymsg/storages/dbs"
	"im-server/services/historymsg/storages/models"
	"im-server/services/historymsg/storages/mongodbs"
)

func NewPrivateHisMsgStorage() models.IPrivateHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.PrivateHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.PrivateHisMsgDao{}
	default:
		return &dbs.PrivateHisMsgDao{}
	}
}

func NewGroupHisMsgStorage() models.IGroupHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.GroupHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.GroupHisMsgDao{}
	default:
		return &dbs.GroupHisMsgDao{}
	}
}

func NewSystemHisMsgStorage() models.ISystemHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.SystemHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.SystemHisMsgDao{}
	default:
		return &dbs.SystemHisMsgDao{}
	}
}

func NewBrdCastHisMsgStorage() models.IBrdCastHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.BrdCastHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.BrdCastHisMsgDao{}
	default:
		return &dbs.BrdCastHisMsgDao{}
	}
}

func NewGrpCastHisMsgStorage() models.IGrpCastHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.GrpCastHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.GrpCastHisMsgDao{}
	default:
		return &dbs.GrpCastHisMsgDao{}
	}
}

func NewMergedMsgStorage() models.IMergedMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.MergedMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.MergedMsgDao{}
	default:
		return &dbs.MergedMsgDao{}
	}
}

func NewHisMsgConverCleanTimeStorage() models.IHisMsgConverCleanTimeStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.HisMsgConverCleanTimeDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.HisMsgConverCleanTimeDao{}
	default:
		return &dbs.HisMsgConverCleanTimeDao{}
	}
}

func NewHisMsgUserCleanTimeStorage() models.IHisMsgUserCleanTimeStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.HisMsgUserCleanTimeDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.HisMsgUserCleanTimeDao{}
	default:
		return &dbs.HisMsgUserCleanTimeDao{}
	}
}

func NewGroupDelHisMsgStorage() models.IGroupDelHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.GroupDelHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.GroupDelHisMsgDao{}
	default:
		return &dbs.GroupDelHisMsgDao{}
	}
}

func NewPrivateDelHisMsgStorage() models.IPrivateDelHisMsgStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.PrivateDelHisMsgDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.PrivateDelHisMsgDao{}
	default:
		return &dbs.PrivateDelHisMsgDao{}
	}
}

func NewReadInfoStorage() models.IReadInfoStorage {
	return &dbs.ReadInfoDao{}
}

func NewMsgExtStorage() models.IMsgExtStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.MsgExtDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.MsgExtDao{}
	default:
		return &dbs.MsgExtDao{}
	}
}

func NewMsgExSetStorage() models.IMsgExSetStorage {
	switch configures.Config.MsgStoreEngine {
	case configures.MsgStoreEngine_MySQL:
		return &dbs.MsgExSetDao{}
	case configures.MsgStoreEngine_Mongo:
		return &mongodbs.MsgExSetDao{}
	default:
		return &dbs.MsgExSetDao{}
	}
}
