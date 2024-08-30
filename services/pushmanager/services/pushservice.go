package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services/hwpush"
	"im-server/services/pushmanager/services/jpush"
	"im-server/services/pushmanager/services/oppopush"
	"im-server/services/pushmanager/services/vivopush"
	"im-server/services/pushmanager/services/xiaomipush"
	"strconv"
	"time"

	"github.com/sideshow/apns2"
)

func SendPush(ctx context.Context, userId string, req *pbobjs.PushData) {
	if req.PushText == "" {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	pushToken, ok := GetPushToken(appkey, userId)
	if ok && pushToken != nil {
		if pushToken.Platform == pbobjs.Platform_iOS {
			if pushToken.PushToken != "" {
				iosPushConf := GetIosPushConf(ctx, appkey, pushToken.PackageName)
				if iosPushConf != nil && iosPushConf.ApnsClient != nil {
					notification := &apns2.Notification{}
					notification.DeviceToken = pushToken.PushToken
					notification.Topic = pushToken.PackageName
					notification.Payload = []byte(fmt.Sprintf(`{"aps":{"alert":{"title":"%s","body":"%s"}},"conver_id":"%s","conver_type":"%d","exts":"%s"}`, req.Title, tools.PureStr(req.PushText), req.ConverId, req.ChannelType, req.PushExtraData))
					resp, err := iosPushConf.ApnsClient.Push(notification)
					if err != nil {
						logs.WithContext(ctx).Errorf("[IOS_FAIL]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
					} else {
						logs.WithContext(ctx).Infof("[IOS_SUCC]user_id:%s\tmsg_id:%s\tapns_id:%s\treason:%s\tcode:%d\ttime:%v", userId, req.MsgId, resp.ApnsID, resp.Reason, resp.StatusCode, resp.Timestamp)
					}
				}
			}
		} else if pushToken.Platform == pbobjs.Platform_Android {
			androidPushConf := GetAndroidPushConf(ctx, appkey, pushToken.PackageName)
			if androidPushConf != nil {
				switch pushToken.PushChannel {
				case pbobjs.PushChannel_Huawei:
					if androidPushConf.HwPushClient != nil {
						androidPushConf.HwPushClient.SendMessage(ctx, &hwpush.MessageRequest{
							Message: &hwpush.Message{
								Notification: &hwpush.Notification{
									Title: req.Title,
									Body:  req.PushText,
								},
								Android: &hwpush.AndroidConfig{
									Notification: &hwpush.AndroidNotification{
										Title:        req.Title,
										Body:         req.PushText,
										DefaultSound: true,
										ClickAction: &hwpush.ClickAction{
											Type: 3,
										},
									},
								},
								Token: []string{pushToken.PushToken},
							},
						})
					}
				case pbobjs.PushChannel_Xiaomi:
					if androidPushConf.XiaomiPushClient != nil {
						androidPushConf.XiaomiPushClient.SendWithContext(ctx, &xiaomipush.SendReq{
							RestrictedPackageName: pushToken.PackageName,
							Title:                 req.Title,
							Description:           req.PushText,
							RegistrationId:        pushToken.PushToken,
							Extra: &xiaomipush.Extra{
								ChannelId: "119572", // TODO
							},
						})
					}
				case pbobjs.PushChannel_Oppo:
					if androidPushConf.OppoPushClient != nil {
						androidPushConf.OppoPushClient.SendWithContext(ctx, &oppopush.SendReq{
							Notification: &oppopush.Notification{
								Title:   req.Title,
								Content: req.PushText,
								//ChannelID: channelId,
							},
							TargetType:  2,
							TargetValue: pushToken.PushToken,
						})
					}
				case pbobjs.PushChannel_Vivo:
					if androidPushConf.VivoPushClient != nil {
						androidPushConf.VivoPushClient.SendWithContext(ctx, &vivopush.SendReq{
							RegId:          pushToken.PushToken,
							NotifyType:     4,
							Title:          req.Title,
							Content:        req.PushText,
							TimeToLive:     24 * 60 * 60,
							SkipType:       1,
							NetworkType:    -1,
							Classification: 1,
							RequestId:      strconv.Itoa(int(time.Now().UnixNano())),
						})
					}
				case pbobjs.PushChannel_JPush:
					if androidPushConf.JpushClient != nil {
						androidPushConf.JpushClient.Push(&jpush.Payload{
							Platform: jpush.NewPlatform().All(),
							Audience: jpush.NewAudience().SetRegistrationId(pushToken.PushToken),
							Notification: &jpush.Notification{
								Alert: req.Title,
								Android: &jpush.AndroidNotification{
									Alert:  req.PushText,
									Title:  req.Title,
									Extras: transfer2Exts(req),
								},
							},
						})
					}
				default:
					fmt.Println("unknown push type", pushToken.PushChannel)
				}
			} else {
				fmt.Println("androidPushConf is nil")
			}
		}
	}
}

func transfer2Exts(pushData *pbobjs.PushData) map[string]interface{} {
	exts := make(map[string]interface{})
	exts["msg_id"] = pushData.MsgId
	exts["sender_id"] = pushData.SenderId
	exts["conver_id"] = pushData.ConverId
	exts["conver_type"] = int32(pushData.ChannelType)
	exts["exts"] = pushData.PushExtraData
	return exts
}
