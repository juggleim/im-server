package interceptor

import (
	"context"
	"github.com/tidwall/gjson"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/sensitivemanager/interceptor/adapters/baidu"
	"im-server/services/sensitivemanager/interceptor/adapters/local"
)

type Interceptor interface {
	CheckMsgInterceptor(ctx context.Context, upMsg *pbobjs.UpMsg) (intercept bool, err error)
}

func BuildInterceptor(name string, conf string, interceptWhenSensitive bool) (interceptor Interceptor, err error) {
	switch name {
	case "baidu":
		interceptor = baidu.NewInterceptor(gjson.Get(conf, "api_key").String(),
			gjson.Get(conf, "secret_key").String(), interceptWhenSensitive)
		return
	case "local":
		return &local.LocalInterceptor{}, nil
	default:
		return &local.LocalInterceptor{}, nil
	}
}
