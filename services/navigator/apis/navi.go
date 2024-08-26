package apis

import (
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/tokens"
	"im-server/services/navigator/models"

	"github.com/gin-gonic/gin"
)

func NaviGet(ctx *gin.Context) {
	appKey := ctx.Request.Header.Get("x-appkey")
	tokenStr := ctx.Request.Header.Get("x-token")
	appInfo, exist := commonservices.GetAppInfo(appKey)
	if exist && appInfo != nil {
		tokenWrap, err := tokens.ParseTokenString(tokenStr)
		if err == nil {
			token, err := tokens.ParseToken(tokenWrap, []byte(appInfo.AppSecureKey))
			if err == nil {
				node := bases.GetCluster().GetTargetNode("connect", token.UserId)
				servers := []string{}
				if node != nil {
					connConf := commonservices.GetConnectAddress()
					if address, ok := connConf.NodeConfs[node.Name]; ok {
						servers = append(servers, address)
					}
					servers = append(servers, connConf.Default...)
				}
				tools.SuccessHttpResp(ctx, models.NaviResp{
					AppKey:  appKey,
					UserId:  token.UserId,
					Servers: servers,
				},
				)
				return
			} else {
				tools.ErrorHttpResp(ctx, errs.IMErrorCode_CONNECT_TOKEN_AUTHFAIL)
				return
			}
		} else {
			tools.ErrorHttpResp(ctx, errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL)
			return
		}
	} else {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_CONNECT_APP_NOT_EXISTED)
		return
	}
}
