package transengines

import (
	"fmt"
	"im-server/commons/tools"
	"net/http"
	"strings"
	"sync"
)

type DeeplTransEngine struct {
	AppKey  string `json:"-"`
	AuthKey string `json:"auth_key"`
}

func (eng *DeeplTransEngine) Translate(content string, langs []string) map[string]string {
	result := map[string]string{}
	if len(langs) <= 0 {
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
				txtAfterTrans := deeplTranslage(language, content, eng.AuthKey)
				if txtAfterTrans != "" {
					lock.Lock()
					defer lock.Unlock()
					result[language] = txtAfterTrans
				}
			}()
		}
		wg.Wait()
	} else {
		txtAfterTrans := deeplTranslage(langs[0], content, eng.AuthKey)
		if txtAfterTrans != "" {
			result[langs[0]] = txtAfterTrans
		}
	}
	return result
}

func deeplTranslage(targetLan string, text string, authKey string) string {
	url := "https://api.deepl.com/v2/translate"
	headers := map[string]string{}
	headers["Authorization"] = fmt.Sprintf("DeepL-Auth-Key %s", authKey)
	headers["Content-Type"] = "application/json"
	bs, _, err := tools.HttpDoBytes(http.MethodPost, url, headers, tools.ToJson(DeeplTransReq{
		Texts:      []string{text},
		TargetLang: targetLan,
	}))
	if err != nil || len(bs) <= 0 {
		return ""
	}
	resp := &DeeplTransResp{}
	err = tools.JsonUnMarshal(bs, resp)
	if err != nil || len(resp.Translations) <= 0 {
		return ""
	}
	results := []string{}
	for _, item := range resp.Translations {
		results = append(results, item.Text)
	}
	return strings.Join(results, "")
}

type DeeplTransReq struct {
	Texts      []string `json:"text"`
	TargetLang string   `json:"target_lang"`
}

type DeeplTransResp struct {
	Translations []*DeeplTransRespTransItem `json:"translations"`
}

type DeeplTransRespTransItem struct {
	Text string `json:"text"`
}
