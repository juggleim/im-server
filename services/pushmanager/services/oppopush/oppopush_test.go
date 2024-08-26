package oppopush_test

import (
	"testing"

	"im-server/services/pushmanager/services/oppopush"
)

var appKey = "f1ff501e580946dea0324f29d1ce69b0"
var masterSecret = "250f5381919046e3b3ce3ddd09e89bd3"
var regId = "OPPO_CN_7a2dd1e9f74d0c93d1e51ae44c5f5243"
var channelId = "message"

func TestSend(t *testing.T) {
	client := oppopush.NewOppoPushClient(appKey, masterSecret)

	sendReq := &oppopush.SendReq{
		Notification: &oppopush.Notification{
			Title:   "test push title",
			Content: "test push content",
			//ChannelID: channelId,
		},
		TargetType:  2,
		TargetValue: regId,
	}
	sendRes, err := client.Send(sendReq)
	t.Log(sendRes, err)
}
