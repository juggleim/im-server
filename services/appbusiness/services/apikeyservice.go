package services

import (
	"encoding/base64"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"time"
)

func CheckApiKey(apiKey string, appkey, secureKey string) bool {
	bs, err := base64.URLEncoding.DecodeString(apiKey)
	if err != nil {
		return false
	}
	decodedBs, err := tools.AesDecrypt(bs, []byte(secureKey))
	if err != nil {
		return false
	}
	var apikey pbobjs.ApiKey
	err = tools.PbUnMarshal(decodedBs, &apikey)
	if err != nil {
		return false
	}
	if apikey.Appkey != appkey {
		return false
	}
	return true
}

func GenerateApiKey(appkey, secureKey string) (string, error) {
	apikey := &pbobjs.ApiKey{
		Appkey:      appkey,
		CreatedTime: time.Now().UnixMilli(),
	}
	bs, _ := tools.PbMarshal(apikey)
	encodedBs, err := tools.AesEncrypt(bs, []byte(secureKey))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedBs), nil
}
