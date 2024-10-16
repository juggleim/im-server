package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"strings"
)

type ConversationDao struct {
	ID                   int64  `gorm:"primary_key"`
	UserId               string `gorm:"user_id"`
	TargetId             string `gorm:"target_id"`
	SortTime             int64  `gorm:"sort_time"`
	ChannelType          int    `gorm:"channel_type"`
	AppKey               string `gorm:"app_key"`
	LatestMsgId          string `gorm:"latest_msg_id"`
	LatestMsg            []byte `gorm:"latest_msg"`
	LatestUnreadMsgIndex int64  `gorm:"latest_unread_msg_index"`
	LatestReadMsgIndex   int64  `gorm:"latest_read_msg_index"`
	LatestReadMsgId      string `gorm:"latest_read_msg_id"`
	LatestReadMsgTime    int64  `gorm:"latest_read_msg_time"`
	IsDeleted            int    `gorm:"is_deleted"`
	IsTop                int    `gorm:"is_top"`
	TopUpdatedTime       int64  `gorm:"top_updated_time"`
	UndisturbType        int32  `gorm:"undisturb_type"`
	SyncTime             int64  `gorm:"sync_time"`
	UnreadTag            int    `gorm:"unread_tag"`
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
	}, nil
}

func (conver *ConversationDao) BatchFind(appkey string, reqConvers []models.Conversation) ([]*models.Conversation, error) {
	if len(reqConvers) <= 0 {
		return []*models.Conversation{}, nil
	}
	var items []*ConversationDao
	var conditionBuilder strings.Builder
	params := []interface{}{}
	conditionBuilder.WriteString("app_key=? and (")
	params = append(params, appkey)
	l := len(reqConvers)
	for i, req := range reqConvers {
		if i == l-1 {
			conditionBuilder.WriteString("(user_id=? and target_id=? and channel_type=?))")
		} else {
			conditionBuilder.WriteString("(user_id=? and target_id=? and channel_type=?) or ")
		}
		params = append(params, req.UserId)
		params = append(params, req.TargetId)
		params = append(params, req.ChannelType)
	}
	err := dbcommons.GetDb().Where(conditionBuilder.String(), params...).Select("user_id,target_id,channel_type,undisturb_type,latest_unread_msg_index").Find(&items).Error
	if err != nil {
		return []*models.Conversation{}, err
	}
	ret := []*models.Conversation{}
	for _, item := range items {
		ret = append(ret, &models.Conversation{
			UserId:               item.UserId,
			TargetId:             item.TargetId,
			SortTime:             item.SortTime,
			ChannelType:          pbobjs.ChannelType(item.ChannelType),
			LatestMsgId:          item.LatestMsgId,
			LatestMsg:            item.LatestMsg,
			LatestUnreadMsgIndex: item.LatestUnreadMsgIndex,
			LatestReadMsgIndex:   item.LatestReadMsgIndex,
			LatestReadMsgId:      item.LatestReadMsgId,
			LatestReadMsgTime:    item.LatestReadMsgTime,
			IsTop:                item.IsTop,
			TopUpdatedTime:       item.TopUpdatedTime,
			UndisturbType:        item.UndisturbType,
			IsDeleted:            item.IsDeleted,
			AppKey:               item.AppKey,
			SyncTime:             item.SyncTime,
		})
	}
	return ret, nil
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

func (conver *ConversationDao) SyncConversations(appkey, userId string, startTime int64, count int32) ([]*models.Conversation, error) {
	var items []*ConversationDao
	var conditionTail string = ""
	// if startTime == 0 {
	// 	conditionTail = " and is_deleted=0"
	// }
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and sync_time>?"+conditionTail, appkey, userId, startTime).Order("sync_time asc").Limit(count).Find(&items).Error
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
			LatestReadMsgIndex:   item.LatestReadMsgIndex,
			LatestReadMsgId:      item.LatestReadMsgId,
			LatestReadMsgTime:    item.LatestReadMsgTime,
			IsTop:                item.IsTop,
			TopUpdatedTime:       item.TopUpdatedTime,
			UndisturbType:        item.UndisturbType,
			IsDeleted:            item.IsDeleted,
			AppKey:               item.AppKey,
			SyncTime:             item.SyncTime,
		})
	}
	return conversations, nil
}

func (conver *ConversationDao) QryConversations(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int32, isPositiveOrder bool, tag string) ([]*models.Conversation, error) {
	var sqlBuilder strings.Builder
	var conditionBuilder strings.Builder
	var items []*ConversationDao
	//sql := fmt.Sprintf("select %s.* from %s", conver.TableName(), conver.TableName())
	sqlBuilder.WriteString("select ")
	sqlBuilder.WriteString(conver.TableName())
	sqlBuilder.WriteString(".* from ")
	sqlBuilder.WriteString(conver.TableName())
	params := []interface{}{}
	//condition := fmt.Sprintf("%s.app_key=?", conver.TableName())
	conditionBuilder.WriteString(conver.TableName())
	conditionBuilder.WriteString(".app_key=?")
	params = append(params, appkey)
	if userId != "" {
		// condition = condition + fmt.Sprintf(" and %s.user_id=?", conver.TableName())
		conditionBuilder.WriteString(" and ")
		conditionBuilder.WriteString(conver.TableName())
		conditionBuilder.WriteString(".user_id=?")
		params = append(params, userId)
	}
	if targetId != "" {
		// condition = condition + fmt.Sprintf(" and %s.target_id=?", conver.TableName())
		conditionBuilder.WriteString(" and ")
		conditionBuilder.WriteString(conver.TableName())
		conditionBuilder.WriteString(".target_id=?")
		params = append(params, targetId)
	}
	if channelType != pbobjs.ChannelType_Unknown {
		// condition = condition + fmt.Sprintf(" and %s.channel_type=?", conver.TableName())
		conditionBuilder.WriteString(" and ")
		conditionBuilder.WriteString(conver.TableName())
		conditionBuilder.WriteString(".channel_type=?")
		params = append(params, int(channelType))
	}
	if tag != "" {
		rel := &ConverTagRelDao{}
		// sql = sql + fmt.Sprintf(" inner join %s on (%s.target_id=%s.target_id and %s.channel_type=%s.channel_type)", rel.TableName(), conver.TableName(), rel.TableName(), conver.TableName(), rel.TableName())
		sqlBuilder.WriteString(" inner join ")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(" on (")
		sqlBuilder.WriteString(conver.TableName())
		sqlBuilder.WriteString(".target_id=")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(".target_id and ")
		sqlBuilder.WriteString(conver.TableName())
		sqlBuilder.WriteString(".channel_type=")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(".channel_type)")
		// condition = condition + fmt.Sprintf(" and %s.tag = ?", rel.TableName())
		conditionBuilder.WriteString(" and ")
		conditionBuilder.WriteString(rel.TableName())
		conditionBuilder.WriteString(".tag=?")
		params = append(params, tag)
	}
	// condition = condition + " and is_deleted=?"
	conditionBuilder.WriteString(" and is_deleted=?")
	params = append(params, 0)
	orderStr := "sort_time desc"
	if isPositiveOrder {
		// condition = condition + " and sort_time>?"
		conditionBuilder.WriteString(" and sort_time>?")
		params = append(params, startTime)
		orderStr = "sort_time asc"
	} else {
		// condition = condition + " and sort_time<?"
		conditionBuilder.WriteString(" and sort_time<?")
		params = append(params, startTime)
	}
	// sql = sql + " where " + condition
	sqlBuilder.WriteString(" where ")
	sqlBuilder.WriteString(conditionBuilder.String())
	//err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	err := dbcommons.GetDb().Raw(sqlBuilder.String(), params...).Order(orderStr).Limit(count).Find(&items).Error
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
		})
	}
	return conversations, nil
}

func (conver *ConversationDao) DelConversation(appkey, userId, targetId string, channelType pbobjs.ChannelType) error {
	upd := map[string]interface{}{}
	upd["is_deleted"] = 1
	upd["latest_msg_id"] = ""
	upd["latest_msg"] = []byte{}
	upd["unread_tag"] = 0
	// upd["latest_unread_msg_index"] = 0
	// upd["latest_read_msg_index"] = 0
	return dbcommons.GetDb().Model(conver).Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Update(upd).Error
}

func (conver *ConversationDao) UpdateLatestReadMsgIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIndex int64, readMsgId string, readMsgTime int64) (int64, error) {
	upd := map[string]interface{}{}
	upd["latest_read_msg_index"] = msgIndex
	upd["latest_read_msg_id"] = readMsgId
	upd["latest_read_msg_time"] = readMsgTime
	upd["unread_tag"] = 0
	result := dbcommons.GetDb().Model(conver).Where("app_key=? and user_id=? and target_id=? and channel_type=? and latest_read_msg_index<=? and latest_unread_msg_index>=?", appkey, userId, targetId, channelType, msgIndex, msgIndex).Update(upd)
	return result.RowsAffected, result.Error
}

func (conver *ConversationDao) UpdateIsTopState(appkey, userId, targetId string, channelType pbobjs.ChannelType, isTop int, optTime int64) (int64, error) {
	upd := map[string]interface{}{}
	upd["is_top"] = isTop
	upd["top_updated_time"] = optTime
	result := dbcommons.GetDb().Model(conver).Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Update(upd)
	return result.RowsAffected, result.Error
}

type UnreadCount struct {
	UnreadCount int64 `gorm:"unread_count"`
}

func (conver *ConversationDao) TotalUnreadCount(appkey, userId string, channelTypes []pbobjs.ChannelType, excludeConvers []*pbobjs.SimpleConversation, tag string) int64 {
	var unreadCount UnreadCount
	params := []interface{}{}
	var sqlBuilder strings.Builder
	var conditionBuilder strings.Builder
	// sql := "SELECT SUM(CASE WHEN latest_unread_msg_index=latest_read_msg_index AND unread_tag=1 THEN 1 ELSE latest_unread_msg_index-latest_read_msg_index END) AS unread_count FROM conversations WHERE app_key=? and user_id=?"
	sqlBuilder.WriteString("SELECT SUM(CASE WHEN latest_unread_msg_index=latest_read_msg_index AND unread_tag=1 THEN 1 ELSE latest_unread_msg_index-latest_read_msg_index END) AS unread_count FROM ")
	sqlBuilder.WriteString(conver.TableName())

	conditionBuilder.WriteString(conver.TableName())
	conditionBuilder.WriteString(".app_key=?")
	params = append(params, appkey)
	conditionBuilder.WriteString(" AND ")
	conditionBuilder.WriteString(conver.TableName())
	conditionBuilder.WriteString(".user_id=?")
	params = append(params, userId)
	// params = append(params, userId)
	if tag != "" {
		rel := &ConverTagRelDao{}
		sqlBuilder.WriteString(" INNER JOIN ")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(" ON (")
		sqlBuilder.WriteString(conver.TableName())
		sqlBuilder.WriteString(".target_id=")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(".target_id and ")
		sqlBuilder.WriteString(conver.TableName())
		sqlBuilder.WriteString(".channel_type=")
		sqlBuilder.WriteString(rel.TableName())
		sqlBuilder.WriteString(".channel_type)")

		conditionBuilder.WriteString(" AND ")
		conditionBuilder.WriteString(rel.TableName())
		conditionBuilder.WriteString(".tag=?")
		params = append(params, tag)
	}

	if len(channelTypes) > 0 {
		channels := []int{}
		for _, c := range channelTypes {
			channels = append(channels, int(c))
		}
		conditionBuilder.WriteString(" AND ")
		conditionBuilder.WriteString(conver.TableName())
		conditionBuilder.WriteString(".channel_type in (?)")
		// sql = sql + " and channel_type in (?)"
		params = append(params, channels)
	}
	if len(excludeConvers) > 0 {
		tableName := conver.TableName()
		for _, conver := range excludeConvers {
			// sql = sql + " and (target_id!=? or channel_type!=?)"
			conditionBuilder.WriteString(" AND ")
			conditionBuilder.WriteString("(")
			conditionBuilder.WriteString(tableName)
			conditionBuilder.WriteString(".target_id!=? OR ")
			conditionBuilder.WriteString(tableName)
			conditionBuilder.WriteString(".channel_type!=?)")
			params = append(params, conver.TargetId)
			params = append(params, int(conver.ChannelType))
		}
	}
	conditionBuilder.WriteString(" AND is_deleted=0")
	conditionBuilder.WriteString(" AND latest_read_msg_index!=latest_unread_msg_index and undisturb_type=0")
	// sql = sql + " and is_deleted=0 and latest_read_msg_index!=latest_unread_msg_index and undisturb_type=0"
	sqlBuilder.WriteString(" WHERE ")
	sqlBuilder.WriteString(conditionBuilder.String())
	err := dbcommons.GetDb().Raw(sqlBuilder.String(), params...).Scan(&unreadCount).Error
	if err != nil {
		return 0
	} else {
		return unreadCount.UnreadCount
	}
}

func (conver *ConversationDao) ClearTotalUnreadCount(appkey, userId string) error {
	return dbcommons.GetDb().Exec("UPDATE conversations set latest_read_msg_index=latest_unread_msg_index, latest_read_msg_time =(UNIX_TIMESTAMP(NOW(3)) * 1000), unread_tag=0 where app_key=? and user_id=?", appkey, userId).Error
}

func (conver *ConversationDao) QryTopConvers(appkey, userId string, startTime, limit int64, sortType pbobjs.TopConverSortType, isPositive bool) ([]*models.Conversation, error) {
	var items []*ConversationDao
	var err error
	params := []interface{}{}
	condition := "app_key=? and user_id=? and is_top=1"
	params = append(params, appkey)
	params = append(params, userId)
	var orderBy string
	if sortType == pbobjs.TopConverSortType_BySortTime {
		if isPositive {
			orderBy = "sort_time asc"
			if startTime > 0 {
				condition = condition + " and sort_time>?"
				params = append(params, startTime)
			}
		} else {
			orderBy = "sort_time desc"
			if startTime > 0 {
				condition = condition + " and sort_time<?"
				params = append(params, startTime)
			}
		}
	} else {
		if isPositive {
			orderBy = "top_updated_time asc"
			if startTime > 0 {
				condition = condition + " and top_updated_time>?"
				params = append(params, startTime)
			}
		} else {
			orderBy = "top_updated_time desc"
			if startTime > 0 {
				condition = condition + " and top_updated_time<?"
				params = append(params, startTime)
			}
		}
	}
	err = dbcommons.GetDb().Where(condition, params...).Order(orderBy).Limit(limit).Find(&items).Error
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
			LatestReadMsgIndex:   item.LatestReadMsgIndex,
			LatestReadMsgId:      item.LatestReadMsgId,
			LatestReadMsgTime:    item.LatestReadMsgTime,
			SyncTime:             item.SyncTime,
			UndisturbType:        item.UndisturbType,
			IsTop:                item.IsTop,
			TopUpdatedTime:       item.TopUpdatedTime,
			IsDeleted:            item.IsDeleted,
			AppKey:               item.AppKey,
		})
	}
	return conversations, nil
}

func (conver *ConversationDao) UpdateUndisturbType(appkey, userId, targetId string, channelType pbobjs.ChannelType, undisturbType int32) (int64, error) {
	result := dbcommons.GetDb().Model(&ConversationDao{}).Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Update("undisturb_type", undisturbType)
	return result.RowsAffected, result.Error
}

func (conver *ConversationDao) FindUndisturb(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*models.Conversation, error) {
	var item ConversationDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Select("app_key,user_id,target_id,channel_type,undisturb_type").Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Conversation{
		AppKey:        item.AppKey,
		UserId:        item.UserId,
		TargetId:      item.TargetId,
		ChannelType:   pbobjs.ChannelType(item.ChannelType),
		UndisturbType: item.UndisturbType,
	}, nil
}
func (conver *ConversationDao) FindUnreadIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*models.Conversation, error) {
	var item ConversationDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Select("app_key,user_id,target_id,channel_type,latest_unread_msg_index").Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.Conversation{
		AppKey:               item.AppKey,
		UserId:               item.UserId,
		TargetId:             item.TargetId,
		ChannelType:          pbobjs.ChannelType(item.ChannelType),
		LatestUnreadMsgIndex: item.LatestUnreadMsgIndex,
	}, nil
}

func (conver *ConversationDao) UpdateLatestMsgBody(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgId string, msgBs []byte) error {
	return dbcommons.GetDb().Model(&ConversationDao{}).Where("app_key=? and user_id=? and target_id=? and channel_type=? and latest_msg_id=?", appkey, userId, targetId, channelType, msgId).Update("latest_msg", msgBs).Error
}

func (conver *ConversationDao) UpdateUnreadTag(appkey, userId, targetId string, channelType pbobjs.ChannelType) (int64, error) {
	result := dbcommons.GetDb().Model(&ConversationDao{}).Where("app_key=? and user_id=? and target_id=? and channel_type=?", appkey, userId, targetId, channelType).Update("unread_tag", 1)
	return result.RowsAffected, result.Error
}
