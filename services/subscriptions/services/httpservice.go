package services

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func HttpDoBytes(method, url string, header map[string]string, body string) (int, []byte, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return 1, []byte{}, err
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
			return 2, []byte{}, err
		}
		return resp.StatusCode, respBody, nil
	}
	return 3, nil, err
}
