package xiaomipush

import (
	"testing"
)

var appSecret = "jhg/BEwp367Bio0LITi0gw=="
var regId = "2/gai2A4Sc2O18ucw0y4zIHIIrCHs0sbjc5P4tHwLkWnVC8yvjnovfB+DOTk13B8"

func TestSend(t *testing.T) {
	client := NewXiaomiPushClient(appSecret)

	sendReq := &SendReq{
		RestrictedPackageName: "com.wahu.oa",
		Title:                 "test push title",
		Description:           "test push description",
		RegistrationId:        regId,
		Extra: &Extra{
			ChannelId: "119572",
		},
	}
	sendRes, err := client.Send(sendReq)
	t.Log(sendRes, err)
}
