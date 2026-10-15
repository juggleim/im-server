package tools

import (
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

var sharedHTTPTransport = func() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	transport.MaxIdleConns = 512
	transport.MaxIdleConnsPerHost = 256
	transport.MaxConnsPerHost = 1024
	transport.IdleConnTimeout = 90 * time.Second
	return transport
}()

func HttpDo(method, url string, header map[string]string, body string) (string, int, error) {
	bs, httpCode, err := HttpDoBytes(method, url, header, body)
	return string(bs), httpCode, err
}

func HttpDoBytes(method, url string, header map[string]string, body string) ([]byte, int, error) {
	return HttpDoBytesWithTimeout(method, url, header, body, 3*time.Second)
}

func HttpDoBytesWithContext(ctx context.Context, method, url string, header map[string]string, body string) ([]byte, int, error) {
	return httpDoBytes(ctx, method, url, header, body, 5*time.Second)
}

func HttpDoBytesWithTimeout(method, url string, header map[string]string, body string, timeout time.Duration) ([]byte, int, error) {
	return httpDoBytes(context.Background(), method, url, header, body, timeout)
}

func httpDoBytes(ctx context.Context, method, url string, header map[string]string, body string, timeout time.Duration) ([]byte, int, error) {
	client := &http.Client{Transport: sharedHTTPTransport, Timeout: timeout}
	request, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return []byte{}, 0, err
	}
	for k, v := range header {
		request.Header.Add(k, v)
	}

	resp, err := client.Do(request)
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
