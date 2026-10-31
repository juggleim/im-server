package services

import (
	"net/http/httptest"
	"testing"

	"im-server/services/commonservices"

	"github.com/gin-gonic/gin"
)

func TestToRpcCtxIncludesDirectClientIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginCtx.Request = httptest.NewRequest("POST", "/apigateway/messages/private/send", nil)
	ginCtx.Request.RemoteAddr = "203.0.113.9:1234"
	ginCtx.Request.Header.Set("X-Forwarded-For", "198.51.100.8")

	rpcCtx := ToRpcCtx(ginCtx, "sender")
	if got := commonservices.GetClientIPFromCtx(rpcCtx); got != "203.0.113.9" {
		t.Fatalf("client IP=%q, want direct peer IP", got)
	}
}
