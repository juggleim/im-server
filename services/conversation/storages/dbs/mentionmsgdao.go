package dbs

import (
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/conversation/storages/models"
	"sort"
	"strings"
)

type MentionMsgDao struct {
	ID          int64  `gorm:"primary_key"`
	UserId      string `gorm:"user_id"`
	TargetId    string `gorm:"target_id"`
	ChannelType int    `gorm:"channel_type"`
	SubChannel  string `gorm:"sub_channel"`
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
		SubChannel:  item.SubChannel,
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

func (mention *MentionMsgDao) QryMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, startIndex int64, cleanTime int64) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao

	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=?"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, channelType)
	params = append(params, subChannel)
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

	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			UserId:      item.UserId,
			TargetId:    item.TargetId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			SubChannel:  item.SubChannel,
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

func (mention *MentionMsgDao) QryUnreadMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startTime int64, count int, isPositiveOrder bool, cleanTime int64) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao
	params := []interface{}{}
	condition := "app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and is_read=0"
	params = append(params, appkey)
	params = append(params, userId)
	params = append(params, targetId)
	params = append(params, channelType)
	params = append(params, subChannel)
	orderStr := "msg_time desc"
	if isPositiveOrder {
		condition = condition + " and msg_time>? and msg_time>?"
		orderStr = "msg_time asc"
	} else {
		condition = condition + " and msg_time<? and msg_time>?"
	}
	params = append(params, startTime)
	params = append(params, cleanTime)

	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			UserId:      item.UserId,
			TargetId:    item.TargetId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			SubChannel:  item.SubChannel,
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

func (mention *MentionMsgDao) QryMentionSenderIdsBaseIndex(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, startIndex int64, count int) ([]*models.MentionMsg, error) {
	var items []MentionMsgDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and msg_index>?", appkey, userId, targetId, int(channelType), subChannel, startIndex).Select("sender_id,msg_id,msg_time,msg_index,mention_type").Order("msg_index desc").Limit(count).Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			SenderId:    item.SenderId,
			MsgTime:     item.MsgTime,
			MsgId:       item.MsgId,
			MsgIndex:    item.MsgIndex,
			MentionType: pbobjs.MentionType(item.MentionType),
			SubChannel:  item.SubChannel,
		})
	}
	sort.Slice(mentionMsgs, func(i, j int) bool {
		return mentionMsgs[i].MsgIndex < mentionMsgs[j].MsgIndex
	})
	return mentionMsgs, nil
}

func (mention *MentionMsgDao) BatchQryMentionSenderIdsBaseIndex(appkey, userId string, convers []models.ConverItem) ([]*models.MentionMsg, error) {
	length := len(convers)
	if length <= 0 {
		return []*models.MentionMsg{}, nil
	}
	var items []MentionMsgDao
	var sqlBuilder strings.Builder
	params := []interface{}{}
	sqlBuilder.WriteString("app_key=? and user_id=? and (")
	params = append(params, appkey)
	params = append(params, userId)
	for i, conver := range convers {
		if i == length-1 {
			sqlBuilder.WriteString("(target_id=? and channel_type=? and sub_channel=? and msg_index>?)")
		} else {
			sqlBuilder.WriteString("(target_id=? and channel_type=? and sub_channel=? and msg_index>?) or ")
		}
		params = append(params, conver.TargetId)
		params = append(params, conver.ChannelType)
		params = append(params, conver.SubChannel)
		params = append(params, conver.MsgIndex)
	}
	sqlBuilder.WriteString(")")
	err := dbcommons.GetDb().Where(sqlBuilder.String(), params...).Select("target_id,channel_type,sub_channel,sender_id,msg_id,msg_time").Order("msg_index asc").Find(&items).Error
	if err != nil {
		return []*models.MentionMsg{}, err
	}
	mentionMsgs := []*models.MentionMsg{}
	for _, item := range items {
		mentionMsgs = append(mentionMsgs, &models.MentionMsg{
			TargetId:    item.TargetId,
			ChannelType: pbobjs.ChannelType(item.ChannelType),
			SubChannel:  item.SubChannel,
			SenderId:    item.SenderId,
			MsgTime:     item.MsgTime,
			MsgId:       item.MsgId,
		})
	}
	return mentionMsgs, nil
}

func (mention *MentionMsgDao) MarkRead(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIds []string) error {
	return dbcommons.GetDb().Model(&MentionMsgDao{}).Where("app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and msg_id in (?)", appkey, userId, targetId, channelType, subChannel, msgIds).Update("is_read", 1).Error
}

func (mention *MentionMsgDao) DelMentionMsgs(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and msg_id in (?)", appkey, userId, targetId, channelType, subChannel, msgIds).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) DelMentionMsg(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and msg_id=?", appkey, userId, targetId, channelType, subChannel, msgId).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) CleanMentionMsgsBaseIndex(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType, msgIndex int64) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=? and target_id=? and channel_type=? and sub_channel=? and msg_index<=?", appkey, userId, targetId, channelType, subChannel, msgIndex).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) CleanMentionMsgsBaseUserId(appkey, userId string) error {
	return dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Delete(&MentionMsgDao{}).Error
}

func (mention *MentionMsgDao) DelOnlyByMsgIds(appkey string, msgIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and msg_id in (?)", appkey, msgIds).Delete(&MentionMsgDao{}).Error
}
