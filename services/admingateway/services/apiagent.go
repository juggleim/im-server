package services

import (
	"fmt"
	"im-server/commons/configures"
	"im-server/commons/tools"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ApiAgent(method, path, body, appkey, secret string) (int, string) {
	ImApiUrl := fmt.Sprintf("http://127.0.0.1:%d", configures.Config.ApiGateway.HttpPort)
	respBs, code, err := tools.HttpDoBytes(method, fmt.Sprintf("%s%s", ImApiUrl, path), getSignatureHeaders(appkey, secret), body)
	if err == nil {
		return code, string(respBs)
	} else {
		return http.StatusBadRequest, ""
	}
}

func getSignatureHeaders(appkey, secret string) map[string]string {
	nonce := fmt.Sprintf("%d", rand.Int31n(10000))
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := tools.SHA1(fmt.Sprintf("%s%s%s", secret, nonce, timestamp))
	headers := map[string]string{
		"Content-Type": "application/json",
		"appkey":       appkey,
		"nonce":        nonce,
		"timestamp":    timestamp,
		"signature":    signature,
	}
	return headers
}
