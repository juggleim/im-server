package tools

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpDo(method, url string, header map[string]string, body string) (string, int, error) {
	bs, httpCode, err := HttpDoBytes(method, url, header, body)
	return string(bs), httpCode, err
}

func HttpDoBytes(method, url string, header map[string]string, body string) ([]byte, int, error) {
	return HttpDoBytesWithTimeout(method, url, header, body, 5*time.Second)
}

func HttpDoBytesWithTimeout(method, url string, header map[string]string, body string, timeout time.Duration) ([]byte, int, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout,
	}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
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
