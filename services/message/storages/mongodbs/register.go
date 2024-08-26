package mongodbs

import "im-server/commons/mongocommons"

func RegistCollections() {
	mongocommons.Register(&BrdInboxMsgDao{})
	mongocommons.Register(&CmdInboxMsgDao{})
	mongocommons.Register(&CmdSendboxMsgDao{})
	mongocommons.Register(&InboxMsgDao{})
	mongocommons.Register(&SendboxMsgDao{})
}
