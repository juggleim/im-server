package dbs

import (
	"fmt"
	"time"

	"im-server/commons/dbcommons"
	"im-server/services/message/storages/models"
)

const (
	cmdMsgSplitPartitionIntervalMillis = int64(24 * time.Hour / time.Millisecond)
)

// partition_day is the UTC-day start derived from send_time. It is both the
// physical partition key and part of the uniq_tag key, so MySQL can perform a
// single-statement upsert without a cross-node advisory lock.
func upsertCmdMsgSplit(tableName string, item models.Msg) error {
	if tableName != "cmdinboxmsgs_split" && tableName != "cmdsendboxmsgs_split" {
		return fmt.Errorf("invalid cmd split table %q", tableName)
	}
	partitionStart := alignCmdMsgSplitPartitionStart(item.SendTime)
	statement := fmt.Sprintf("INSERT INTO `%s` (user_id,send_time,partition_day,msg_id,channel_type,msg_body,app_key,target_id,msg_type,uniq_tag) VALUES (?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE send_time=VALUES(send_time), msg_id=VALUES(msg_id), msg_body=VALUES(msg_body)", tableName)
	return dbcommons.GetDb().Exec(statement, item.UserId, item.SendTime, partitionStart, item.MsgId, item.ChannelType, item.MsgBody, item.AppKey, item.TargetId, item.MsgType, item.UniqTag).Error
}

func alignCmdMsgSplitPartitionStart(sendTime int64) int64 {
	return sendTime - sendTime%cmdMsgSplitPartitionIntervalMillis
}
