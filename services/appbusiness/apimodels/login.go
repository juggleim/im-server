package apimodels

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type SmsLoginReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type EmailLoginReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type LoginUserResp struct {
	UserId        string `json:"user_id"`
	Authorization string `json:"authorization"`
	NickName      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	Status        int    `json:"status"`
	ImToken       string `json:"im_token,omitempty"`
}

type QrCode struct {
	Id string `json:"id"`
}
