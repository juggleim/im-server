package xiaomipush

import (
	"context"
	"fmt"
	"im-server/services/pushmanager/services/httputil"
	"net/http"
)

type XiaomiPushClient struct {
	httpClient *http.Client
	host       string
	appSecret  string
}

func NewXiaomiPushClient(appSecret string) *XiaomiPushClient {
	return &XiaomiPushClient{
		host:      Host,
		appSecret: appSecret,
	}
}

func (c *XiaomiPushClient) SetHost(host string) {
	c.host = host
}

func (c *XiaomiPushClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *XiaomiPushClient) Send(req *SendReq) (*SendRes, error) {
	return c.SendWithContext(context.Background(), req)
}

func (c *XiaomiPushClient) SendWithContext(ctx context.Context, req *SendReq) (*SendRes, error) {
	res := &SendRes{}

	params := httputil.StructToUrlValues(req)

	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded;charset=UTF-8",
		"Authorization": fmt.Sprintf("key=%s", c.appSecret),
	}

	code, resBody, err := httputil.PostForm(ctx, c.httpClient, c.host+SendURL, params, res, headers)
	if err != nil {
		return nil, fmt.Errorf("code=%d body=%s err=%v", code, resBody, err)
	}

	if code != http.StatusOK || res.Code != 0 {
		return nil, fmt.Errorf("code=%d body=%s", code, resBody)
	}

	return res, nil
}
