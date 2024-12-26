package transengines

import (
	"fmt"
	"im-server/commons/logs"
	"im-server/commons/tools"
	"io"
	"net/http"
	"strings"
	"sync"
)

type BdTransEngine struct {
	AppKey      string `json:"-"`
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	accessToken string `json:"-"`
}

func (eng *BdTransEngine) Translate(content string, langs []string) map[string]string {
	result := map[string]string{}
	if len(langs) <= 0 {
		return result
	}
	if eng.accessToken == "" {
		token, err := getBdToken(eng.ApiKey, eng.SecretKey)
		if err == nil && token != "" {
			eng.accessToken = token
		} else {
			logs.Errorf("failed to get  bd translate access token:%v", err)
			eng.accessToken = "NO"
		}
	}
	if eng.accessToken == "NO" {
		logs.Errorf("have no bd translate access token. appkey:%s", eng.AppKey)
		return result
	}
	if len(langs) > 1 {
		wg := &sync.WaitGroup{}
		lock := &sync.RWMutex{}
		for _, lang := range langs {
			wg.Add(1)
			language := lang
			go func() {
				defer wg.Done()
				txtAfterTrans := bdTranslate(language, content, eng.accessToken)
				if txtAfterTrans != "" {
					lock.Lock()
					defer lock.Unlock()
					result[language] = txtAfterTrans
				}
			}()
		}
		wg.Wait()
	} else {
		txtAfterTrans := bdTranslate(langs[0], content, eng.accessToken)
		if txtAfterTrans != "" {
			result[langs[0]] = txtAfterTrans
		}
	}
	return result
}

func bdTranslate(targetLan string, text string, accessToken string) string {
	url := fmt.Sprintf("https://aip.baidubce.com/rpc/2.0/mt/texttrans/v1?access_token=%s", accessToken)
	headers := map[string]string{}
	headers["Content-Type"] = "application/json"
	bs, _, err := tools.HttpDoBytes(http.MethodPost, url, headers, tools.ToJson(BdTransReq{
		From: "auto",
		To:   targetLan,
		Q:    text,
	}))
	if err != nil || len(bs) <= 0 {
		return ""
	}
	resp := &BdTransResp{}
	err = tools.JsonUnMarshal(bs, resp)
	if err != nil || len(resp.Result.TransResultItems) <= 0 {
		return ""
	}
	results := []string{}
	for _, item := range resp.Result.TransResultItems {
		results = append(results, item.Dst)
	}
	return strings.Join(results, "\n")
}

type BdTransReq struct {
	From string `json:"from"`
	To   string `json:"to"`
	Q    string `json:"q"`
}

type BdTransResp struct {
	Result *BdTransRespResult `json:"result"`
}
type BdTransRespResult struct {
	TransResultItems []*BdTransRespResultItem `json:"trans_result"`
}
type BdTransRespResultItem struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func getBdToken(apiKey, secret string) (string, error) {
	url := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?client_id=%s&client_secret=%s&grant_type=client_credentials", apiKey, secret)
	payload := strings.NewReader(``)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resp := &BdAccessTokenResp{}
	err = tools.JsonUnMarshal(body, resp)
	if err != nil || resp.AccessToken == "" {
		return "", err
	}

	return resp.AccessToken, nil
}

type BdAccessTokenResp struct {
	AccessToken string `json:"access_token"`
}
