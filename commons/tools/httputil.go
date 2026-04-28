package tools

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

var defaultHTTPTransport = &http.Transport{
	TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	MaxIdleConns:        512,
	MaxIdleConnsPerHost: 128,
	MaxConnsPerHost:     256,
	IdleConnTimeout:     90 * time.Second,
}

var defaultHTTPClient = &http.Client{
	Transport: defaultHTTPTransport,
}

func HttpDo(method, url string, header map[string]string, body string) (string, int, error) {
	bs, httpCode, err := HttpDoBytes(method, url, header, body)
	return string(bs), httpCode, err
}

func HttpDoBytes(method, url string, header map[string]string, body string) ([]byte, int, error) {
	return HttpDoBytesWithTimeout(method, url, header, body, 5*time.Second)
}

func HttpDoBytesWithTimeout(method, url string, header map[string]string, body string, timeout time.Duration) ([]byte, int, error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return []byte{}, 0, err
	}
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		request = request.WithContext(ctx)
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}

	resp, err := defaultHTTPClient.Do(request)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err == nil && resp != nil && resp.Body != nil {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return []byte{}, 0, err
		}
		return respBody, resp.StatusCode, nil
	}
	return []byte{}, 0, err
}
