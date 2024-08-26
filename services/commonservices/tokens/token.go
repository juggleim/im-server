package tokens

import (
	"encoding/base64"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
)

type ImToken struct {
	AppKey    string
	UserId    string
	DeviceId  string
	TokenTime int64
}

func (t ImToken) ToTokenString(secureKey []byte) (string, error) {
	tokenValue := &pbobjs.TokenValue{
		UserId:    t.UserId,
		DeviceId:  t.DeviceId,
		TokenTime: t.TokenTime,
	}
	tokenBs, err := tools.PbMarshal(tokenValue)
	if err == nil {
		encryptToken, err := encrypt(tokenBs, secureKey)
		if err == nil {
			tokenWrap := &pbobjs.TokenWrap{
				AppKey:     t.AppKey,
				TokenValue: encryptToken,
			}
			tokenWrapBs, err := tools.PbMarshal(tokenWrap)
			if err == nil {
				bas64TokenStr := base64.URLEncoding.EncodeToString(tokenWrapBs)
				return bas64TokenStr, nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "", err
}

func encrypt(dataBs, secureKeyBs []byte) ([]byte, error) {
	return tools.AesEncrypt(dataBs, secureKeyBs)
}
func decrypt(cryptedData, secureKeyBs []byte) ([]byte, error) {
	return tools.AesDecrypt(cryptedData, secureKeyBs)
}

func ParseTokenString(tokenStr string) (*pbobjs.TokenWrap, error) {
	tokenWrap := &pbobjs.TokenWrap{}
	tokenWrapBs, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		tokenWrapBs, err = base64.StdEncoding.DecodeString(tokenStr)
	}
	if err == nil {
		err = tools.PbUnMarshal(tokenWrapBs, tokenWrap)
	}
	return tokenWrap, err
}
func ParseToken(tokenWrap *pbobjs.TokenWrap, secureKey []byte) (ImToken, error) {
	token := ImToken{
		AppKey: tokenWrap.AppKey,
	}
	cryptedToken := tokenWrap.TokenValue
	tokenBs, err := decrypt(cryptedToken, secureKey)
	if err == nil {
		tokenValue := &pbobjs.TokenValue{}
		err = tools.PbUnMarshal(tokenBs, tokenValue)
		if err == nil {
			token.UserId = tokenValue.UserId
			token.DeviceId = tokenValue.DeviceId
			token.TokenTime = tokenValue.TokenTime
		} else {
			return token, err
		}
	}
	return token, err
}
