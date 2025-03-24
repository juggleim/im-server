package sms

import (
	"fmt"

	"github.com/baidubce/bce-sdk-go/services/sms"
	"github.com/baidubce/bce-sdk-go/services/sms/api"
)

type BdSmsEngine struct {
	smsClient   *sms.Client
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	Endpoint    string `json:"endpoint"`
	Template    string `json:"template"`
	SignatureId string `json:"signature_id"`
}

func (eng *BdSmsEngine) SmsSend(phone string, params map[string]interface{}) error {
	if eng.smsClient == nil {
		var err error
		eng.smsClient, err = sms.NewClient(eng.ApiKey, eng.SecretKey, eng.Endpoint)
		if err != nil {
			eng.smsClient = nil
			return err
		}
	}
	sendSmsArgs := &api.SendSmsArgs{
		Mobile:      phone,
		Template:    eng.Template,
		SignatureId: eng.SignatureId,
		ContentVar:  params,
	}
	result, err := eng.smsClient.SendSms(sendSmsArgs)
	if err != nil {
		return err
	}
	fmt.Println("result:", result)
	return nil
}
