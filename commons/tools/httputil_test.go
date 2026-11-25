package tools

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHttpDoBytesReusesConnection(t *testing.T) {
	var newConnections int64
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.Copy(io.Discard, request.Body)
		_, _ = writer.Write([]byte("ok"))
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateNew {
			atomic.AddInt64(&newConnections, 1)
		}
	}
	server.Start()
	defer server.Close()

	const requestCount = 10
	for i := 0; i < requestCount; i++ {
		body, statusCode, err := HttpDoBytes(http.MethodPost, server.URL, nil, fmt.Sprintf("request-%d", i))
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		if statusCode != http.StatusOK {
			t.Fatalf("request %d returned status %d", i, statusCode)
		}
		if string(body) != "ok" {
			t.Fatalf("request %d returned body %q", i, body)
		}
	}

	if got := atomic.LoadInt64(&newConnections); got != 1 {
		t.Fatalf("expected one reused connection for %d requests, got %d connections", requestCount, got)
	}
}

func TestHttpDoBytesWithContextCancelsRequest(t *testing.T) {
	requestStarted := make(chan struct{})
	releaseHandler := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		close(requestStarted)
		<-releaseHandler
	}))
	defer func() {
		close(releaseHandler)
		server.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, _, err := HttpDoBytesWithContext(ctx, http.MethodPost, server.URL, nil, "request")
		result <- err
	}()

	select {
	case <-requestStarted:
	case <-time.After(time.Second):
		t.Fatal("request did not reach server")
	}
	cancel()

	select {
	case err := <-result:
		if err == nil {
			t.Fatal("expected canceled request to return an error")
		}
	case <-time.After(time.Second):
		t.Fatal("canceled request did not return promptly")
	}
}
