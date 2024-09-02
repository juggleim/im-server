package serversdk

import "net/http"

type Message struct {
}

func (sdk *JuggleIMSdk) SendSystemMsg(msg Message) (ApiCode, string, error) {
	url := sdk.ApiUrl + "apigateway/messages/system/send"
	code, traceId, err := sdk.HttpCall(http.MethodPost, url, msg, nil)
	return code, traceId, err
}
