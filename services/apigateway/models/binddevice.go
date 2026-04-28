package models

type BindDevice struct {
	UserId        string `json:"user_id"`
	DeviceId      string `json:"device_id"`
	Platform      string `json:"platform"`
	DeviceCompany string `json:"device_company"`
	DeviceModel   string `json:"device_model"`
	CreatedTime   int64  `json:"created_time"`
}

type BindDevicesResp struct {
	Items []*BindDevice `json:"items"`
}
