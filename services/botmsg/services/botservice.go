package services

import (
	"bytes"
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/botmsg/services/botengines"
	"im-server/services/botmsg/storages"
	"im-server/services/botmsg/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"strings"
	"time"
)

var botCache *caches.LruCache
var botLocks *tools.SegmentatedLocks

func init() {
	botCache = caches.NewLruCacheWithAddReadTimeout(10000, nil, 5*time.Second, 5*time.Second)
	botLocks = tools.NewSegmentatedLocks(128)
}

type BotInfo struct {
	AppKey    string
	BotId     string
	Nickname  string
	Portrait  string
	BotType   models.BotType
	BotEngine botengines.IBotEngine
}

func GetBotInfo(ctx context.Context, botId string) *BotInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getKey(appkey, botId)
	if val, exist := botCache.Get(key); exist {
		return val.(*BotInfo)
	} else {
		l := botLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()
		if val, exist := botCache.Get(key); exist {
			return val.(*BotInfo)
		} else {
			botInfo := &BotInfo{
				AppKey:    appkey,
				BotId:     botId,
				BotEngine: &botengines.NilBotEngine{},
			}
			storage := storages.NewBotConfStorage()
			bot, err := storage.FindById(appkey, botId)
			if err == nil {
				botInfo.Nickname = bot.Nickname
				botInfo.Portrait = bot.BotPortrait
				botInfo.BotType = bot.BotType
				switch botInfo.BotType {
				case models.BotType_Dify:
					difyBot := &botengines.DifyBotEngine{}
					err = tools.JsonUnMarshal([]byte(bot.BotConf), difyBot)
					if err == nil && difyBot.ApiKey != "" && difyBot.Url != "" {
						botInfo.BotEngine = difyBot
					}
				case models.BotType_Coze:
					cozeBot := &botengines.CozeBotEngine{}
					err = tools.JsonUnMarshal([]byte(bot.BotConf), cozeBot)
					if err == nil && cozeBot.BotId != "" && cozeBot.Url != "" && cozeBot.Token != "" {
						botInfo.BotEngine = cozeBot
					}
				case models.BotType_SiliconFlow:
					sfBot := &botengines.SiliconFlowEngine{}
					err = tools.JsonUnMarshal([]byte(bot.BotConf), sfBot)
					if err == nil && sfBot.ApiKey != "" && sfBot.Model != "" && sfBot.Url != "" {
						botInfo.BotEngine = sfBot
					}
				}
			} else {
				botInfo.BotEngine = &botengines.NilBotEngine{}
			}
			botCache.Add(key, botInfo)
			return botInfo
		}
	}
}

func getKey(appkey, botId string) string {
	return strings.Join([]string{appkey, botId}, "_")
}

type Combiner struct {
	StreamMsg *pbobjs.DownMsg
	tmpbuf    *bytes.Buffer
	finalbuf  *bytes.Buffer
	ts        int64
	interval  int64
	subSeq    int64
}

func (combiner *Combiner) Append(part string) string {
	if combiner.finalbuf == nil {
		combiner.finalbuf = bytes.NewBuffer([]byte{})
	}
	if combiner.tmpbuf == nil {
		combiner.tmpbuf = bytes.NewBuffer([]byte{})
	}
	combiner.tmpbuf.WriteString(part)
	combiner.finalbuf.WriteString(part)
	cur := time.Now().UnixMilli()
	if combiner.ts == 0 {
		combiner.ts = cur
	}
	if cur-combiner.ts > combiner.interval {
		ret := combiner.tmpbuf.String()
		combiner.tmpbuf = nil
		combiner.ts = cur
		return ret
	}
	return ""
}

func (combiner *Combiner) GetLeft() string {
	if combiner.tmpbuf != nil {
		return combiner.tmpbuf.String()
	}
	return ""
}

func (combiner *Combiner) GetFinal() string {
	if combiner.finalbuf != nil {
		return combiner.finalbuf.String()
	}
	return ""
}

func (combiner *Combiner) GetSubSeq() int64 {
	combiner.subSeq = combiner.subSeq + 1
	return combiner.subSeq
}

func HandleBotMsg(ctx context.Context, msg *pbobjs.DownMsg) {
	if msg.MsgType != msgdefines.InnerMsgType_Text || (msg.ChannelType != pbobjs.ChannelType_Private && msg.ChannelType != pbobjs.ChannelType_Group) {
		return
	}
	txtMsg := &msgdefines.TextMsg{}
	err := tools.JsonUnMarshal(msg.MsgContent, txtMsg)
	if err != nil {
		logs.WithContext(ctx).Errorf("text msg illigal. content:%s", string(msg.MsgContent))
		return
	}
	botId := bases.GetTargetIdFromCtx(ctx)
	botInfo := GetBotInfo(ctx, botId)
	if botInfo.BotEngine != nil {
		converKey := ""
		if msg.ChannelType == pbobjs.ChannelType_Private {
			converKey = commonservices.GetConversationId(msg.SenderId, botId, pbobjs.ChannelType_Private)
			converKey = fmt.Sprintf("%s_%d", converKey, pbobjs.ChannelType_Private)
		} else if msg.ChannelType == pbobjs.ChannelType_Group {
			converKey = commonservices.GetConversationId(msg.SenderId, msg.TargetId, pbobjs.ChannelType_Group)
			converKey = fmt.Sprintf("%s_%d", converKey, pbobjs.ChannelType_Group)
		}
		msgFlag := msgdefines.SetStoreMsg(0)
		msgFlag = msgdefines.SetCountMsg(msgFlag)

		var combiner *Combiner
		botUserInfo := commonservices.GetTargetDisplayUserInfo(ctx, botId)
		botInfo.BotEngine.StreamChat(ctx, msg.SenderId, converKey, msg.ChannelType, txtMsg.Content, func(answerPart string, sectionStart, sectionEnd, isEnd bool) {
			if sectionStart {
				curr := time.Now().UnixMilli()
				streamFlag := msgdefines.SetStreamMsg(0)
				streamFlag = msgdefines.SetStateMsg(streamFlag)
				combiner = &Combiner{
					interval: 50,
					StreamMsg: &pbobjs.DownMsg{
						TargetId:       botId,
						ChannelType:    msg.ChannelType,
						MsgType:        "jgs:text",
						SenderId:       botId,
						Flags:          streamFlag,
						ClientUid:      tools.GenerateUUIDShort22(),
						TargetUserInfo: botUserInfo,
						MsgTime:        curr,
						MsgId:          tools.GenerateMsgId(curr, int32(msg.ChannelType), msg.TargetId),
						StreamMsgParts: []*pbobjs.StreamMsgItem{},
					},
				}
			}
			if isEnd || sectionEnd {
				if combiner != nil {
					streamMsg := combiner.StreamMsg
					if msg.ChannelType == pbobjs.ChannelType_Private {
						part := combiner.GetLeft()
						if part != "" {
							streamMsg.MsgSeqNo = combiner.GetSubSeq()
							streamMsg.MsgContent = tools.ToJsonBs(&StreamMsg{
								Content:     part,
								StreamMsgId: streamMsg.MsgId,
							})
							ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, botId)
							MsgDirect(ctx, msg.SenderId, streamMsg)
						}
					}
					finalContent := combiner.GetFinal()
					if finalContent != "" {
						if msg.ChannelType == pbobjs.ChannelType_Private {
							commonservices.SyncPrivateMsgOverUpstream(ctx, botId, msg.SenderId, &pbobjs.UpMsg{
								MsgType: msgdefines.InnerMsgType_Text,
								MsgContent: tools.ToJsonBs(&msgdefines.TextMsg{
									Content: combiner.GetFinal(),
									Extra: tools.ToJson(&StreamMsg{
										StreamMsgId: streamMsg.MsgId,
									}),
								}),
								Flags: msgFlag,
							})
						} else if msg.ChannelType == pbobjs.ChannelType_Group {
							commonservices.SyncGroupMsgOverUpstream(ctx, botId, msg.TargetId, &pbobjs.UpMsg{
								MsgType: msgdefines.InnerMsgType_Text,
								MsgContent: tools.ToJsonBs(&msgdefines.TextMsg{
									Content: combiner.GetFinal(),
									Extra: tools.ToJson(&StreamMsg{
										StreamMsgId: streamMsg.MsgId,
									}),
								}),
								Flags: msgFlag,
							})
						}
					}
					combiner = nil
				}
			} else {
				if combiner != nil {
					part := combiner.Append(answerPart)
					if part != "" {
						if msg.ChannelType == pbobjs.ChannelType_Private {
							streamMsg := combiner.StreamMsg
							streamMsg.MsgContent = tools.ToJsonBs(&StreamMsg{
								Content:     part,
								StreamMsgId: combiner.StreamMsg.MsgId,
							})
							streamMsg.MsgSeqNo = combiner.GetSubSeq()
							ctx = context.WithValue(ctx, bases.CtxKey_RequesterId, botId)
							MsgDirect(ctx, msg.SenderId, streamMsg)
						}
					}
				}
			}
		})
	}
}

type StreamMsg struct {
	Content     string `json:"content,omitempty"`
	StreamMsgId string `json:"stream_msg_id"`
}

func MsgDirect(ctx context.Context, targetId string, downMsg *pbobjs.DownMsg) {
	rpcMsg := bases.CreateServerPubWraper(ctx, bases.GetRequesterIdFromCtx(ctx), targetId, "msg", downMsg)
	if downMsg.IsSend {
		rpcMsg.PublishType = int32(commonservices.PublishType_AllSessionExceptSelf)
	}
	rpcMsg.Qos = 0
	bases.UnicastRouteWithNoSender(rpcMsg)
}
