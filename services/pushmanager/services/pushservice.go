package services

import (
	"context"
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
	"github.com/sideshow/apns2/payload"
)

func SendPush(ctx context.Context, userId string, req *pbobjs.PushData) {
	if req.PushText == "" {
		return
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	pushToken := GetPushToken(appkey, userId)
	if pushToken.PushToken != "" || pushToken.VoipPushToken != "" {
		if pushToken.Platform == pbobjs.Platform_iOS {
			if pushToken.PushToken != "" {
				iosPushConf := GetIosPushConf(ctx, appkey, pushToken.PackageName)
				if iosPushConf != nil && (iosPushConf.ApnsClient != nil || iosPushConf.ApnsVoipClient != nil) {
					notification := &apns2.Notification{}
					notification.DeviceToken = pushToken.PushToken
					notification.Topic = pushToken.PackageName
					notification.Payload = iosPushPayload(req)
					var client *apns2.Client
					if req.IsVoip && iosPushConf.ApnsVoipClient != nil {
						client = iosPushConf.ApnsVoipClient
						notification.Topic = notification.Topic + ".voip"
						notification.DeviceToken = pushToken.VoipPushToken
						notification.PushType = apns2.PushTypeVOIP
					} else {
						client = iosPushConf.ApnsClient
					}
					if client != nil {
						resp, err := client.Push(notification)
						if err != nil {
							logs.WithContext(ctx).Infof("[IOS_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							if resp.StatusCode == 200 {
								logs.WithContext(ctx).Infof("[IOS_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
							} else {
								logs.WithContext(ctx).Infof("[IOS_FAIL]user_id:%s\tmsg_id:%s\tcode:%d\treason:%s\tapns_id:%s", userId, req.MsgId, resp.StatusCode, resp.Reason, resp.ApnsID)
							}
						}
					} else {
						logs.WithContext(ctx).Infof("[IOS_ERR]user_id:%s\tnot init apns client")
					}
				} else {
					logs.WithContext(ctx).Infof("[IOS_CONF_NIL]app_key=%s\tpackage:%s", appkey, pushToken.PackageName)
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
				case pbobjs.PushChannel_FCM:
					if androidPushConf.FcmPushClient != nil {
						err := androidPushConf.FcmPushClient.SendPush(req.Title, req.PushText, pushToken.PushToken, transfer2Exts(req))
						if err != nil {
							logs.WithContext(ctx).Infof("[FCM_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[FCM_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[FCM_FAIL]have no init fcm push client")
					}
				default:
					logs.WithContext(ctx).Infof("unknown push type %s", pushToken.PushChannel)
				}
			} else {
				logs.WithContext(ctx).Info("android_push_conf is nil")
			}
		} else {
			logs.WithContext(ctx).Info("not support platform")
		}
	} else {
		logs.WithContext(ctx).Info("have no push token")
		//ntf close push
		bases.AsyncRpcCall(ctx, "upd_push_status", userId, &pbobjs.UserPushStatus{
			CanPush: false,
		})
	}
}

func transfer2Exts(pushData *pbobjs.PushData) map[string]interface{} {
	exts := make(map[string]interface{})
	if pushData.IsVoip {
		if pushData.RtcRoomId != "" {
			exts["room_id"] = pushData.RtcRoomId
		}
		if pushData.RtcInviterId != "" {
			exts["inviter_id"] = pushData.RtcInviterId
		}
		exts["is_multi"] = pushData.RtcRoomType
		exts["media_type"] = pushData.RtcMediaType
	} else {
		if pushData.MsgId != "" {
			exts["msg_id"] = pushData.MsgId
		}
		if pushData.SenderId != "" {
			exts["sender_id"] = pushData.SenderId
		}
		if pushData.ConverId != "" {
			exts["conver_id"] = pushData.ConverId
		}
		if pushData.ChannelType != pbobjs.ChannelType_Unknown {
			exts["conver_type"] = int32(pushData.ChannelType)
		}
	}
	if pushData.PushExtraData != "" {
		exts["exts"] = pushData.PushExtraData
	}
	return exts
}

func iosPushPayload(req *pbobjs.PushData) interface{} {
	iosPayload := payload.NewPayload()
	iosPayload.AlertTitle(req.Title)
	iosPayload.AlertBody(tools.PureStr(req.PushText))
	iosPayload.Sound("default")
	jcExts := transfer2Exts(req)
	for k, v := range jcExts {
		iosPayload.Custom(k, v)
	}
	if req.Badge > 0 {
		iosPayload.Badge(int(req.Badge))
	}
	return iosPayload
}
