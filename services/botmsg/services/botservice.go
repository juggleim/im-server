package services

import (
	"bytes"
	"context"
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
		// create stream msg
		msgFlag := msgdefines.SetStoreMsg(0)
		msgFlag = msgdefines.SetCountMsg(msgFlag)
		msgFlag = msgdefines.SetStreamMsg(msgFlag)
		ctx = bases.SetRequesterId2Ctx(ctx, botId)
		shellContent := &msgdefines.StreamMsg{
			Content: "",
		}
		bs, _ := tools.JsonMarshal(shellContent)
		_, msgId, _, _ := commonservices.SyncPrivateMsgOverUpstream(ctx, botId, msg.SenderId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_StreamText,
			MsgContent: bs,
			Flags:      msgFlag,
		})
		var ts int64 = 0
		buf := bytes.NewBuffer([]byte{})
		botInfo.BotEngine.StreamChat(ctx, msg.SenderId, "", txtMsg.Content, func(answerPart string, isEnd bool) {
			curr := time.Now().UnixMilli()
			if ts == 0 {
				ts = curr
			}
			if !isEnd {
				buf.WriteString(answerPart)
				if curr-ts > 50 {
					partBs, _ := tools.JsonMarshal(&msgdefines.StreamMsg{
						Content: buf.String(),
					})
					bases.SyncRpcCall(ctx, "pri_stream", msg.SenderId, &pbobjs.StreamDownMsg{
						TargetId:    msg.SenderId,
						ChannelType: msg.ChannelType,
						MsgId:       msgId,
						MsgItems: []*pbobjs.StreamMsgItem{
							{
								Event:          pbobjs.StreamEvent_StreamMessage,
								PartialContent: partBs,
							},
						},
						MsgType: msgdefines.InnerMsgType_StreamText,
					}, nil)
					buf = bytes.NewBuffer([]byte{})
					ts = curr
				}
			} else {
				lastPart := buf.String()
				items := []*pbobjs.StreamMsgItem{}
				if lastPart != "" {
					partBs, _ := tools.JsonMarshal(&msgdefines.StreamMsg{
						Content: lastPart,
					})
					items = append(items, &pbobjs.StreamMsgItem{
						Event:          pbobjs.StreamEvent_StreamMessage,
						PartialContent: partBs,
					})
				}
				items = append(items, &pbobjs.StreamMsgItem{
					Event: pbobjs.StreamEvent_StreamComplete,
				})
				bases.SyncRpcCall(ctx, "pri_stream", msg.SenderId, &pbobjs.StreamDownMsg{
					TargetId:    msg.SenderId,
					ChannelType: msg.ChannelType,
					MsgId:       msgId,
					MsgItems:    items,
					MsgType:     msgdefines.InnerMsgType_StreamText,
				}, nil)
			}
		})
	}
}
