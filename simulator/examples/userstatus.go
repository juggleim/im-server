package examples

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SubUsers(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.UserIdsReq{
			UserIds: []string{"userid2", "userid3"},
		}

		code := wsClient.SubUsers(req)
		fmt.Println("code:", code)
	}
}

func UnSubUsers(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.UserIdsReq{
			UserIds: []string{"userid2", "userid3"},
		}
		code := wsClient.UnSubUsers(req)
		fmt.Println("code:", code)
	}
}
