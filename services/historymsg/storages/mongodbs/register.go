package mongodbs

import "im-server/commons/mongocommons"

func RegistCollections() {
	mongocommons.Register(&BrdCastHisMsgDao{})
	mongocommons.Register(&HisMsgConverCleanTimeDao{})
	// mongocommons.Register(&GroupDelHisMsgDao{})
	mongocommons.Register(&GroupHisMsgDao{})
	mongocommons.Register(&GrpCastHisMsgDao{})
	mongocommons.Register(&MergedMsgDao{})
	mongocommons.Register(&MsgExtDao{})
	// mongocommons.Register(&PrivateDelHisMsgDao{})
	mongocommons.Register(&PrivateHisMsgDao{})
	mongocommons.Register(&ReadInfoDao{})
	mongocommons.Register(&SystemHisMsgDao{})
	mongocommons.Register(&HisMsgUserCleanTimeDao{})
}
