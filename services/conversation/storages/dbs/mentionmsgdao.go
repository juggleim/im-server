package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"sort"
)

type MentionMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	UserId      string `gorm:"user_id"`
	TargetId    string `gorm:"target_id"`
	ChannelType int    `gorm:"channel_type"`
	SenderId    string `gorm:"sender_id"`
	MentionType int    `gorm:"mention_type"`
	MsgId       string `gorm:"msg_id"`
	MsgTime     int64  `gorm:"msg_time"`
	MsgIndex    int64  `gorm:"msg_index"`
	IsRead      int    `gorm:"is_read"`
	AppKey      string `gorm:"app_key"`
}

func (mention *MentionMsgDao) TableName() string {
	return "mentionmsgs"
}
func (mention *MentionMsgDao) SaveMentionMsg(item models.MentionMsg) error {
	daoItem := MentionMsgDao{
		UserId:      item.UserId,
		TargetId:    item.TargetId,
		ChannelType: int(item.ChannelType),
		SenderId:    item.SenderId,
		MentionType: int(item.MentionType),
		MsgId:       item.MsgId,
		MsgTime:     item.MsgTime,
		MsgIndex:    item.MsgIndex,
		AppKey:      item.AppKey,
		IsRead:      item.IsRead,
	}
	err := dbcommons.GetDb().Create(&daoItem).Error
	return err
}

func (mention *MentionMsgDao) QryMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, startIndex int64, cleanTime int64) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao

	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and channel_type=?"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, channelType)
	if startIndex > 0 {
		condition = condition + " and msg_index>?"
		params = append(params, startIndex)
	}
	orderStr := "msg_time desc"
	if isPositiveOrder {
		condition = condition + " and msg_time>? and msg_time>?"
		orderStr = "msg_time asc"
	} else {
		condition = condition + " and msg_time<? and msg_time>?"
	}
	params = append(params, startTime)
	params = append(params, cleanTime)

	err := dbcommons.GetDb().Debug().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			UserId:      item.UserId,
			TargetId:    item.TargetId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			SenderId:    item.SenderId,
			MentionType: pbobjs.MentionType(item.MentionType),
			MsgId:       item.MsgId,
			MsgTime:     item.MsgTime,
			MsgIndex:    item.MsgIndex,
			AppKey:      item.AppKey,
			IsRead:      item.IsRead,
		})
	}
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].MsgTime < items[j].MsgTime
		})
	}
	return mentionMsgs, nil
}

func (mention *MentionMsgDao) QryUnreadMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, cleanTime int64) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao
	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and channel_type=? and is_read=0"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, channelType)
	orderStr := "msg_time desc"
	if isPositiveOrder {
		condition = condition + " and msg_time>? and msg_time>?"
		orderStr = "msg_time asc"
	} else {
		condition = condition + " and msg_time<? and msg_time>?"
	}
	params = append(params, startTime)
	params = append(params, cleanTime)

	err := dbcommons.GetDb().Debug().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			UserId:      item.UserId,
			TargetId:    item.TargetId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			SenderId:    item.SenderId,
			MentionType: pbobjs.MentionType(item.MentionType),
			MsgId:       item.MsgId,
			MsgTime:     item.MsgTime,
			MsgIndex:    item.MsgIndex,
			AppKey:      item.AppKey,
			IsRead:      item.IsRead,
		})
	}
	if !isPositiveOrder {
		sort.Slice(items, func(i, j int) bool {
			return items[i].MsgTime < items[j].MsgTime
		})
	}
	return mentionMsgs, nil
}

func (mention *MentionMsgDao) QryMentionSenderIdsBaseIndex(appkey, userId, targetId string, channelType pbobjs.ChannelType, startIndex int64, count int) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and msg_index>?", appkey, userId, targetId, int(channelType), startIndex).Select("sender_id,msg_id,msg_time").Order("msg_index asc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			SenderId: item.SenderId,
			MsgTime:  item.MsgTime,
			MsgId:    item.MsgId,
		})
	}
	return mentionMsgs, nil
}

func (mention *MentionMsgDao) MarkRead(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIds []string) error {
	return dbcommons.GetDb().Model(&MentionMsgDao{}).Where("app_key=? and user_id=? and target_id=? and channel_type=? and msg_id in (?)", appkey, userId, targetId, channelType, msgIds).Update("is_read", 1).Error
}

func (mention *MentionMsgDao) DelMentionMsgs(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and msg_id in (?)", appkey, userId, targetId, channelType, msgIds).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) DelMentionMsg(appkey, userId, targetId string, channelType pbobjs.ChannelType, msgId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and msg_id=?", appkey, userId, targetId, channelType, msgId).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) CleanMentionMsgsBaseTime(appkey, userId, targetId string, channelType pbobjs.ChannelType, cleanTime int64) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and msg_time<=?", appkey, userId, targetId, channelType, cleanTime).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) DelOnlyByMsgIds(appkey string, msgIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and msg_id in (?)", appkey, msgIds).Delete(&MentionMsgDao{}).Error
}
