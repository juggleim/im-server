package examples

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"
)

func SubUsers(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.SubUsersReq{
			UserIds: []string{"userid2", "userid3"},
		}

		code, list := wsClient.SubUsers(req)
		fmt.Println("code:", code)
		if list != nil {
			fmt.Println("UserStatusList items:", len(list.Items))
		}
	}
}

func UnSubUsers(wsClient *wsclients.WsImClient) {
	if wsClient != nil && wsClient.GetState() == utils.State_Connected {
		req := &pbobjs.SubUsersReq{
			UserIds: []string{"userid2", "userid3"},
		}
		code := wsClient.UnSubUsers(req)
		fmt.Println("code:", code)
	}
}
