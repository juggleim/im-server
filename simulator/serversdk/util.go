package serversdk

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"im-server/commons/tools"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type ApiCode int32

var (
	ApiCode_Success          ApiCode = 0
	ApiCode_HttpTimeout      ApiCode = 1
	ApiCode_DecodeFail       ApiCode = 2
	ApiCode_NotSupportMethod ApiCode = 3
)

type ApiResp struct {
	Code ApiCode     `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

func HttpDo(method, url string, header map[string]string, body string) (string, error) {
	bs, err := HttpDoBytes(method, url, header, body)
	return string(bs), err
}

/*
respBs, err := HttpDoBytes("POST", url, headers, string(bodyBs))

	if err != nil {
		return nil, ApiCode_HttpTimeout, "", err
	}
	resp := &ApiResp{
		Data: &UserRegResp{},
	}
	err = json.Unmarshal(respBs, resp)
	if err != nil {
		return nil, ApiCode_DecodeFail, "", err
	}
	if resp.Code != ApiCode_Success {
		return nil, ApiCode(resp.Code), "", fmt.Errorf(resp.Msg)
	}
	if resp.Data == nil {
		return nil, ApiCode_DecodeFail, "", fmt.Errorf("decode fail.")
	}

	return resp.Data.(*UserRegResp), ApiCode_Success, "", nil
*/
func (sdk *JuggleIMSdk) HttpCall(method, url string, req interface{}, resp interface{}) (ApiCode, string, error) {
	traceId := tools.GenerateUUIDShort11()
	headers := sdk.getHeaders()
	var respBs []byte
	var err error
	if method == http.MethodPost {
		bodyBs, _ := json.Marshal(req)
		respBs, err = HttpDoBytes(http.MethodPost, url, headers, string(bodyBs))
		if err != nil {
			return ApiCode_HttpTimeout, traceId, err
		}
	} else if method == http.MethodGet {

	} else {
		return ApiCode_NotSupportMethod, traceId, fmt.Errorf("not support method:%s", method)
	}
	apiResp := &ApiResp{
		Data: resp,
	}
	err = json.Unmarshal(respBs, apiResp)
	if err != nil {
		return ApiCode_DecodeFail, traceId, err
	}
	if apiResp.Code != ApiCode_Success {
		return ApiCode(apiResp.Code), traceId, fmt.Errorf(apiResp.Msg)
	}
	if resp != nil && apiResp.Data == nil {
		return ApiCode_DecodeFail, traceId, fmt.Errorf("decode fail.")
	}
	return ApiCode_Success, traceId, nil
}

func HttpDoBytes(method, url string, header map[string]string, body string) ([]byte, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return []byte{}, err
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}

	resp, err := client.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err == nil && resp != nil && resp.Body != nil {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		return respBody, nil
	}
	return []byte{}, err
}

func (sdk *JuggleIMSdk) getHeaders() map[string]string {
	nonce := fmt.Sprintf("%d", rand.Int31n(10000))
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := SHA1(fmt.Sprintf("%s%s%s", sdk.Secret, nonce, timestamp))

	return map[string]string{
		"Content-Type": "application/json",
		"appkey":       sdk.Appkey,
		"nonce":        nonce,
		"timestamp":    timestamp,
		"signature":    signature,
	}
}
