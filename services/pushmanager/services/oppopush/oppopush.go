package oppopush

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"im-server/services/pushmanager/services/httputil"
)

type OppoPushClient struct {
	httpClient        *http.Client
	host              string
	appKey            string
	masterSecret      string
	authToken         string
	authTokenExpireAt int64
}

func NewOppoPushClient(appKey, masterSecret string) *OppoPushClient {
	return &OppoPushClient{
		host:         Host,
		appKey:       appKey,
		masterSecret: masterSecret,
	}
}

func (c *OppoPushClient) SetHost(host string) {
	c.host = host
}

func (c *OppoPushClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *OppoPushClient) auth(ctx context.Context) (string, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if c.authToken != "" && c.authTokenExpireAt > now {
		return c.authToken, nil
	}

	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	shaByte := sha256.Sum256([]byte(c.appKey + timestamp + c.masterSecret))
	sign := hex.EncodeToString(shaByte[:])

	req := &AuthReq{
		AppKey:    c.appKey,
		Sign:      sign,
		Timestamp: timestamp,
	}
	res := &AuthRes{}

	params := httputil.StructToUrlValues(req)
	code, resBody, err := httputil.PostForm(ctx, c.httpClient, c.host+AuthURL, params, res, nil)
	if err != nil {
		return "", fmt.Errorf("code=%d body=%s err=%v", code, resBody, err)
	}

	if code != http.StatusOK || res.Code != 0 || res.Data.AuthToken == "" {
		return "", fmt.Errorf("code=%d body=%s", code, resBody)
	}

	c.authToken = res.Data.AuthToken
	c.authTokenExpireAt = now + 60*60*1000 // 一个小时后更新
	return c.authToken, nil
}

// Send 单推-通知栏消息推送
func (c *OppoPushClient) Send(req *SendReq) (*SendRes, error) {
	return c.SendWithContext(context.Background(), req)
}

func (c *OppoPushClient) SendWithContext(ctx context.Context, req *SendReq) (*SendRes, error) {
	res := &SendRes{}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, err
	}

	message, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("message", string(message))
	params.Add("auth_token", token)

	code, resBody, err := httputil.PostForm(ctx, c.httpClient, c.host+SendURL, params, res, nil)
	if err != nil {
		return nil, fmt.Errorf("code=%d body=%s err=%v", code, resBody, err)
	}

	if code != http.StatusOK || res.Code != 0 {
		return nil, fmt.Errorf("code=%d body=%s", code, resBody)
	}

	return res, nil
}
