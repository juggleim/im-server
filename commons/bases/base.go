package bases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"im-server/commons/configures"
	"im-server/commons/errs"
	"im-server/commons/gmicro"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"

	"google.golang.org/protobuf/proto"
)

const (
	NodeTag_Nav     = "nav"
	NodeTag_Api     = "api"
	NodeTag_Connect = "connect"
	NodeTag_Admin   = "admin"
)

type HttpNodeExt struct {
	Port int `json:"port"`
}
type ConnectNodeExt struct {
	WsPort int `json:"wsport"`
}

var cluster *gmicro.Cluster

func InitImServer(exts map[string]string) error {
	cluster = gmicro.NewCluster(configures.Config.NodeName, configures.Config.NodeHost, exts)
	return nil
}

type IRoute interface {
	GetMethod() string
	GetTargetId() string
}

func GetCluster() *gmicro.Cluster {
	return cluster
}

type ApiCallbackActor struct {
	actorsystem.UntypedActor
	respChan chan *ApiRespWraper
}
type ApiRespWraper struct {
	Msg *pbobjs.RpcMessageWraper
	Err error
}

func (actor *ApiCallbackActor) OnReceive(ctx context.Context, input proto.Message) {
	if rpcMsg, ok := input.(*pbobjs.RpcMessageWraper); ok {
		actor.respChan <- &ApiRespWraper{
			Msg: rpcMsg,
			Err: nil,
		}
	} else {
		fmt.Println("need log.")
	}
}
func (actor *ApiCallbackActor) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}
func (actor *ApiCallbackActor) OnTimeout() {
	actor.respChan <- &ApiRespWraper{
		Msg: nil,
		Err: errors.New("time out1"),
	}
}

func SyncUnicastRoute(msg IRoute, ttl time.Duration) (*pbobjs.RpcMessageWraper, error) {
	respChan := make(chan *ApiRespWraper, 1)
	sender := cluster.CallbackActorOf(ttl, &ApiCallbackActor{
		respChan: respChan,
	})
	cluster.UnicastRoute(msg.GetMethod(), msg.GetTargetId(), msg.(*pbobjs.RpcMessageWraper), sender)

	select {
	case resp := <-respChan:
		return resp.Msg, resp.Err
	case <-time.After(ttl + time.Millisecond*1000):
		return nil, errors.New("time out2")
	}
}

func SyncRpcCall(ctx context.Context, method, targetId string, req proto.Message, respFactory func() proto.Message, opts ...BaseActorOption) (errs.IMErrorCode, proto.Message, error) {
	result, err := SyncOriginalRpcCall(ctx, method, targetId, req, opts...)
	if err != nil {
		return errs.IMErrorCode_DEFAULT, nil, err
	}

	code := errs.IMErrorCode(result.ResultCode)
	if respFactory == nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil, nil
	}

	respObj := respFactory()
	err = tools.PbUnMarshal(result.AppDataBytes, respObj)
	if err != nil {
		return errs.IMErrorCode_DEFAULT, nil, err
	}
	return code, respObj, nil
}

func SyncOriginalRpcCall(ctx context.Context, method, targetId string, req proto.Message, opts ...BaseActorOption) (*pbobjs.RpcMessageWraper, error) {
	if len(opts) > 0 {
		for _, opt := range opts {
			ctx = opt.HandleCtx(ctx)
		}
	}
	dataBytes, _ := tools.PbMarshal(req)
	result, err := SyncUnicastRoute(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_QueryMsg,
		AppKey:       GetAppKeyFromCtx(ctx),
		Session:      GetSessionFromCtx(ctx),
		Method:       method,
		RequesterId:  GetRequesterIdFromCtx(ctx),
		Qos:          GetQosFromCtx(ctx),
		ReqIndex:     GetSeqIndexFromCtx(ctx),
		TargetId:     targetId,
		AppDataBytes: dataBytes,
		OnlySendbox:  GetOnlySendboxFromCtx(ctx),
		NoSendbox:    GetNoSendboxFromCtx(ctx),
		IsFromApi:    GetIsFromApiFromCtx(ctx),
		IsFromApp:    GetIsFromAppFromCtx(ctx),
		TargetIds:    GetTargetIdsFromCtx(ctx),
		ExtParams:    GetExtsFromCtx(ctx),
		MsgId:        GetMsgIdFromCtx(ctx),
		DelMsgId:     GetDelMsgIdFromCtx(ctx),
	}, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AsyncRpcCall(ctx context.Context, method, targetId string, req proto.Message, opts ...BaseActorOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			ctx = opt.HandleCtx(ctx)
		}
	}
	dataBytes, _ := tools.PbMarshal(req)
	UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
		AppKey:       GetAppKeyFromCtx(ctx),
		Session:      GetSessionFromCtx(ctx),
		Method:       method,
		RequesterId:  GetRequesterIdFromCtx(ctx),
		ReqIndex:     GetSeqIndexFromCtx(ctx),
		Qos:          GetQosFromCtx(ctx),
		AppDataBytes: dataBytes,
		TargetId:     targetId,
		OnlySendbox:  GetOnlySendboxFromCtx(ctx),
		NoSendbox:    GetNoSendboxFromCtx(ctx),
		IsFromApi:    GetIsFromApiFromCtx(ctx),
		IsFromApp:    GetIsFromAppFromCtx(ctx),
		TargetIds:    GetTargetIdsFromCtx(ctx),
		ExtParams:    GetExtsFromCtx(ctx),
		MsgId:        GetMsgIdFromCtx(ctx),
		DelMsgId:     GetDelMsgIdFromCtx(ctx),
	})
}

func AsyncRpcCallWithSender(ctx context.Context, method, targetId string, req proto.Message, sender actorsystem.ActorRef, opts ...BaseActorOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			ctx = opt.HandleCtx(ctx)
		}
	}
	dataBytes, _ := tools.PbMarshal(req)
	UnicastRouteWithSenderActor(&pbobjs.RpcMessageWraper{
		RpcMsgType:   GetRpcTypeFromCtx(ctx),
		AppKey:       GetAppKeyFromCtx(ctx),
		Session:      GetSessionFromCtx(ctx),
		Method:       method,
		RequesterId:  GetRequesterIdFromCtx(ctx),
		ReqIndex:     GetSeqIndexFromCtx(ctx),
		Qos:          GetQosFromCtx(ctx),
		AppDataBytes: dataBytes,
		TargetId:     targetId,
		OnlySendbox:  GetOnlySendboxFromCtx(ctx),
		NoSendbox:    GetNoSendboxFromCtx(ctx),
		IsFromApi:    GetIsFromApiFromCtx(ctx),
		IsFromApp:    GetIsFromAppFromCtx(ctx),
		TargetIds:    GetTargetIdsFromCtx(ctx),
		ExtParams:    GetExtsFromCtx(ctx),
		MsgId:        GetMsgIdFromCtx(ctx),
		DelMsgId:     GetDelMsgIdFromCtx(ctx),
	}, sender)
}

func GroupRpcCall(ctx context.Context, method string, targetIds []string, req proto.Message, opts ...BaseActorOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			ctx = opt.HandleCtx(ctx)
		}
	}
	dataBytes, _ := tools.PbMarshal(req)
	disMap := GroupTargets(method, targetIds)
	for _, ids := range disMap {
		if len(ids) > 0 {
			UnicastRouteWithNoSender(&pbobjs.RpcMessageWraper{
				RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
				AppKey:       GetAppKeyFromCtx(ctx),
				Session:      GetSessionFromCtx(ctx),
				Method:       method,
				RequesterId:  GetRequesterIdFromCtx(ctx),
				ReqIndex:     GetSeqIndexFromCtx(ctx),
				Qos:          GetQosFromCtx(ctx),
				AppDataBytes: dataBytes,
				TargetId:     ids[0],
				OnlySendbox:  GetOnlySendboxFromCtx(ctx),
				NoSendbox:    GetNoSendboxFromCtx(ctx),
				IsFromApi:    GetIsFromApiFromCtx(ctx),
				IsFromApp:    GetIsFromAppFromCtx(ctx),
				TargetIds:    ids,
				ExtParams:    GetExtsFromCtx(ctx),
				MsgId:        GetMsgIdFromCtx(ctx),
				DelMsgId:     GetDelMsgIdFromCtx(ctx),
			})
		}
	}
}

func Broadcast(ctx context.Context, method string, req proto.Message, opts ...BaseActorOption) {
	if len(opts) > 0 {
		for _, opt := range opts {
			ctx = opt.HandleCtx(ctx)
		}
	}
	dataBytes, _ := tools.PbMarshal(req)
	BroadcastRouteWithNoSender(&pbobjs.RpcMessageWraper{
		RpcMsgType:   pbobjs.RpcMsgType_ServerPub,
		AppKey:       GetAppKeyFromCtx(ctx),
		Session:      GetSessionFromCtx(ctx),
		Method:       method,
		RequesterId:  GetRequesterIdFromCtx(ctx),
		ReqIndex:     GetSeqIndexFromCtx(ctx),
		Qos:          GetQosFromCtx(ctx),
		AppDataBytes: dataBytes,
		OnlySendbox:  GetOnlySendboxFromCtx(ctx),
		NoSendbox:    GetNoSendboxFromCtx(ctx),
		IsFromApi:    GetIsFromApiFromCtx(ctx),
		IsFromApp:    GetIsFromAppFromCtx(ctx),
		TargetIds:    GetTargetIdsFromCtx(ctx),
		ExtParams:    GetExtsFromCtx(ctx),
		MsgId:        GetMsgIdFromCtx(ctx),
		DelMsgId:     GetDelMsgIdFromCtx(ctx),
	})
}

func UnicastRouteWithCallback(msg IRoute, callbackActor actorsystem.ICallbackUntypedActor, ttl time.Duration) {
	sender := cluster.CallbackActorOf(ttl, callbackActor)
	cluster.UnicastRoute(msg.GetMethod(), msg.GetTargetId(), msg.(*pbobjs.RpcMessageWraper), sender)
}

func UnicastRoute(msg IRoute, sendMethod string) bool {
	sender := cluster.LocalActorOf(sendMethod)
	return cluster.UnicastRoute(msg.GetMethod(), msg.GetTargetId(), msg.(*pbobjs.RpcMessageWraper), sender)
}

func UnicastRouteWithNoSender(msg IRoute) {
	cluster.UnicastRouteWithNoSender(msg.GetMethod(), msg.GetTargetId(), msg.(*pbobjs.RpcMessageWraper))
}

func BroadcastRouteWithNoSender(msg IRoute) {
	cluster.BroadcastWithNoSender(msg.GetMethod(), msg.(*pbobjs.RpcMessageWraper))
}

func UnicastRouteWithSenderActor(msg IRoute, sender actorsystem.ActorRef) bool {
	return cluster.UnicastRoute(msg.GetMethod(), msg.GetTargetId(), msg.(*pbobjs.RpcMessageWraper), sender)
}

func GroupTargets(method string, targetIds []string) map[string][]string {
	rets := map[string][]string{}
	for _, targetId := range targetIds {
		node := cluster.GetTargetNode(method, targetId)
		if node != nil {
			var arr []string
			var ok bool
			if arr, ok = rets[node.Name]; ok {
				arr = append(arr, targetId)
			} else {
				arr = []string{targetId}
			}
			rets[node.Name] = arr
		}
	}
	return rets
}

func Startup() {
	if cluster != nil {
		cluster.Startup()
	}
}

func Shutdown() {

}
