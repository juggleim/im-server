package apis

import (
	"im-server/commons/bases"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/navigator/models"

	"github.com/gin-gonic/gin"
)

func NaviGet(ctx *gin.Context) {
	appKey := ctx.GetString(CtxKey_AppKey)
	userId := ctx.GetString(CtxKey_UserId)

	node := bases.GetCluster().GetTargetNode("connect", userId)
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
		UserId:  userId,
		Servers: servers,
	})
}
