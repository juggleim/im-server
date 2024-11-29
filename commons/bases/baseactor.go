package bases

import (
	"context"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/logs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"reflect"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	preProcess func(ctx context.Context, sender actorsystem.ActorRef) bool
)

func SetPreProcess(f func(ctx context.Context, sender actorsystem.ActorRef) bool) {
	preProcess = f
}

func BaseProcessActor(actor actorsystem.IUntypedActor, serviceName string) actorsystem.IUntypedActor {
	tags := map[string]string{
		"service_name": serviceName,
	}
	return &baseProcessActor{
		exeActor: actor,
		Tags:     tags,
	}
}

type baseProcessActor struct {
	Tags     map[string]string
	exeActor actorsystem.IUntypedActor
}

type BaseActor struct {
	actorsystem.UntypedActor
}
type CtxKey string

const (
	CtxKey_Tags          CtxKey = "CtxKey_Tags"
	CtxKey_RpcType       CtxKey = "CtxKey_RpcType"
	CtxKey_SeqIndex      CtxKey = "CtxKey_SeqIndex"
	CtxKey_AppKey        CtxKey = "CtxKey_AppKey"
	CtxKey_Qos           CtxKey = "CtxKey_Qos"
	CtxKey_Session       CtxKey = "CtxKey_Session"
	CtxKey_DeviceId      CtxKey = "CtxKey_DeviceId"
	CtxKey_InstanceId    CtxKey = "CtxKey_InstanceId"
	CtxKey_Platform      CtxKey = "CtxKey_Platform"
	CtxKey_Method        CtxKey = "CtxKey_Method"
	CtxKey_SourceMethod  CtxKey = "CtxKey_SourceMethod"
	CtxKey_RequesterId   CtxKey = "CtxKey_RequesterId"
	CtxKey_TargetId      CtxKey = "CtxKey_TargetId"
	CtxKey_PublishType   CtxKey = "CtxKey_PublishType"
	CtxKey_IsFromApi     CtxKey = "CtxKey_IsFromApi"
	CtxKey_TerminalCount CtxKey = "CtxKey_TerminalCount"
	CtxKey_GroupId       CtxKey = "CtxKey_GroupId"
	CtxKey_TargetIds     CtxKey = "CtxKey_TargetIds"
	CtxKey_SenderInfo    CtxKey = "CtxKey_SenderInfo"
	CtxKey_NoSendbox     CtxKey = "CtxKey_NoSendbox"
	CtxKey_OnlySendbox   CtxKey = "CtxKey_OnlySendbox"
	CtxKey_Exts          CtxKey = "CtxKey_Exts"
	CtxKey_MsgId         CtxKey = "CtxKey_MsgId"

	CtxKey_StartTime CtxKey = "CtxKey_StartTime"
)

func (actor *baseProcessActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}

func (actor *baseProcessActor) OnReceive(ctx context.Context, input proto.Message) {
	ctx = setCtxValue(ctx, CtxKey_Tags, actor.Tags)
	var err error
	if input != nil {
		ssRequest, ok := input.(*pbobjs.RpcMessageWraper)
		if ok {
			startTime := time.Now()
			ctx = setCtxValue(ctx, CtxKey_RpcType, ssRequest.RpcMsgType)
			ctx = setCtxValue(ctx, CtxKey_SeqIndex, ssRequest.ReqIndex)
			ctx = setCtxValue(ctx, CtxKey_AppKey, ssRequest.AppKey)
			ctx = setCtxValue(ctx, CtxKey_Qos, ssRequest.Qos)
			ctx = setCtxValue(ctx, CtxKey_Session, ssRequest.Session)
			ctx = setCtxValue(ctx, CtxKey_DeviceId, ssRequest.DeviceId)
			ctx = setCtxValue(ctx, CtxKey_InstanceId, ssRequest.InstanceId)
			ctx = setCtxValue(ctx, CtxKey_Platform, ssRequest.Platform)
			ctx = setCtxValue(ctx, CtxKey_Method, ssRequest.Method)
			ctx = setCtxValue(ctx, CtxKey_SourceMethod, ssRequest.SourceMethod)
			ctx = setCtxValue(ctx, CtxKey_RequesterId, ssRequest.RequesterId)
			ctx = setCtxValue(ctx, CtxKey_TargetId, ssRequest.TargetId)
			ctx = setCtxValue(ctx, CtxKey_PublishType, ssRequest.PublishType)
			ctx = setCtxValue(ctx, CtxKey_IsFromApi, ssRequest.IsFromApi)
			ctx = setCtxValue(ctx, CtxKey_TerminalCount, ssRequest.TerminalNum)
			ctx = setCtxValue(ctx, CtxKey_NoSendbox, ssRequest.NoSendbox)
			ctx = setCtxValue(ctx, CtxKey_OnlySendbox, ssRequest.OnlySendbox)

			ctx = setCtxValue(ctx, CtxKey_GroupId, ssRequest.GroupId)
			ctx = setCtxValue(ctx, CtxKey_TargetIds, ssRequest.TargetIds)
			exts := ssRequest.ExtParams
			if exts == nil {
				exts = map[string]string{}
			}
			ctx = setCtxValue(ctx, CtxKey_Exts, exts)
			ctx = setCtxValue(ctx, CtxKey_MsgId, ssRequest.MsgId)
			ctx = setCtxValue(ctx, CtxKey_StartTime, startTime.UnixMilli())

			ctx = setCtxValue(ctx, CtxKey_SenderInfo, ssRequest.SenderInfo)
			if preProcess != nil {
				isContinue := preProcess(ctx, actor.GetSender())
				if !isContinue {
					return
				}
			}
			msgBytes := ssRequest.AppDataBytes
			createInputHandler, ok := actor.exeActor.(actorsystem.ICreateInputHandler)
			if ok {
				msg := createInputHandler.CreateInputObj()
				err = tools.PbUnMarshal(msgBytes, msg)
				if err == nil {
					receiveHandler, ok := actor.exeActor.(actorsystem.IReceiveHandler)
					if ok {
						receiveHandler.OnReceive(ctx, msg)
					}
				} else {
					logs.Errorf("decode input failed. err:%v\tmethod:%s\tsession:%s\tseq_index:%d\tinput_type:%v", err, ssRequest.Method, ssRequest.Session, ssRequest.ReqIndex, reflect.TypeOf(msg))
				}
			}
			consume := time.Since(startTime)
			if consume > 1*time.Second {
				logs.Errorf("RPC_Timeout:%v\tsession:%s\tseq_index:%d\tmethod:%s\trequest_id:%s\ttarget_id:%s", consume, ssRequest.Session, ssRequest.ReqIndex, ssRequest.Method, ssRequest.RequesterId, ssRequest.TargetId)
			}
		}
	}
}

func (actor *baseProcessActor) SetSender(sender actorsystem.ActorRef) {
	senderHandler, ok := actor.exeActor.(actorsystem.ISenderHandler)
	if ok {
		senderHandler.SetSender(sender)
	}
}

func (actor *baseProcessActor) GetSender() actorsystem.ActorRef {
	senderHandler, ok := actor.exeActor.(actorsystem.ISenderHandler)
	if ok {
		return senderHandler.GetSender()
	}
	return nil
}

func (actor *baseProcessActor) SetSelf(self actorsystem.ActorRef) {
	selfHandler, ok := actor.exeActor.(actorsystem.ISelfHandler)
	if ok {
		selfHandler.SetSelf(self)
	}
}

func (actor *baseProcessActor) GetSelf() actorsystem.ActorRef {
	selfHandler, ok := actor.exeActor.(actorsystem.ISelfHandler)
	if ok {
		return selfHandler.GetSelf()
	}
	return nil
}

func (actor *baseProcessActor) OnTimeout() {
	timeoutHandler, ok := actor.exeActor.(actorsystem.ITimeoutHandler)
	if ok {
		timeoutHandler.OnTimeout()
	}
}

func setCtxValue(ctx context.Context, key CtxKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetRpcTypeFromCtx(ctx context.Context) pbobjs.RpcMsgType {
	if rpcMsgType, ok := ctx.Value(CtxKey_RpcType).(pbobjs.RpcMsgType); ok {
		return rpcMsgType
	}
	return -1
}

func GetSeqIndexFromCtx(ctx context.Context) int32 {
	if seqIndex, ok := ctx.Value(CtxKey_SeqIndex).(int32); ok {
		return seqIndex
	}
	return 0
}

func GetQosFromCtx(ctx context.Context) int32 {
	if qos, ok := ctx.Value(CtxKey_Qos).(int32); ok {
		return qos
	}
	return 0
}

func GetAppKeyFromCtx(ctx context.Context) string {
	if appKey, ok := ctx.Value(CtxKey_AppKey).(string); ok {
		return appKey
	}
	return ""
}
func GetSessionFromCtx(ctx context.Context) string {
	if session, ok := ctx.Value(CtxKey_Session).(string); ok {
		return session
	}
	return ""
}
func GetDeviceIdFromCtx(ctx context.Context) string {
	if deviceId, ok := ctx.Value(CtxKey_DeviceId).(string); ok {
		return deviceId
	}
	return ""
}
func GetInstanceIdFromCtx(ctx context.Context) string {
	if instanceId, ok := ctx.Value(CtxKey_InstanceId).(string); ok {
		return instanceId
	}
	return ""
}
func GetPlatformFromCtx(ctx context.Context) string {
	if platform, ok := ctx.Value(CtxKey_Platform).(string); ok {
		return platform
	}
	return ""
}

func GetRequesterIdFromCtx(ctx context.Context) string {
	if requesterId, ok := ctx.Value(CtxKey_RequesterId).(string); ok {
		return requesterId
	}
	return ""
}

func SetRequesterId2Ctx(ctx context.Context, requestId string) context.Context {
	return setCtxValue(ctx, CtxKey_RequesterId, requestId)
}

func GetTargetIdFromCtx(ctx context.Context) string {
	if targetId, ok := ctx.Value(CtxKey_TargetId).(string); ok {
		return targetId
	}
	return ""
}

func GetMethodFromCtx(ctx context.Context) string {
	if method, ok := ctx.Value(CtxKey_Method).(string); ok {
		return method
	}
	return ""
}

func GetSourceMethodFromCtx(ctx context.Context) string {
	if sourceMethod, ok := ctx.Value(CtxKey_SourceMethod).(string); ok {
		return sourceMethod
	}
	return ""
}

func GetIsFromApiFromCtx(ctx context.Context) bool {
	if isFromApi, ok := ctx.Value(CtxKey_IsFromApi).(bool); ok {
		return isFromApi
	}
	return false
}

func GetNoSendboxFromCtx(ctx context.Context) bool {
	if noSendbox, ok := ctx.Value(CtxKey_NoSendbox).(bool); ok {
		return noSendbox
	}
	return false
}

func GetOnlySendboxFromCtx(ctx context.Context) bool {
	if onlySendbox, ok := ctx.Value(CtxKey_OnlySendbox).(bool); ok {
		return onlySendbox
	}
	return false
}

func SetOnlySendbox2Ctx(ctx context.Context, onlySendbox bool) context.Context {
	return setCtxValue(ctx, CtxKey_OnlySendbox, onlySendbox)
}

func GetPublishTypeFromCtx(ctx context.Context) int32 {
	if publishType, ok := ctx.Value(CtxKey_PublishType).(int32); ok {
		return publishType
	}
	return 0
}

func GetTerminalNumFromCtx(ctx context.Context) int32 {
	if terminalNum, ok := ctx.Value(CtxKey_TerminalCount).(int32); ok {
		return terminalNum
	}
	return 0
}

func GetExpendFromCtx(ctx context.Context) int64 {
	if start, ok := ctx.Value(CtxKey_StartTime).(int64); ok {
		now := time.Now().UnixMilli()
		return now - start
	}
	return 0
}

func GetGroupIdFromCtx(ctx context.Context) string {
	if groupId, ok := ctx.Value(CtxKey_GroupId).(string); ok {
		return groupId
	}
	return ""
}

func GetTargetIdsFromCtx(ctx context.Context) []string {
	if targetIds, ok := ctx.Value(CtxKey_TargetIds).([]string); ok {
		return targetIds
	}
	return []string{}
}

func SetTargetIds2Ctx(ctx context.Context, ids []string) context.Context {
	return setCtxValue(ctx, CtxKey_TargetIds, ids)
}

func GetSenderInfoFromCtx(ctx context.Context) *pbobjs.UserInfo {
	if senderInfo, ok := ctx.Value(CtxKey_SenderInfo).(*pbobjs.UserInfo); ok {
		return senderInfo
	}
	return nil
}

func GetTagsFromCtx(ctx context.Context) map[string]string {
	if tags, ok := ctx.Value(CtxKey_Tags).(map[string]string); ok {
		return tags
	}
	return map[string]string{}
}

func GetExtsFromCtx(ctx context.Context) map[string]string {
	if exts, ok := ctx.Value(CtxKey_Exts).(map[string]string); ok {
		return exts
	}
	return map[string]string{}
}

func GetMsgIdFromCtx(ctx context.Context) string {
	if msgId, ok := ctx.Value(CtxKey_MsgId).(string); ok {
		return msgId
	}
	return ""
}
