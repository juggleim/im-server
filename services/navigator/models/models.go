package models

type NaviReq struct {
	AppKey string `json:"app_key"`
	Token  string `json:"token"`
}

type NaviResp struct {
	AppKey  string   `json:"app_key"`
	UserId  string   `json:"user_id"`
	Servers []string `json:"servers"`
}

type DefaultResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func DefaultErr(errMsg string) *DefaultResp {
	if errMsg == "" {
		errMsg = "Default Error Msg."
	}
	return &DefaultResp{
		Code: 10000,
		Msg:  errMsg,
	}
}
func Success(obj interface{}) *DefaultResp {
	return &DefaultResp{
		Code: 0,
		Data: obj,
	}
}
