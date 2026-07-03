package getuipush

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	neturl "net/url"
	"strconv"
	"time"

	"im-server/services/pushmanager/services/httputil"
)

type GetuiPushClient struct {
	httpClient   *http.Client
	host         string
	appID        string
	appKey       string
	masterSecret string

	token         string
	tokenExpireAt int64 // ms
}

func NewGetuiPushClient(appID, appKey, masterSecret string) *GetuiPushClient {
	return &GetuiPushClient{
		host:         Host,
		appID:        appID,
		appKey:       appKey,
		masterSecret: masterSecret,
	}
}

func (c *GetuiPushClient) SetHost(host string) {
	c.host = host
}

func (c *GetuiPushClient) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

func (c *GetuiPushClient) auth(ctx context.Context) (string, error) {
	nowMs := time.Now().UnixMilli()
	if c.token != "" && c.tokenExpireAt > nowMs {
		return c.token, nil
	}

	timestamp := strconv.FormatInt(nowMs, 10)
	sum := sha256.Sum256([]byte(c.appKey + timestamp + c.masterSecret))
	sign := hex.EncodeToString(sum[:])

	req := &AuthReq{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    c.appKey,
	}
	res := &AuthResp{}

	endpoint := fmt.Sprintf("%s/%s%s", c.host, neturl.PathEscape(c.appID), AuthURL)
	code, body, err := httputil.PostJSON(ctx, c.httpClient, endpoint, req, res, nil)
	if err != nil {
		return "", fmt.Errorf("getui auth: code=%d body=%s err=%w", code, body, err)
	}
	if code != http.StatusOK || res.Code != 0 || res.Data.Token == "" {
		return "", fmt.Errorf("getui auth failed: code=%d body=%s", code, body)
	}

	expireMs, parseErr := strconv.ParseInt(res.Data.ExpireTime, 10, 64)
	if parseErr != nil || expireMs <= nowMs {
		// 文档默认 1 天有效期，兜底防止解析异常
		expireMs = nowMs + 24*60*60*1000
	}
	// 提前 1 分钟刷新
	expireMs -= 60 * 1000
	if expireMs <= nowMs {
		expireMs = nowMs + 60*1000
	}

	c.token = res.Data.Token
	c.tokenExpireAt = expireMs
	return c.token, nil
}

// ToSingleCID 执行 cid 单推
func (c *GetuiPushClient) ToSingleCID(req *ToSingleCIDReq) (*ToSingleCIDResp, error) {
	return c.ToSingleCIDWithContext(context.Background(), req)
}

func (c *GetuiPushClient) ToSingleCIDWithContext(ctx context.Context, req *ToSingleCIDReq) (*ToSingleCIDResp, error) {
	if req == nil || req.Audience == nil || len(req.Audience.CID) == 0 {
		return nil, fmt.Errorf("getui push: empty audience.cid")
	}
	if req.PushMessage == nil {
		return nil, fmt.Errorf("getui push: empty push_message")
	}

	token, err := c.auth(ctx)
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("%s/%s%s", c.host, neturl.PathEscape(c.appID), PushSingleCID)
	res := &ToSingleCIDResp{}
	headers := map[string]string{"token": token}

	code, body, err := httputil.PostJSON(ctx, c.httpClient, endpoint, req, res, headers)
	if err != nil {
		return nil, fmt.Errorf("getui push single cid: code=%d body=%s err=%w", code, body, err)
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("getui push single cid http=%d body=%s", code, body)
	}
	if res.Code == 10001 {
		// token 失效：被动刷新后重试一次
		c.token = ""
		c.tokenExpireAt = 0
		return c.retryToSingleCID(ctx, req)
	}
	if res.Code != 0 {
		return nil, fmt.Errorf("getui push single cid code=%d msg=%s body=%s", res.Code, res.Msg, body)
	}
	return res, nil
}

func (c *GetuiPushClient) retryToSingleCID(ctx context.Context, req *ToSingleCIDReq) (*ToSingleCIDResp, error) {
	token, err := c.auth(ctx)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("%s/%s%s", c.host, neturl.PathEscape(c.appID), PushSingleCID)
	res := &ToSingleCIDResp{}
	headers := map[string]string{"token": token}

	code, body, err := httputil.PostJSON(ctx, c.httpClient, endpoint, req, res, headers)
	if err != nil {
		return nil, fmt.Errorf("getui push single cid retry: code=%d body=%s err=%w", code, body, err)
	}
	if code != http.StatusOK || res.Code != 0 {
		return nil, fmt.Errorf("getui push single cid retry failed: http=%d code=%d msg=%s body=%s", code, res.Code, res.Msg, body)
	}
	return res, nil
}
