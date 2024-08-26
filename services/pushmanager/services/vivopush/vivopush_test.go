package vivopush_test

import (
	"im-server/services/pushmanager/services/vivopush"
	"strconv"
	"testing"
	"time"
)

var appId = "105736684"
var appKey = "6283aa880794282f004be2c8562f5859"
var appSecret = "54ad6a73-b9d8-4f67-9edf-cc9827a35bfc"
var regId = "v2-CQe1xHjP9-ExtXny_EeUsXWsEKmfjEXOPntCDFmRyr_TjbYQwfRevgqK"

func TestSend(t *testing.T) {
	client := vivopush.NewVivoPushClient(appId, appKey, appSecret)

	sendReq := &vivopush.SendReq{
		RegId:          regId,
		NotifyType:     4,
		Title:          "test push title",
		Content:        "test push content",
		TimeToLive:     24 * 60 * 60,
		SkipType:       1,
		NetworkType:    -1,
		Classification: 1,
		RequestId:      strconv.Itoa(int(time.Now().UnixNano())),
	}
	sendRes, err := client.Send(sendReq)
	t.Log(sendRes, err)
}
