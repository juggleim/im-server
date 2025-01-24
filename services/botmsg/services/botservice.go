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
	tmpbuf   *bytes.Buffer
	finalbuf *bytes.Buffer
	ts       int64
	interval int64
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

func HandleBotMsg(ctx context.Context, msg *pbobjs.DownMsg) {
	if msg.MsgType != msgdefines.InnerMsgType_Text || msg.ChannelType != pbobjs.ChannelType_Private {
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
		converId := ""
		if msg.ChannelType == pbobjs.ChannelType_Private {
			converId = commonservices.GetConversationId(msg.SenderId, botId, pbobjs.ChannelType_Private)
			converId = fmt.Sprintf("%s_%d", converId, pbobjs.ChannelType_Private)
		} else if msg.ChannelType == pbobjs.ChannelType_Group {
			converId = commonservices.GetConversationId(msg.SenderId, msg.TargetId, pbobjs.ChannelType_Group)
			converId = fmt.Sprintf("%s_%d", converId, pbobjs.ChannelType_Group)
		}
		msgFlag := msgdefines.SetStoreMsg(0)
		msgFlag = msgdefines.SetCountMsg(msgFlag)

		var combiner *Combiner
		botInfo.BotEngine.StreamChat(ctx, msg.SenderId, converId, txtMsg.Content, func(answerPart string, sectionStart, sectionEnd, isEnd bool) {
			if sectionStart {
				combiner = &Combiner{
					interval: 50,
				}
			}
			if isEnd || sectionEnd {
				if combiner != nil {
					commonservices.SyncPrivateMsgOverUpstream(ctx, botId, msg.SenderId, &pbobjs.UpMsg{
						MsgType: msgdefines.InnerMsgType_Text,
						MsgContent: tools.ToJsonBs(&msgdefines.TextMsg{
							Content: combiner.GetFinal(),
						}),
						Flags: msgFlag,
					})
					combiner = nil
				}
			} else {
				if combiner != nil {
					combiner.Append(answerPart)
				}
			}
		})
	}
}
