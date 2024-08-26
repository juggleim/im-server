package baidu

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/atomic"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"net/http"
	"net/url"
	"time"
)

type (
	AccessTokenResult struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		Error       string `json:"error"`
	}
	TextResult struct {
		LogId          int64  `json:"log_id"`
		Conclusion     string `json:"conclusion"`
		ConclusionType int    `json:"conclusionType"`
		Data           []struct {
			Type           int    `json:"type"`
			SubType        int    `json:"subType"`
			Conclusion     string `json:"conclusion"`
			ConclusionType int    `json:"conclusionType"`
			Msg            string `json:"msg"`
			Hits           []struct {
				DatasetName       string      `json:"datasetName"`
				Words             []string    `json:"words"`
				ModelHitPositions [][]float64 `json:"modelHitPositions,omitempty"`
				WordHitPositions  []struct {
					Keyword   string  `json:"keyword"`
					Positions [][]int `json:"positions"`
					Label     string  `json:"label"`
				} `json:"wordHitPositions,omitempty"`
				Probability float64 `json:"probability,omitempty"`
			} `json:"hits"`
		} `json:"data"`
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}
	ImgResult struct {
		LogId          int64  `json:"log_id"`
		Conclusion     string `json:"conclusion"`
		ConclusionType int    `json:"conclusionType"`
		ErrorCode      int    `json:"error_code"`
		ErrorMsg       string `json:"error_msg"`
	}
)

type Interceptor struct {
	apiKey      string
	secretKey   string
	accessToken *atomic.String
	expireAt    *atomic.Int64

	interceptWhenSensitive bool
}

func NewInterceptor(apiKey, secretKey string, interceptWhenSensitive bool) *Interceptor {
	return &Interceptor{
		apiKey:      apiKey,
		secretKey:   secretKey,
		accessToken: atomic.NewString(""),
		expireAt:    atomic.NewInt64(0),

		interceptWhenSensitive: interceptWhenSensitive,
	}
}

func (i *Interceptor) getAccessToken() (accessToken string, err error) {
	if i.accessToken.Load() != "" && i.expireAt.Load() > time.Now().Unix() {
		return i.accessToken.Load(), nil
	}

	uri := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s", i.apiKey, i.secretKey)
	resp, err := http.Post(uri, "", nil)
	if err != nil {
		return "", err
	}

	var result *AccessTokenResult
	result, err = parseResponse(resp, func(result AccessTokenResult) (err error) {
		if result.Error != "" {
			err = fmt.Errorf("failed to get access token, %s", result.Error)
			return
		}
		return
	})

	if err != nil {
		return
	}
	i.accessToken.Store(result.AccessToken)
	i.expireAt.Store(time.Now().Unix() + result.ExpiresIn - 60)

	return result.AccessToken, nil
}

func (i *Interceptor) InterceptText(text string) (interceptText string, containsSensitive bool, err error) {
	var (
		response *http.Response
		result   *TextResult
	)
	uri := "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined?access_token=%s"

	accessToken, err := i.getAccessToken()
	if err != nil {
		err = errors.Wrap(err, "failed to get access token")
		return
	}
	uri = fmt.Sprintf(uri, accessToken)

	response, err = http.PostForm(uri, map[string][]string{
		"text": {text},
	})

	result, err = parseResponse(response, func(res TextResult) (err error) {
		if res.ErrorCode != 0 {
			err = fmt.Errorf("error code: %d, error message: %s", res.ErrorCode, res.ErrorMsg)
			return
		}
		return
	})
	if err != nil {
		return
	}
	for _, item := range result.Data {
		for _, hit := range item.Hits {
			for _, position := range hit.WordHitPositions {
				for _, pos := range position.Positions {
					if len(pos) == 2 {
						t := []rune(text)
						for i := pos[0]; i <= pos[1]; i++ {
							t[i] = '*'
						}
						text = string(t)

						containsSensitive = true
					}
				}
			}
		}
	}
	interceptText = text
	return
}

func (i *Interceptor) InterceptImage(imageUrl string) (intercept bool, err error) {
	var (
		response    *http.Response
		result      *ImgResult
		accessToken string
	)
	uri := "https://aip.baidubce.com/rest/2.0/solution/v1/img_censor/v2/user_definedaccess_token=%s"

	accessToken, err = i.getAccessToken()
	if err != nil {
		err = errors.Wrap(err, "failed to get access token")
		return
	}
	uri = fmt.Sprintf(uri, accessToken)

	response, err = http.PostForm(uri, map[string][]string{
		"imgUrl": {url.QueryEscape(imageUrl)},
	})

	result, err = parseResponse(response, func(res ImgResult) (err error) {
		if res.ErrorCode != 0 {
			err = fmt.Errorf("error code: %d, error message: %s", res.ErrorCode, res.ErrorMsg)
			return
		}
		return
	})
	if err != nil {
		return
	}

	return result.ConclusionType == 1, nil
}

func (i *Interceptor) CheckMsgInterceptor(ctx context.Context, upMsg *pbobjs.UpMsg) (intercept bool, err error) {
	if upMsg.MsgType == "jg:text" {
		var (
			text              string
			containsSensitive bool
		)
		txtMsg := &struct {
			Content string `json:"content"`
		}{}
		err = tools.JsonUnMarshal(upMsg.MsgContent, txtMsg)
		if err != nil {
			return
		}
		text, containsSensitive, err = i.InterceptText(txtMsg.Content)

		if err != nil {
			return
		}
		txtMsg.Content = text
		upMsg.MsgContent, _ = tools.JsonMarshal(txtMsg)
		if containsSensitive && i.interceptWhenSensitive {
			intercept = true
		}
	}
	return
}

func parseResponse[T any](response *http.Response, checkOk func(res T) error) (*T, error) {
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", response.StatusCode)
	}
	var result T
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}
	if err := checkOk(result); err != nil {
		return nil, err
	}

	return &result, nil
}
