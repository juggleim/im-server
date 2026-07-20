package benchmark

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type apiClient struct {
	baseURL   string
	appKey    string
	appSecret string
	http      *http.Client
}

type apiResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
}

type registeredUser struct {
	ID    string
	Token string
}

func newAPIClient(config Config) *apiClient {
	return &apiClient{
		baseURL:   strings.TrimRight(config.APIURL, "/"),
		appKey:    config.AppKey,
		appSecret: config.AppSecret,
		http:      &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *apiClient) registerUser(ctx context.Context, userID string) (registeredUser, error) {
	request := struct {
		UserID   string `json:"user_id"`
		Nickname string `json:"nickname"`
	}{UserID: userID, Nickname: userID}
	var response struct {
		UserID string `json:"user_id"`
		Token  string `json:"token"`
	}
	if err := c.post(ctx, "/users/register", request, &response); err != nil {
		return registeredUser{}, err
	}
	if response.Token == "" {
		return registeredUser{}, fmt.Errorf("register user %s returned an empty token", userID)
	}
	return registeredUser{ID: response.UserID, Token: response.Token}, nil
}

func (c *apiClient) createGroup(ctx context.Context, groupID string, memberIDs []string) error {
	request := struct {
		GroupID   string   `json:"group_id"`
		GroupName string   `json:"group_name"`
		MemberIDs []string `json:"member_ids"`
	}{GroupID: groupID, GroupName: "JuggleIM benchmark group", MemberIDs: memberIDs}
	return c.post(ctx, "/groups/add", request, nil)
}

func (c *apiClient) post(ctx context.Context, endpoint string, payload, output any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode %s request: %w", endpoint, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create %s request: %w", endpoint, err)
	}
	nonce, err := randomNonce()
	if err != nil {
		return fmt.Errorf("create request nonce: %w", err)
	}
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("appkey", c.appKey)
	req.Header.Set("nonce", nonce)
	req.Header.Set("timestamp", timestamp)
	req.Header.Set("signature", signature(c.appSecret, nonce, timestamp))

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call %s: %w", endpoint, err)
	}
	defer resp.Body.Close()
	var envelope apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("decode %s response (HTTP %d): %w", endpoint, resp.StatusCode, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 || envelope.Code != 0 {
		return fmt.Errorf("%s failed: HTTP %d, code %d, message %q", endpoint, resp.StatusCode, envelope.Code, envelope.Msg)
	}
	if output != nil && len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, output); err != nil {
			return fmt.Errorf("decode %s data: %w", endpoint, err)
		}
	}
	return nil
}

func signature(secret, nonce, timestamp string) string {
	sum := sha1.Sum([]byte(secret + nonce + timestamp))
	return hex.EncodeToString(sum[:])
}

func randomNonce() (string, error) {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}
