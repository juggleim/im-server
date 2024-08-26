package examples

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func GetFileCred(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.QryFileCredReq{
			FileType: pbobjs.FileType_File,
		}

		code, resp := wsClient.GetFileCred(req)
		fmt.Println("code:", code, "resp:", tools.ToJson(resp))
	}
}
