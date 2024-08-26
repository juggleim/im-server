package hwpush

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type AuthClient struct {
	endpoint  string
	appId     string
	appSecret string
	client    *HTTPClient
}
type TokenMsg struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// NewClient creates a instance of the huawei cloud auth client
// It's contained in huawei cloud app and provides service through huawei cloud app
// If AuthUrl is null using default auth url address
func NewAuthClient(conf *Config) (*AuthClient, error) {
	if conf.AppId == "" || conf.AppSecret == "" {
		return nil, errors.New("appId or appSecret is null")
	}

	c, err := NewHTTPClient()
	if err != nil {
		return nil, errors.New("failed to get http client")
	}

	if conf.AuthUrl == "" {
		return nil, errors.New("authUrl can't be empty")
	}

	return &AuthClient{
		endpoint:  conf.AuthUrl,
		appId:     conf.AppId,
		appSecret: conf.AppSecret,
		client:    c,
	}, nil
}

// GetAuthToken gets token from huawei cloud
// the developer can access the app by using this token
func (ac *AuthClient) GetAuthToken(ctx context.Context) (string, error) {
	if ac.appId == "" || ac.appSecret == "" {
		return "", errors.New("appId or appSecret is null")
	}
	body := fmt.Sprintf("grant_type=client_credentials&client_secret=%s&client_id=%s", ac.appSecret, ac.appId)

	request := &PushRequest{
		Method: http.MethodPost,
		URL:    ac.endpoint,
		Body:   []byte(body),
		Header: []HTTPOption{SetHeader("Content-Type", "application/x-www-form-urlencoded")},
	}

	resp, err := ac.client.DoHttpRequest(ctx, request)
	if err != nil {
		return "", err
	}

	var token TokenMsg
	if resp.Status == 200 {
		err = json.Unmarshal(resp.Body, &token)
		if err != nil {
			return "", err
		}
		return token.AccessToken, nil
	}
	return "", fmt.Errorf("get token failed, status: %d, body: %s", resp.Status, string(resp.Body))
}
