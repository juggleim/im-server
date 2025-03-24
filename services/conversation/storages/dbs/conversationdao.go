package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/conversation/storages/models"
	"strings"
)

type ConversationDao struct {
	ID          int64  `gorm:"primary_key"`
	AppKey      string `gorm:"app_key"`
	UserId      string `gorm:"user_id"`
	TargetId    string `gorm:"target_id"`
	ChannelType int    `gorm:"channel_type"`

	SortTime int64 `gorm:"sort_time"`
	SyncTime int64 `gorm:"sync_time"`

	LatestMsgId          string `gorm:"latest_msg_id"`
	LatestMsg            []byte `gorm:"latest_msg"`
	LatestUnreadMsgIndex int64  `gorm:"latest_unread_msg_index"`

	LatestReadMsgIndex int64  `gorm:"latest_read_msg_index"`
	LatestReadMsgId    string `gorm:"latest_read_msg_id"`
	LatestReadMsgTime  int64  `gorm:"latest_read_msg_time"`

	IsTop          int   `gorm:"is_top"`
	TopUpdatedTime int64 `gorm:"top_updated_time"`
	UndisturbType  int32 `gorm:"undisturb_type"`

	UnreadTag  int    `gorm:"unread_tag"`
	ConverExts []byte `gorm:"conver_exts"`
	IsDeleted  int    `gorm:"is_deleted"`
}

func (conver *ConversationDao) TableName() string {
	return "conversations"
}

func (conver *ConversationDao) FindOne(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*models.Conversation, error) {
	var item ConversationDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Conversation{
		UserId:               item.UserId,
		TargetId:             item.TargetId,
		SortTime:             item.SortTime,
		ChannelType:          pbobjs.ChannelType(item.ChannelType),
		AppKey:               item.AppKey,
		LatestMsgId:          item.LatestMsgId,
		LatestMsg:            item.LatestMsg,
		LatestUnreadMsgIndex: item.LatestUnreadMsgIndex,

		LatestReadMsgIndex: item.LatestReadMsgIndex,
		LatestReadMsgId:    item.LatestReadMsgId,
		LatestReadMsgTime:  item.LatestReadMsgTime,
		IsTop:              item.IsTop,
		TopUpdatedTime:     item.TopUpdatedTime,
		UndisturbType:      item.UndisturbType,
		ConverExts:         parseConverExts(item.ConverExts),
	}, nil
}

func (conver *ConversationDao) UpsertConversation(item models.Conversation) error {
	var err error
	if item.SortTime > 0 {
		if item.LatestUnreadMsgIndex > 0 {
			err = dbcommons.GetDb().Exec("INSERT INTO conversations (app_key, user_id, target_id, channel_type, sort_time, latest_msg_id, latest_msg, latest_unread_msg_index, sync_time)VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE sort_time=?, latest_msg_id=?, latest_msg=?, latest_unread_msg_index=?, is_deleted=0, sync_time=?",
				item.AppKey, item.UserId, item.TargetId, item.ChannelType, item.SortTime, item.LatestMsgId, item.LatestMsg, item.LatestUnreadMsgIndex, item.SyncTime, item.SortTime, item.LatestMsgId, item.LatestMsg, item.LatestUnreadMsgIndex, item.SyncTime).Error
		} else {
			err = dbcommons.GetDb().Exec("INSERT INTO conversations (app_key, user_id, target_id, channel_type, sort_time, latest_msg_id, latest_msg, sync_time)VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE sort_time=?, latest_msg_id=?, latest_msg=?, is_deleted=0, sync_time=?",
				item.AppKey, item.UserId, item.TargetId, item.ChannelType, item.SortTime, item.LatestMsgId, item.LatestMsg, item.SyncTime, item.SortTime, item.LatestMsgId, item.LatestMsg, item.SyncTime).Error
		}
	} else {
		if item.LatestUnreadMsgIndex > 0 {
			err = dbcommons.GetDb().Exec("INSERT INTO conversations (app_key, user_id, target_id, channel_type, latest_msg_id, latest_msg, latest_unread_msg_index, sync_time)VALUES(?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE latest_msg_id=?, latest_msg=?, latest_unread_msg_index=?, is_deleted=0, sync_time=?",
				item.AppKey, item.UserId, item.TargetId, item.ChannelType, item.LatestMsgId, item.LatestMsg, item.LatestUnreadMsgIndex, item.SyncTime, item.LatestMsgId, item.LatestMsg, item.LatestUnreadMsgIndex, item.SyncTime).Error
		} else {
			err = dbcommons.GetDb().Exec("INSERT INTO conversations (app_key, user_id, target_id, channel_type, latest_msg_id, latest_msg, sync_time)VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE latest_msg_id=?, latest_msg=?, is_deleted=0, sync_time=?",
				item.AppKey, item.UserId, item.TargetId, item.ChannelType, item.LatestMsgId, item.LatestMsg, item.SyncTime, item.LatestMsgId, item.LatestMsg, item.SyncTime).Error
		}
	}
	return err
}

func (conver *ConversationDao) Upsert(item models.Conversation) error {
	var sqlBuilder strings.Builder
	params := []interface{}{}
	sqlBuilder.WriteString("INSERT INTO ")
	sqlBuilder.WriteString(conver.TableName())
	sqlBuilder.WriteString(" (app_key,user_id,target_id,channel_type,sort_time,sync_time,latest_msg_id,latest_msg,latest_unread_msg_index,latest_read_msg_index,latest_read_msg_id,latest_read_msg_time,is_top,top_updated_time,undisturb_type,unread_tag,is_deleted,conver_exts) ")
	sqlBuilder.WriteString("VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE ")
	params = append(params,
		item.AppKey, item.UserId, item.TargetId, item.ChannelType,
		item.SortTime, item.SyncTime,
		item.LatestMsgId, item.LatestMsg, item.LatestUnreadMsgIndex,
		item.LatestReadMsgIndex, item.LatestReadMsgId, item.LatestReadMsgTime,
		item.IsTop, item.TopUpdatedTime, item.UndisturbType,
		item.UnreadTag,
		item.IsDeleted,
		converExts2Bs(item.ConverExts),
	)
	sqlBuilder.WriteString("sort_time=VALUES(sort_time),sync_time=VALUES(sync_time),latest_msg_id=VALUES(latest_msg_id),latest_msg=VALUES(latest_msg),latest_unread_msg_index=VALUES(latest_unread_msg_index),")
	sqlBuilder.WriteString("latest_read_msg_index=VALUES(latest_read_msg_index),latest_read_msg_id=VALUES(latest_read_msg_id),latest_read_msg_time=VALUES(latest_read_msg_time),")
	sqlBuilder.WriteString("is_top=VALUES(is_top),top_updated_time=VALUES(top_updated_time),undisturb_type=VALUES(undisturb_type),")
	sqlBuilder.WriteString("unread_tag=VALUES(unread_tag),is_deleted=VALUES(is_deleted),conver_exts=VALUES(conver_exts)")
	return dbcommons.GetDb().Exec(sqlBuilder.String(), params...).Error
}

func (conver *ConversationDao) QryConvers(appkey, userId string, startTime int64, count int32) ([]*models.Conversation, error) {
	var items []*ConversationDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and sync_time<?", appkey, userId, startTime).Order("sync_time desc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.Conversation{}, err
	}
	conversations := []*models.Conversation{}
	for _, item := range items {
		conversations = append(conversations, &models.Conversation{
			UserId:               item.UserId,
			TargetId:             item.TargetId,
			SortTime:             item.SortTime,
			ChannelType:          pbobjs.ChannelType(item.ChannelType),
			LatestMsgId:          item.LatestMsgId,
			LatestMsg:            item.LatestMsg,
			LatestUnreadMsgIndex: item.LatestUnreadMsgIndex,
			AppKey:               item.AppKey,
			LatestReadMsgIndex:   item.LatestReadMsgIndex,
			LatestReadMsgId:      item.LatestReadMsgId,
			LatestReadMsgTime:    item.LatestReadMsgTime,
			SyncTime:             item.SyncTime,
			UndisturbType:        item.UndisturbType,
			IsTop:                item.IsTop,
			TopUpdatedTime:       item.TopUpdatedTime,
			IsDeleted:            item.IsDeleted,
			UnreadTag:            item.UnreadTag,
			ConverExts:           parseConverExts(item.ConverExts),
		})
	}
	return conversations, nil
}

func (conver *ConversationDao) ClearTotalUnreadCount(appkey, userId string) error {
	return dbcommons.GetDb().Exec("UPDATE conversations set latest_read_msg_index=latest_unread_msg_index, latest_read_msg_time =(UNIX_TIMESTAMP(NOW(3)) * 1000), unread_tag=0 where app_key=? and user_id=?", appkey, userId).Error
}

func parseConverExts(bs []byte) *pbobjs.ConverExts {
	if len(bs) > 0 {
		var tags pbobjs.ConverExts
		err := tools.PbUnMarshal(bs, &tags)
		if err == nil && len(tags.ConverTags) > 0 {
			return &tags
		}
	}
	return nil
}

func converExts2Bs(exts *pbobjs.ConverExts) []byte {
	if exts != nil {
		bs, _ := tools.PbMarshal(exts)
		return bs
	}
	return []byte{}
}
