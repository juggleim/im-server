package storages

import (
	"testing"

	"im-server/commons/configures"
	"im-server/services/message/storages/dbs"
	"im-server/services/message/storages/mongodbs"
)

func TestMessageStorageUsesSplitTablesOnlyForMySQL(t *testing.T) {
	original := configures.Config.MsgStoreEngine
	t.Cleanup(func() { configures.Config.MsgStoreEngine = original })

	configures.Config.MsgStoreEngine = configures.MsgStoreEngine_MySQL
	if _, ok := NewInboxMsgStorage().(*dbs.InboxMsgSplitDao); !ok {
		t.Fatal("MySQL inbox storage should use the split table")
	}
	if _, ok := NewSendboxMsgStorage().(*dbs.SendboxMsgSplitDao); !ok {
		t.Fatal("MySQL sendbox storage should use the split table")
	}
	if _, ok := NewCmdInboxMsgStorage().(*dbs.CmdInboxMsgSplitDao); !ok {
		t.Fatal("MySQL cmd inbox storage should use the split table")
	}
	if _, ok := NewCmdSendboxMsgStorage().(*dbs.CmdSendboxMsgSplitDao); !ok {
		t.Fatal("MySQL cmd sendbox storage should use the split table")
	}

	configures.Config.MsgStoreEngine = configures.MsgStoreEngine_Mongo
	if _, ok := NewInboxMsgStorage().(*mongodbs.InboxMsgDao); !ok {
		t.Fatal("Mongo inbox storage should keep the original implementation")
	}
	if _, ok := NewSendboxMsgStorage().(*mongodbs.SendboxMsgDao); !ok {
		t.Fatal("Mongo sendbox storage should keep the original implementation")
	}
}
