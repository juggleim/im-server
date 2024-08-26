package apis

import (
	"im-server/services/admingateway/services"
	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
)

type AccessAddress struct {
	Original *OriginalAddress `json:"original"`
	Proxy    *ProxyAddress    `json:"proxy"`
}
type OriginalAddress struct {
	Nav     map[string]string `json:"nav"`
	Api     map[string]string `json:"api"`
	Connect map[string]string `json:"connect"`
}
type ProxyAddress struct {
	Nav     *commonservices.AddressConf `json:"nav"`
	Api     *commonservices.AddressConf `json:"api"`
	Connect *commonservices.AddressConf `json:"connect"`
}

func GetAccessAddress(ctx *gin.Context) {
	services.SuccessHttpResp(ctx, &AccessAddress{
		Original: &OriginalAddress{
			Nav:     commonservices.GetOriginalNavAddress().NodeConfs,
			Api:     commonservices.GetOriginalApiAddress().NodeConfs,
			Connect: commonservices.GetOriginalConnectAddress().NodeConfs,
		},
		Proxy: &ProxyAddress{
			Nav:     commonservices.GetProxyNavAddress(),
			Api:     commonservices.GetProxyApiAddress(),
			Connect: commonservices.GetProxyConnectAddress(),
		},
	})
}
