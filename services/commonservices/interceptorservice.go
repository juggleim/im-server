package commonservices

import (
	"context"
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"im-server/services/sensitivemanager/interceptor"
	"im-server/services/sensitivemanager/interceptor/adapters/local"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	interceptorCache *caches.LruCache
	interceptorLocks *tools.SegmentatedLocks

	defaultLocalInterceptor = &local.LocalInterceptor{}
)

func init() {
	interceptorCache = caches.NewLruCacheWithAddReadTimeout(1000, nil, 10*time.Minute, 10*time.Minute)
	interceptorLocks = tools.NewSegmentatedLocks(128)
}

func CheckMsgInterceptor(ctx context.Context, senderId, receiverId string, channelType pbobjs.ChannelType, upMsg *pbobjs.UpMsg) errs.IMErrorCode {
	intercept, _ := defaultLocalInterceptor.CheckMsgInterceptor(ctx, upMsg)
	if intercept {
		return errs.IMErrorCode_MSG_Hit_Sensitive
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	msgInterceptors := GetMsgInterceptors(appkey)
	for _, msgInterceptor := range msgInterceptors.Interceptors {
		if msgInterceptor.Match(senderId, receiverId, channelType, upMsg.MsgType) {
			if msgInterceptor.interceptor != nil {
				intercept, _ := msgInterceptor.interceptor.CheckMsgInterceptor(ctx, upMsg)
				if intercept {
					return errs.IMErrorCode_MSG_Hit_Sensitive
				}
			}
		}
	}
	return errs.IMErrorCode_SUCCESS
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
	Interceptors []*MsgInterceptor
}

type MsgInterceptor struct {
	Id              int64
	Name            string
	Sort            int
	RequestUrl      string
	RequestTemplate string
	SuccTemplate    string
	IsAsync         bool
	AppKey          string
	Conf            string
	interceptor     interceptor.Interceptor
	Conditions      []*Condition
}

func (interceptor *MsgInterceptor) Match(senderId, receiverId string, channelType pbobjs.ChannelType, msgType string) bool {
	for _, condition := range interceptor.Conditions {
		if condition.ChannelTypeChecker.Match(strconv.Itoa(int(channelType))) &&
			condition.MsgTypeChecker.Match(msgType) &&
			condition.SenderIdChecker.Match(senderId) &&
			condition.ReceiverIdChecker.Match(receiverId) {
			return true
		}
	}
	return false
}

type Condition struct {
	ChannelTypeChecker Matcher
	MsgTypeChecker     Matcher
	SenderIdChecker    Matcher
	ReceiverIdChecker  Matcher
}

func LoadInterceptors(appkey string) *MsgInterceptors {
	ret := &MsgInterceptors{
		Interceptors: []*MsgInterceptor{},
	}
	dao := dbs.InterceptorDao{}
	dbInterceptors, err := dao.QryInterceptors(appkey)
	if err == nil {
		for _, dbInterceptor := range dbInterceptors {
			interceptorHandler, _ := interceptor.BuildInterceptor(dbInterceptor.Name, dbInterceptor.Conf, dbInterceptor.InterceptType != 0)
			msgInterceptor := &MsgInterceptor{
				Id:              dbInterceptor.ID,
				Name:            dbInterceptor.Name,
				Sort:            dbInterceptor.Sort,
				RequestUrl:      dbInterceptor.RequestUrl,
				RequestTemplate: dbInterceptor.RequestTemplate,
				SuccTemplate:    dbInterceptor.SuccTemplate,
				IsAsync:         dbInterceptor.IsAsync > 0,
				AppKey:          dbInterceptor.AppKey,
				Conf:            dbInterceptor.Conf,
				Conditions:      LoadIcConditions(appkey, dbInterceptor.ID),
				interceptor:     interceptorHandler,
			}
			ret.Interceptors = append(ret.Interceptors, msgInterceptor)
		}
	}
	return ret
}

func LoadIcConditions(appkey string, interceptorId int64) []*Condition {
	ret := []*Condition{}
	dao := dbs.IcConditionDao{}
	dbConditions, err := dao.QryConditions(appkey, interceptorId)
	if err == nil {
		for _, dbCondition := range dbConditions {
			condition := &Condition{
				ChannelTypeChecker: CreateMatcher(dbCondition.ChannelType),
				MsgTypeChecker:     CreateMatcher(dbCondition.MsgType),
				SenderIdChecker:    CreateMatcher(dbCondition.SenderId),
				ReceiverIdChecker:  CreateMatcher(dbCondition.ReceiverId),
			}
			ret = append(ret, condition)
		}
	}
	return ret
}

func CreateMatcher(val string) Matcher {
	if val == "" || val == "*" {
		return &NilMatcher{}
	} else if strings.Contains(val, "contains") {
		values, err := extractContainsValues(val)
		if err != nil {
			return &NilMatcher{}
		}
		return NewContainsChecker(values)
	} else {
		return &EqualMatcher{
			value: val,
		}
	}
}

func extractContainsValues(input string) ([]string, error) {
	re := regexp.MustCompile(`contains\(([^)]+)\)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no matches found")
	}
	values := strings.Split(matches[1], ",")

	return values, nil
}

type Matcher interface {
	Match(val string) bool
}

type NilMatcher struct {
}

func (checker *NilMatcher) Match(val string) bool {
	return true
}

type EqualMatcher struct {
	value string
}

func NewEqualChecker(val string) *EqualMatcher {
	return &EqualMatcher{
		value: val,
	}
}

func (checker *EqualMatcher) Match(val string) bool {
	return checker.value == val
}

type ContainsMatcher struct {
	values map[string]struct{}
}

func NewContainsChecker(vals []string) *ContainsMatcher {
	m := &ContainsMatcher{
		values: make(map[string]struct{}, len(vals)),
	}
	for _, val := range vals {
		m.values[val] = struct{}{}
	}
	return m
}

func (checker *ContainsMatcher) Match(val string) bool {
	if _, ok := checker.values[val]; ok {
		return true
	}
	return false
}
