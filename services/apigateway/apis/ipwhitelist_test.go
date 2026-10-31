package apis

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestIPWhitelistMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name       string
		entries    []string
		remoteAddr string
		wantStatus int
	}{
		{name: "not configured", remoteAddr: "203.0.113.9:1234", wantStatus: http.StatusOK},
		{name: "exact IP", entries: []string{"203.0.113.9"}, remoteAddr: "203.0.113.9:1234", wantStatus: http.StatusOK},
		{name: "CIDR", entries: []string{"10.20.0.0/16"}, remoteAddr: "10.20.3.4:1234", wantStatus: http.StatusOK},
		{name: "IPv6", entries: []string{"2001:db8::/32"}, remoteAddr: "[2001:db8::1]:1234", wantStatus: http.StatusOK},
		{name: "forbidden", entries: []string{"203.0.113.9"}, remoteAddr: "203.0.113.10:1234", wantStatus: http.StatusForbidden},
		{name: "invalid remote address", entries: []string{"203.0.113.9"}, remoteAddr: "unknown", wantStatus: http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			middleware, err := NewIPWhitelistMiddleware(test.entries)
			if err != nil {
				t.Fatal(err)
			}
			engine := gin.New()
			engine.Use(middleware)
			engine.GET("/apigateway/test", func(ctx *gin.Context) {
				ctx.Status(http.StatusOK)
			})

			request := httptest.NewRequest(http.MethodGet, "/apigateway/test", nil)
			request.RemoteAddr = test.remoteAddr
			request.Header.Set("X-Forwarded-For", "203.0.113.9")
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)

			if response.Code != test.wantStatus {
				t.Fatalf("status=%d body=%s, want=%d", response.Code, response.Body.String(), test.wantStatus)
			}
			if response.Code == http.StatusForbidden && !strings.Contains(response.Body.String(), `"code":10012`) {
				t.Fatalf("unexpected forbidden response: %s", response.Body.String())
			}
		})
	}
}

func TestIPWhitelistRejectsInvalidConfiguration(t *testing.T) {
	for _, entries := range [][]string{{"not-an-ip"}, {""}, {"10.0.0.0/99"}} {
		if _, err := NewIPWhitelistMiddleware(entries); err == nil {
			t.Fatalf("expected configuration error for %q", entries)
		}
	}
}
