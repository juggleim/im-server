package services

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages/models"
)

func resolveDelMsgTime(msg *pbobjs.SimpleMsg) (int64, error) {
	if msg == nil {
		return 0, tools.ErrInvalidMsgId
	}
	if msg.MsgTime > 0 {
		return msg.MsgTime, nil
	}
	msgTime, err := tools.ParseTimeFromMsgIdStrict(msg.MsgId)
	if err != nil {
		return 0, fmt.Errorf("parse delete msg time from msg_id %q: %w", msg.MsgId, err)
	}
	return msgTime, nil
}

func normalizeDelMsgTimes(msgs []*pbobjs.SimpleMsg) error {
	for _, msg := range msgs {
		if msg == nil {
			continue
		}
		msgTime, err := resolveDelMsgTime(msg)
		if err != nil {
			return err
		}
		msg.MsgTime = msgTime
	}
	return nil
}

func newPrivateDelHisMsg(userId, targetId, subChannel, appkey string, msg *pbobjs.SimpleMsg) models.PrivateDelHisMsg {
	return models.PrivateDelHisMsg{
		UserId: userId, TargetId: targetId, SubChannel: subChannel,
		MsgId: msg.GetMsgId(), MsgTime: msg.GetMsgTime(), AppKey: appkey,
	}
}

func newGroupDelHisMsg(userId, targetId, subChannel, appkey string, msg *pbobjs.SimpleMsg) models.GroupDelHisMsg {
	return models.GroupDelHisMsg{
		UserId: userId, TargetId: targetId, SubChannel: subChannel,
		MsgId: msg.GetMsgId(), MsgTime: msg.GetMsgTime(), AppKey: appkey,
	}
}
