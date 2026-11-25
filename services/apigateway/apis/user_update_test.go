package apis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func TestUpdateUserForwardsServerAPIRequest(t *testing.T) {
	const (
		method = "upd_user_info"
		path   = "/apigateway/users/update"
		appKey = "test-app"
		userID = "test-user"
	)
	t.Setenv("IM_ACTOR_CALLBACK_WORKERS", "1")
	t.Setenv("IM_ACTOR_EXECUTOR_WORKERS", "1")
	if err := bases.InitImServer(nil); err != nil {
		t.Fatal(err)
	}
	received := make(chan *pbobjs.RpcMessageWraper, 1)
	bases.GetCluster().RegisterActor(method, func() actorsystem.IUntypedActor {
		return &updateUserActorStub{received: received}
	})

	engine := gin.New()
	engine.POST(path, func(ctx *gin.Context) {
		// Model the context produced by successful Server API authentication.
		ctx.Set(string(bases.CtxKey_AppKey), appKey)
	}, UpdateUser)
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(`{"user_id":"test-user","nickname":"Updated name"}`))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	engine.ServeHTTP(resp, req)
	var result struct{ Code *errs.IMErrorCode }
	if resp.Code != http.StatusOK || json.Unmarshal(resp.Body.Bytes(), &result) != nil ||
		result.Code == nil || *result.Code != errs.IMErrorCode_SUCCESS {
		t.Fatalf("unexpected API response: HTTP %d %s", resp.Code, resp.Body.String())
	}

	select {
	case rpc := <-received:
		if !rpc.GetIsFromApi() || rpc.GetRpcMsgType() != pbobjs.RpcMsgType_ServerPub ||
			rpc.GetAppKey() != appKey || rpc.GetMethod() != method || rpc.GetTargetId() != userID {
			t.Fatalf("unexpected RPC envelope: %v", rpc)
		}
		var info pbobjs.UserInfo
		if err := proto.Unmarshal(rpc.GetAppDataBytes(), &info); err != nil {
			t.Fatal(err)
		}
		if info.GetUserId() != userID || info.GetNickname() != "Updated name" {
			t.Fatalf("unexpected profile update: %v", &info)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("profile update actor did not receive the Server API request")
	}
}

type updateUserActorStub struct {
	actorsystem.UntypedActor
	received chan<- *pbobjs.RpcMessageWraper
}

func (*updateUserActorStub) CreateInputObj() proto.Message { return &pbobjs.RpcMessageWraper{} }

func (a *updateUserActorStub) OnReceive(_ context.Context, input proto.Message) {
	a.received <- input.(*pbobjs.RpcMessageWraper)
}
