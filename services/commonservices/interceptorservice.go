package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/interceptors"
	"im-server/services/commonservices/msgdefines"
	"time"

	"go.uber.org/atomic"
)

var (
	interceptorCache *caches.LruCache
	interceptorLocks *tools.SegmentatedLocks

	sensitiveInterceptor = &interceptors.MsgInterceptor{
		Interceptor: &interceptors.SensitiveInterceptor{},
	}
)

func init() {
	interceptorCache = caches.NewLruCacheWithAddReadTimeout("interceptor_cache", 1000, nil, 10*time.Minute, 10*time.Minute)
	interceptorLocks = tools.NewSegmentatedLocks(128)
}

func CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, upMsg *pbobjs.UpMsg) (interceptors.InterceptorResult, int64) {
	if msgdefines.IsCmdMsg(upMsg.Flags) || msgdefines.IsStateMsg(upMsg.Flags) {
		return interceptors.InterceptorResult_Pass, 0
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	appInfo, exist := GetAppInfo(appkey)
	if exist && appInfo != nil && appInfo.OpenSensitive {
		result, code := sensitiveInterceptor.CheckMsgInterceptor(ctx, senderId, receiverId, channelType, upMsg)
		if result != interceptors.InterceptorResult_Pass {
			return result, code
		}
	}

	//other
	msgInterceptors := GetMsgInterceptors(appkey)
	for _, msgInterceptor := range msgInterceptors.Interceptors {
		result, code := msgInterceptor.CheckMsgInterceptor(ctx, senderId, receiverId, channelType, upMsg)
		if result != interceptors.InterceptorResult_Pass {
			return result, code
		}
	}
	return interceptors.InterceptorResult_Pass, 0
}

func GetMsgInterceptors(appkey string) *MsgInterceptors {
	if val, exist := interceptorCache.Get(appkey); exist {
		return val.(*MsgInterceptors)
	} else {
		lock := interceptorLocks.GetLocks(appkey)
		lock.Lock()
		defer lock.Unlock()
		if val, exist := interceptorCache.Get(appkey); exist {
			return val.(*MsgInterceptors)
		} else {
			msgInterceptors := LoadInterceptors(appkey)
			interceptorCache.Add(appkey, msgInterceptors)
			return msgInterceptors
		}
	}
}

type MsgInterceptors struct {
	Interceptors []*interceptors.MsgInterceptor
}

func LoadInterceptors(appkey string) *MsgInterceptors {
	ret := &MsgInterceptors{
		Interceptors: []*interceptors.MsgInterceptor{},
	}
	dao := dbs.InterceptorDao{}
	dbInterceptors, err := dao.QryInterceptors(appkey)
	if err == nil {
		for _, dbInterceptor := range dbInterceptors {
			if dbInterceptor.InterceptType == dbs.InterceptorType_Custom {
				appInfo, exist := GetAppInfo(appkey)
				if exist && appInfo != nil {
					ret.Interceptors = append(ret.Interceptors, &interceptors.MsgInterceptor{
						Interceptor: &interceptors.CustomInterceptor{
							AppKey:     dbInterceptor.AppKey,
							AppSecret:  appInfo.AppSecret,
							RequestUrl: dbInterceptor.RequestUrl,
							Conditions: LoadIcConditions(appkey, dbInterceptor.ID),
						},
					})
				}
			} else if dbInterceptor.InterceptType == dbs.InterceptorType_Baidu {
				bdConf := interceptors.BdInterceptorConf{}
				err := tools.JsonUnMarshal([]byte(dbInterceptor.Conf), &bdConf)
				if err == nil && bdConf.ApiKey != "" && bdConf.SecretKey != "" {
					ret.Interceptors = append(ret.Interceptors, &interceptors.MsgInterceptor{
						Interceptor: &interceptors.BdInterceptor{
							Conf:        &bdConf,
							AccessToken: atomic.NewString(""),
							ExpireAt:    atomic.NewInt64(0),
							Conditions:  LoadIcConditions(appkey, dbInterceptor.ID),
						},
					})
				}
			}
		}
	}
	return ret
}

type BdInterceptorConf struct {
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

func LoadIcConditions(appkey string, interceptorId int64) []*interceptors.Condition {
	ret := []*interceptors.Condition{}
	dao := dbs.IcConditionDao{}
	dbConditions, err := dao.QryConditions(appkey, interceptorId)
	if err == nil {
		for _, dbCondition := range dbConditions {
			condition := &interceptors.Condition{
				ChannelTypeChecker: interceptors.CreateMatcher(dbCondition.ChannelType),
				MsgTypeChecker:     interceptors.CreateMatcher(dbCondition.MsgType),
				SenderIdChecker:    interceptors.CreateMatcher(dbCondition.SenderId),
				ReceiverIdChecker:  interceptors.CreateMatcher(dbCondition.ReceiverId),
			}
			ret = append(ret, condition)
		}
	}
	return ret
}
