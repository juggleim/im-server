package apimodels

type CommonError struct {
	Code     int    `json:"code"`
	ErrorMsg string `json:"msg"`
}

type CommonResp struct {
	CommonError
	Data interface{} `json:"data,omitempty"`
}
