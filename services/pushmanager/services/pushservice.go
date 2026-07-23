package services

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/pushmanager/services/getuipush"
	"im-server/services/pushmanager/services/honorpush"
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
	params := map[string]string{
		"senderName": req.SenderName,
		"pushText":   req.PushText,
		"groupName":  req.GroupName,
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	pushToken := GetPushToken(appkey, userId)
	if pushToken.PushToken != "" || pushToken.VoipPushToken != "" {
		if pushToken.Platform == pbobjs.Platform_iOS {
			iosPushConf := GetIosPushConf(ctx, appkey, pushToken.PackageName)
			if iosPushConf != nil {
				notification := &apns2.Notification{}
				notification.DeviceToken = pushToken.PushToken
				notification.Topic = pushToken.PackageName
				notification.Payload = iosPushPayload(req)
				var client *apns2.Client
				if req.IsVoip && iosPushConf.ApnsVoipClient != nil && pushToken.VoipPushToken != "" {
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
		} else if pushToken.Platform == pbobjs.Platform_Android {
			androidPushConf := GetAndroidPushConf(ctx, appkey, pushToken.PackageName)
			if androidPushConf != nil {
				switch pushToken.PushChannel {
				case pbobjs.PushChannel_Huawei:
					if androidPushConf.HwPushClient != nil {
						hwNotification := &hwpush.AndroidNotification{
							Title:        req.Title,
							Body:         req.PushText,
							DefaultSound: true,
							ClickAction: &hwpush.ClickAction{
								Type: 3,
							},
						}
						// 华为角标依赖入口 Activity 全类名（badge.class），未配置时不下发角标
						if badgeClass := androidPushConf.HwPushClient.BadgeClass; badgeClass != "" {
							hwNotification.Badge = &hwpush.BadgeNotification{
								Class: badgeClass,
							}
							if req.Badge > 0 {
								hwNotification.Badge.SetNum = int(req.Badge)
							} else {
								hwNotification.Badge.AddNum = 1
							}
						}
						_, err := androidPushConf.HwPushClient.SendMessage(ctx, &hwpush.MessageRequest{
							Message: &hwpush.Message{
								Notification: &hwpush.Notification{
									Title: req.Title,
									Body:  req.PushText,
								},
								Android: &hwpush.AndroidConfig{
									Notification: hwNotification,
								},
								Token: []string{pushToken.PushToken},
							},
						})
						if err != nil {
							logs.WithContext(ctx).Infof("[Huawei_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Huawei_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Huawei_FAIL]have no init jpush client")
					}
				case pbobjs.PushChannel_Xiaomi:
					if androidPushConf.XiaomiPushClient != nil {
						_, err := androidPushConf.XiaomiPushClient.SendWithContext(ctx, &xiaomipush.SendReq{
							RestrictedPackageName: pushToken.PackageName,
							Title:                 req.Title,
							Description:           req.PushText,
							RegistrationId:        pushToken.PushToken,
							Extra: &xiaomipush.Extra{
								ChannelId: androidPushConf.XiaomiPushClient.ChannelId,
							},
						})
						if err != nil {
							logs.WithContext(ctx).Infof("[Xiaomi_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Xiaomi_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Xiaomi_FAIL]have no init jpush client")
					}
				case pbobjs.PushChannel_Oppo:
					if androidPushConf.OppoPushClient != nil {
						_, err := androidPushConf.OppoPushClient.SendWithContext(ctx, &oppopush.SendReq{
							Notification: &oppopush.Notification{
								Title:     req.Title,
								Content:   req.PushText,
								ChannelID: androidPushConf.OppoPushClient.ChannelId,
							},
							TargetType:  2,
							TargetValue: pushToken.PushToken,
						})
						if err != nil {
							logs.WithContext(ctx).Infof("[Oppo_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Oppo_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Oppo_FAIL]have no init jpush client")
					}
				case pbobjs.PushChannel_Vivo:
					if androidPushConf.VivoPushClient != nil {
						_, err := androidPushConf.VivoPushClient.SendWithContext(ctx, &vivopush.SendReq{
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
						if err != nil {
							logs.WithContext(ctx).Infof("[Vivo_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Vivo_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Vivo_FAIL]have no init jpush client")
					}
				case pbobjs.PushChannel_JPush:
					if androidPushConf.JpushClient != nil {
						intentUrl := "intent:#Intent;action=com.j.im.intent.MESSAGE_CLICK;%send"
						intentComponent := ""
						if androidPushConf.Package != "" && androidPushConf.JpushClient.BadgeClass != "" {
							intentComponent = fmt.Sprintf("component=%s/%s;", androidPushConf.Package, androidPushConf.JpushClient.BadgeClass)
						}
						intentUrl = fmt.Sprintf(intentUrl, intentComponent)
						androidNotification := &jpush.AndroidNotification{
							Alert:  req.PushText,
							Title:  req.Title,
							Extras: transfer2Exts(req),
							Intent: map[string]interface{}{
								"url": intentUrl,
							},
						}
						if req.Badge > 0 {
							androidNotification.BadgeSetNum = int(req.Badge)
						} else {
							androidNotification.BadgeAddNum = 1
						}
						androidNotification.BadgeClass = androidPushConf.JpushClient.BadgeClass
						jPushPayload := &jpush.Payload{
							Platform: jpush.NewPlatform().All(),
							Audience: jpush.NewAudience().SetRegistrationId(pushToken.PushToken),
							Notification: &jpush.Notification{
								Alert:   req.Title,
								Android: androidNotification,
							},
						}
						if req.JPushOptions != "" {
							jPushPayload.Options = &jpush.Options{}
							tools.JsonUnMarshal([]byte(req.JPushOptions), jPushPayload.Options)
						} else {
							jPushPayload.Options = handleJPushOptions(androidPushConf.JpushClient.Options, params)
						}
						_, err := androidPushConf.JpushClient.Push(jPushPayload)
						if err != nil {
							logs.WithContext(ctx).Infof("[JPush_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[JPush_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[JPush_FAIL]have no init jpush client")
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
				case pbobjs.PushChannel_Honor:
					if androidPushConf.HonorPushClient != nil {
						honorNotification := &honorpush.PushAndroidNotification{
							Title: req.Title,
							Body:  req.PushText,
							ClickAction: &honorpush.PushAndroidClickAction{
								ActionType: honorpush.ClickActionTypeLaunchApp,
							},
						}
						// 荣耀角标依赖入口 Activity 全类名（badge.badgeClass），未配置时不下发角标
						if badgeClass := androidPushConf.HonorPushClient.BadgeClass; badgeClass != "" {
							honorNotification.Badge = &honorpush.PushBadge{
								BadgeClass: badgeClass,
							}
							if req.Badge > 0 {
								honorNotification.Badge.SetNum = uint16(req.Badge)
							} else {
								honorNotification.Badge.AddNum = 1
							}
						}
						_, err := androidPushConf.HonorPushClient.SendMessage(&honorpush.SendMessageReq{
							Token: []string{pushToken.PushToken},
							Android: &honorpush.PushAndroidConfig{
								Notification: honorNotification,
							},
						})
						if err != nil {
							logs.WithContext(ctx).Infof("[Honor_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Honor_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Honor_FAIL]have no init honor push client")
					}
				case pbobjs.PushChannel_Getui:
					if androidPushConf.GetuiPushClient != nil {
						_, err := androidPushConf.GetuiPushClient.ToSingleCIDWithContext(ctx, &getuipush.ToSingleCIDReq{
							RequestID: strconv.FormatInt(time.Now().UnixNano(), 10),
							Audience: &getuipush.AudienceCID{
								CID: []string{pushToken.PushToken},
							},
							PushMessage: &getuipush.PushMessage{
								Notification: &getuipush.Notification{
									Title: req.Title,
									Body:  req.PushText,
								},
							},
						})
						if err != nil {
							logs.WithContext(ctx).Infof("[Getui_ERROR]user_id:%s\tmsg_id:%s\t%s", userId, req.MsgId, err.Error())
						} else {
							logs.WithContext(ctx).Infof("[Getui_SUCC]user_id:%s\tmsg_id:%s", userId, req.MsgId)
						}
					} else {
						logs.WithContext(ctx).Infof("[Getui_FAIL]have no init getui push client")
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

func handleJPushOptions(jpushOptions *commonservices.JPushOptions, params map[string]string) *jpush.Options {
	if jpushOptions != nil {
		options := &jpush.Options{
			Classification: jpushOptions.Classification,
		}
		if jpushOptions.ThirdPartyChannel != nil {
			options.ThirdPartyChannel = &jpush.ThirdPartyChannel{
				Huawei: jpushOptions.ThirdPartyChannel.Huawei,
				Xiaomi: jpushOptions.ThirdPartyChannel.Xiaomi,
				Honor:  jpushOptions.ThirdPartyChannel.Honor,
				Oppo:   jpushOptions.ThirdPartyChannel.Oppo,
				Vivo:   jpushOptions.ThirdPartyChannel.Vivo,
				Meizu:  jpushOptions.ThirdPartyChannel.Meizu,
				Fcm:    jpushOptions.ThirdPartyChannel.Fcm,
				Nio:    jpushOptions.ThirdPartyChannel.Nio,
				Asus:   jpushOptions.ThirdPartyChannel.Asus,
				Hmos:   jpushOptions.ThirdPartyChannel.Hmos,
			}
			if options.ThirdPartyChannel.Xiaomi != nil && options.ThirdPartyChannel.Xiaomi.MiTemplateParam != "" {
				options.ThirdPartyChannel.Xiaomi.MiTemplateParam = tools.ReplaceTemplateParams(options.ThirdPartyChannel.Xiaomi.MiTemplateParam, params)
			}
			if options.ThirdPartyChannel.Oppo != nil {
				if len(options.ThirdPartyChannel.Oppo.PrivateContentParameters) > 0 {
					for k, v := range options.ThirdPartyChannel.Oppo.PrivateContentParameters {
						options.ThirdPartyChannel.Oppo.PrivateContentParameters[k] = tools.ReplaceTemplateParams(v, params)
					}
				}
				if len(options.ThirdPartyChannel.Oppo.PrivateTitleParameters) > 0 {
					for k, v := range options.ThirdPartyChannel.Oppo.PrivateTitleParameters {
						options.ThirdPartyChannel.Oppo.PrivateTitleParameters[k] = tools.ReplaceTemplateParams(v, params)
					}
				}
			}
		}
		return options
	}
	return nil
}
