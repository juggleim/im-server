package vivopush

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"im-server/services/pushmanager/services/httputil"
	"net/http"
	"strconv"
	"time"
)

type VivoPushClient struct {
	httpClient        *http.Client
	host              string
	appId             string
	appKey            string
	appSecret         string
	authToken         string
	authTokenExpireAt int64
}

func NewVivoPushClient(appId, appKey, appSecret string) *VivoPushClient {
	return &VivoPushClient{
		host:      Host,
		appId:     appId,
		appKey:    appKey,
		appSecret: appSecret,
	}
}

func (c *VivoPushClient) SetHost(host string) {
	c.host = host
}

func (c *VivoPushClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *VivoPushClient) auth(ctx context.Context) (string, error) {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if c.authToken != "" && c.authTokenExpireAt > now {
		return c.authToken, nil
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(c.appId + c.appKey + strconv.FormatInt(now, 10) + c.appSecret))
	sign := hex.EncodeToString(md5Ctx.Sum(nil))

	req := &AuthReq{
		AppId:     c.appId,
		AppKey:    c.appKey,
		Timestamp: now,
		Sign:      sign,
	}
	res := &AuthRes{}

	code, resBody, err := httputil.PostJSON(ctx, c.httpClient, c.host+AuthURL, req, res, nil)
	if err != nil {
		return "", fmt.Errorf("code=%d body=%s err=%v", code, resBody, err)
	}

	if code != http.StatusOK || res.Result != 0 || res.AuthToken == "" {
		return "", fmt.Errorf("code=%d body=%s", code, resBody)
	}

	c.authToken = res.AuthToken
	c.authTokenExpireAt = now + 60*60*1000 // 一个小时后更新
	return c.authToken, nil
}

func (c *VivoPushClient) Send(req *SendReq) (*SendRes, error) {
	return c.SendWithContext(context.Background(), req)
}

func (c *VivoPushClient) SendWithContext(ctx context.Context, req *SendReq) (*SendRes, error) {
	res := &SendRes{}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, err
	}

	code, resBody, err := httputil.PostJSON(ctx, c.httpClient, c.host+SendURL, req, res, map[string]string{"authToken": token})
	if err != nil {
		return nil, fmt.Errorf("code=%d body=%s err=%v", code, resBody, err)
	}

	if code != http.StatusOK || res.Result != 0 {
		return nil, fmt.Errorf("code=%d body=%s", code, resBody)
	}

	return res, nil
}
