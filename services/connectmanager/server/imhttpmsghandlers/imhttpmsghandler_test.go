package imhttpmsghandlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/codec"
)

func TestHandlePublishClientTopic(t *testing.T) {
	const (
		clientTopic = "p_msg"
		serverTopic = "g_del_members"
		appKey      = "test-app"
		token       = "test-token"
		userID      = "test-user"
		targetID    = "test-target"
		index       = 42
		messageID   = "test-message"
	)
	tests := []struct {
		name, topic, target string
		tokenCode, wantCode errs.IMErrorCode
		wantRoute           bool
	}{
		{name: "unknown topic", topic: "test.unknown", target: targetID, wantCode: errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC},
		{name: "server-only topic", topic: serverTopic, target: targetID, wantCode: errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC},
		{name: "query-only topic", topic: "sync_msgs", target: targetID, wantCode: errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC},
		{name: "client topic", topic: clientTopic, target: targetID, wantRoute: true},
		{name: "invalid token", topic: serverTopic, target: targetID, tokenCode: errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL, wantCode: errs.IMErrorCode_CONNECT_TOKEN_ILLEGAL},
		{name: "missing topic", target: targetID, wantCode: errs.IMErrorCode_CONNECT_PARAM_REQUIRED},
		{name: "missing target", topic: clientTopic, wantCode: errs.IMErrorCode_CONNECT_PARAM_REQUIRED},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := tools.PbMarshal(codec.NewUserPublishMessage(&codec.PublishMsgBody{
				Index: index, Topic: tt.topic, TargetId: tt.target,
			}).ToImWebsocketMsg())
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Set("x-appkey", appKey)
			req.Header.Set("x-token", token)
			w := httptest.NewRecorder()

			var tokenCalls, routeCalls int
			handlePublish(w, req, func(app, tok string) (string, errs.IMErrorCode) {
				tokenCalls++
				if app != appKey || tok != token {
					t.Fatal("unexpected app key or token passed to validator")
				}
				return userID, tt.tokenCode
			}, func(msg bases.IRoute, _ time.Duration) (*pbobjs.RpcMessageWraper, error) {
				routeCalls++
				if !tt.wantRoute || tokenCalls != 1 || routeCalls != 1 {
					t.Fatalf("unexpected route call: token checks=%d", tokenCalls)
				}
				rpc, ok := msg.(*pbobjs.RpcMessageWraper)
				if !ok {
					t.Fatalf("route message type = %T", msg)
				}
				if rpc.Method != "upstream" || rpc.ExtParams[commonservices.RpcExtKey_RealMethod] != clientTopic ||
					rpc.AppKey != appKey || rpc.TargetId != userID || rpc.RequesterId != tt.target || rpc.ReqIndex != index {
					t.Errorf("unexpected RPC envelope: %v", rpc)
				}
				return &pbobjs.RpcMessageWraper{ResultCode: int32(errs.IMErrorCode_SUCCESS), MsgId: messageID}, nil
			})

			if tokenCalls != 1 || (routeCalls == 1) != tt.wantRoute {
				t.Fatalf("unexpected call counts: token=%d, route=%d", tokenCalls, routeCalls)
			}
			if w.Code != http.StatusOK {
				t.Fatalf("HTTP status = %d, want %d", w.Code, http.StatusOK)
			}
			var response codec.ImWebsocketMsg
			if err := tools.PbUnMarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			ack := response.GetPubAckMsgBody()
			if ack == nil || ack.Index != index || ack.Code != int32(tt.wantCode) {
				t.Fatalf("ack = %v, want index %d and code %v", ack, index, tt.wantCode)
			}
			if tt.wantRoute && ack.MsgId != messageID {
				t.Errorf("ack message ID = %q, want %q", ack.MsgId, messageID)
			}
		})
	}
}
