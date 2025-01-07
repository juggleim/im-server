package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgtypes"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/breaker"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

func init() {
	logx.SetLevel(logx.ErrorLevel)
}

var retryStrategy = []retry.Option{
	retry.Delay(100 * time.Millisecond),
	retry.Attempts(3),
	retry.LastErrorOnly(true),
	retry.DelayType(retry.BackOffDelay),
}

func MsgSubHandle(ctx context.Context, msgs *pbobjs.SubMsgs) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := commonservices.GetAppInfo(appKey)
	if ok {
		//check subscription config
		if appInfo.EventSubConfigObj == nil || appInfo.EventSubConfigObj.EventSubUrl == "" {
			logs.WithContext(ctx).Errorf("no event_sub_url or event_sub_auth. [%s]", appInfo.EventSubConfig)
			return
		}
		event := createEvent(msgs)
		if event == nil {
			return
		}

		msgIds := lo.Map(msgs.SubMsgs, func(msg *pbobjs.SubMsg, i int) string {
			return msg.Msg.MsgId
		})

		nonce := tools.RandStr(8)
		tsStr := fmt.Sprintf("%d", time.Now().UnixMilli())
		headers := map[string]string{
			"Content-Type": "application/json",
			"appkey":       appKey,
			"nonce":        nonce,
			"timestamp":    tsStr,
			"signature":    tools.SHA1(fmt.Sprintf("%s%s%s", appInfo.AppSecret, nonce, tsStr)),
		}
		body := tools.ToJson(event)

		// 以推送路由为key，建立熔断机制，滑动窗口内失败率过高，则熔断，直到窗口内成功率达到阈值，恢复
		err := breaker.Do(appInfo.EventSubConfigObj.EventSubUrl, func() error {
			return retry.Do(func() error {
				return notify(appInfo.EventSubConfigObj.EventSubUrl, headers, body)
			}, retryStrategy...)
		})

		if err == nil { //success
			logs.WithContext(ctx).Infof("msg sub success:%v", msgIds)
		} else { //failed
			logs.WithContext(ctx).Errorf("msg sub failed:%v\terr:%v", msgIds, err)
		}
	}
}

func OnlineOfflineHandle(ctx context.Context, msgs *pbobjs.OnlineOfflineMsg) {
	appKey := bases.GetAppKeyFromCtx(ctx)
	appInfo, ok := commonservices.GetAppInfo(appKey)
	if ok {
		//check subscription config
		if appInfo.EventSubConfigObj == nil || appInfo.EventSubConfigObj.EventSubUrl == "" {
			logs.WithContext(ctx).Errorf("no event_sub_url or event_sub_auth. [%s]", appInfo.EventSubConfig)
			return
		}
		event := createEvent(msgs)
		if event == nil {
			return
		}

		nonce := tools.RandStr(8)
		tsStr := fmt.Sprintf("%d", time.Now().UnixMilli())
		headers := map[string]string{
			"Content-Type": "application/json",
			"appkey":       appKey,
			"nonce":        nonce,
			"timestamp":    tsStr,
			"signature":    tools.SHA1(fmt.Sprintf("%s%s%s", appInfo.AppSecret, nonce, tsStr)),
			"ext":          msgs.ConnectionExt,
		}
		body := tools.ToJson(event)

		// 以推送路由为key，建立熔断机制，滑动窗口内失败率过高，则熔断，直到窗口内成功率达到阈值，恢复
		err := breaker.Do(appInfo.EventSubConfigObj.EventSubUrl, func() error {
			return retry.Do(func() error {
				return notify(appInfo.EventSubConfigObj.EventSubUrl, headers, body)
			}, retryStrategy...)
		})

		logs.WithContext(ctx).Infof("online sub result, userId:%v, eventType:%v, err:%v", msgs.UserId, msgs.Type, err)
	}
}

func notify(url string, headers map[string]string, body string) error {
	code, res, err := HttpDoBytes("POST", url, headers, body)
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("msg sub failed:%v\tcode:%d body:%s", url, code, res)
	}
	return nil
}

func createEvent(msg proto.Message) *SubEvent {
	switch msg.(type) {
	case *pbobjs.SubMsgs:
		msgs := msg.(*pbobjs.SubMsgs)
		event := &SubEvent{
			EventType: EventType_Message,
			Timestamp: time.Now().UnixMilli(),
			Payload:   []interface{}{},
		}
		for _, msg := range msgs.SubMsgs {
			event.Payload = append(event.Payload, &MsgEvent{
				Platform:    msg.Platform,
				Sender:      msg.Msg.SenderId,
				Receiver:    msg.Msg.TargetId,
				ConverType:  int(msg.Msg.ChannelType),
				MsgType:     msg.Msg.MsgType,
				MsgContent:  string(msg.Msg.MsgContent),
				MsgId:       msg.Msg.MsgId,
				MsgTime:     msg.Msg.MsgTime,
				MentionInfo: transMentionInfo(msg.Msg.MentionInfo),
			})
		}
		return event
	case *pbobjs.OnlineOfflineMsg:
		msg := msg.(*pbobjs.OnlineOfflineMsg)
		var eventType EventType
		if msg.Type == pbobjs.OnlineType_Online {
			eventType = EventType_Online
		} else {
			eventType = EventType_Offline
		}
		onlineEvent := &OnlineEvent{
			Type:             int32(msg.Type),
			DepUserId:        msg.UserId,
			UserId:           msg.UserId,
			DepDeviceId:      msg.DeviceId,
			DeviceId:         msg.DeviceId,
			Platform:         msg.Platform,
			DepClientIp:      msg.ClientIp,
			ClientIp:         msg.ClientIp,
			DepSessionId:     msg.SessionId,
			SessionId:        msg.SessionId,
			Timestamp:        msg.Timestamp,
			DepConnectionExt: msg.ConnectionExt,
			ConnectionExt:    msg.ConnectionExt,
			InstanceId:       msg.InstanceId,
		}
		event := &SubEvent{
			EventType: eventType,
			Timestamp: time.Now().UnixMilli(),
			Payload: []interface{}{
				onlineEvent,
			},
		}
		return event
	}

	return nil
}

func transMentionInfo(mention *pbobjs.MentionInfo) *MentionInfo {
	if mention != nil {
		mentionType := ""
		if mention.MentionType == pbobjs.MentionType_All {
			mentionType = msgtypes.MentionType_All
		} else if mention.MentionType == pbobjs.MentionType_Someone {
			mentionType = msgtypes.MentionType_Someone
		} else if mention.MentionType == pbobjs.MentionType_AllAndSomeone {
			mentionType = msgtypes.MentionType_AllSomeone
		}
		userIds := []string{}
		for _, user := range mention.TargetUsers {
			userIds = append(userIds, user.UserId)
		}
		return &MentionInfo{
			MentionType:   mentionType,
			TargetUserIds: userIds,
		}
	}
	return nil
}
