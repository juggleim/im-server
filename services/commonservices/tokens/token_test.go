package tokens

import (
	"fmt"
	"testing"
	"time"
)

func TestToken(t *testing.T) {
	token := ImToken{
		AppKey:    "appkey3",
		UserId:    "1793842846692876290",
		DeviceId:  "deviceid",
		TokenTime: time.Now().UnixMilli(),
	}
	secureKey := []byte("abcdefghijk1mn0p")
	tokenStr, err := token.ToTokenString(secureKey)
	fmt.Println(tokenStr, err)
	if err == nil {
		tokenWrap, err := ParseTokenString(tokenStr)
		if err == nil {
			newToken, err := ParseToken(tokenWrap, secureKey)
			if err == nil {
				if newToken.UserId == token.UserId {
					return
				}
			}
		}
	}
	t.Error("Failed.")
}
