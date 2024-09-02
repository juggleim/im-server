package serversdk

import (
	"net/http"
)

type User struct {
	UserId       string            `json:"user_id"`
	Nickname     string            `json:"nickname"`
	UserPortrait string            `json:"user_portrait"`
	ExtFields    map[string]string `json:"ext_fields"`
}
type UserRegResp struct {
	UserId string `json:"user_id"`
	Token  string `json:"token"`
}

func (sdk *JuggleIMSdk) Register(user User) (*UserRegResp, ApiCode, string, error) {
	url := sdk.ApiUrl + "/apigateway/users/register"
	resp := &UserRegResp{}
	code, traceId, err := sdk.HttpCall(http.MethodPost, url, user, resp)
	return resp, code, traceId, err
}
