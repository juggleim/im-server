package tools

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"strings"
)

var defaultEmptyMessagesLimit uint = 300

type HttpStream struct {
	readCloser io.ReadCloser
	reader     *bufio.Reader
	isFinished bool
}

func (stream *HttpStream) Receive() (string, error) {
	if stream.isFinished {
		return "", errors.New("stream is close")
	}
	if stream.reader == nil {
		return "", errors.New("reader is nil")
	}
	var emptyMsgCount uint = 0
	for {
		line, err := stream.reader.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			emptyMsgCount++
			if emptyMsgCount > defaultEmptyMessagesLimit {
				return "", errors.New("stream has sent too many empty messages")
			}
			continue
		}
		lineStr := string(line)
		if lineStr == "[DONE]" {
			stream.Close()
			return "", io.EOF
		}
		return lineStr, nil
	}
}
func (stream *HttpStream) Close() {
	stream.isFinished = true
	if stream.readCloser != nil {
		stream.readCloser.Close()
	}
}

func CreateStream(method, url string, headers map[string]string, body string) (*HttpStream, int, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Connection", "keep-alive")
	resp, err := client.Do(request)
	if err != nil || resp == nil || resp.Body == nil {
		return nil, 0, err
	}
	reader := bufio.NewReader(resp.Body)
	return &HttpStream{
		readCloser: resp.Body,
		reader:     reader,
	}, resp.StatusCode, nil
}
