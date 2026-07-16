package honorpush

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"time"

	"im-server/services/pushmanager/services/httputil"
)

// HonorPushClient 荣耀推送 REST 客户端（下行消息 sendMessage）。
// appId：应用 ID（控制台）；clientId / clientSecret：OAuth 客户端凭证，通常对应控制台 Client ID / Client Secret。
// 与 commonservices.HonorPushConf 对应关系：AppId→appId，AppKey→clientId，AppSecret→clientSecret（以控制台实际字段为准）。
type HonorPushClient struct {
	httpClient   *http.Client
	iamHost      string
	pushHost     string
	appID        string
	clientID     string
	clientSecret string

	accessToken   string
	tokenExpireAt int64 // unix 毫秒，提前刷新

	// BadgeClass 桌面角标对应的入口 Activity 全类名，荣耀要求设置角标时必须携带该字段
	BadgeClass string
}

func NewHonorPushClient(appID, clientID, clientSecret string) *HonorPushClient {
	return &HonorPushClient{
		iamHost:      IAMHost,
		pushHost:     PushAPIHost,
		appID:        appID,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *HonorPushClient) SetIAMHost(host string) {
	c.iamHost = host
}

func (c *HonorPushClient) SetPushHost(host string) {
	c.pushHost = host
}

func (c *HonorPushClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *HonorPushClient) auth(ctx context.Context) (string, error) {
	nowMs := time.Now().UnixMilli()
	if c.accessToken != "" && c.tokenExpireAt > nowMs {
		return c.accessToken, nil
	}

	form := neturl.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", c.clientID)
	form.Set("client_secret", c.clientSecret)

	var raw json.RawMessage
	code, body, err := httputil.PostForm(ctx, c.httpClient, c.iamHost+AuthTokenPath, form, &raw, nil)
	if err != nil {
		return "", fmt.Errorf("honor auth post: code=%d body=%s err=%w", code, body, err)
	}

	var errObj authErrorResp
	if err := json.Unmarshal(raw, &errObj); err == nil && errObj.Error != "" {
		return "", fmt.Errorf("honor auth error: %s %s", errObj.Error, errObj.ErrorDescription)
	}

	var ok authTokenResp
	if err := json.Unmarshal(raw, &ok); err != nil {
		return "", fmt.Errorf("honor auth decode: body=%s err=%w", string(raw), err)
	}
	if ok.AccessToken == "" {
		return "", fmt.Errorf("honor auth: empty access_token body=%s", string(raw))
	}

	expireMs := nowMs + ok.ExpiresIn*1000
	// 提前 5 分钟刷新，避免边界过期
	if ok.ExpiresIn > 300 {
		expireMs -= 5 * 60 * 1000
	}
	c.accessToken = ok.AccessToken
	c.tokenExpireAt = expireMs
	return c.accessToken, nil
}

// SendMessage 调用下行消息接口：POST /api/v1/{appId}/sendMessage
func (c *HonorPushClient) SendMessage(req *SendMessageReq) (*SendMessageResp, error) {
	return c.SendMessageWithContext(context.Background(), req)
}

// SendMessageWithContext 带 context 的下行消息发送
func (c *HonorPushClient) SendMessageWithContext(ctx context.Context, req *SendMessageReq) (*SendMessageResp, error) {
	if req == nil || len(req.Token) == 0 {
		return nil, fmt.Errorf("honor push: empty token")
	}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, err
	}

	ts := strconv.FormatInt(time.Now().UnixMilli(), 10)
	endpoint := fmt.Sprintf("%s/api/v1/%s/sendMessage", c.pushHost, neturl.PathEscape(c.appID))
	res := &SendMessageResp{}
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"timestamp":     ts,
	}

	code, body, err := httputil.PostJSON(ctx, c.httpClient, endpoint, req, res, headers)
	if err != nil {
		return nil, fmt.Errorf("honor sendMessage: code=%d body=%s err=%w", code, body, err)
	}

	if code != http.StatusOK {
		return nil, fmt.Errorf("honor sendMessage http=%d body=%s", code, body)
	}
	// 业务码以文档为准，常见成功 code 为 0
	if res.Code != 0 {
		return nil, fmt.Errorf("honor sendMessage code=%d message=%s body=%s", res.Code, res.Message, body)
	}

	return res, nil
}
