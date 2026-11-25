package apis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func TestGroupDelMembersForwardsServerAPIRequest(t *testing.T) {
	const (
		method  = "g_del_members"
		path    = "/apigateway/groups/members/del"
		appKey  = "test-app"
		session = "api_test_session"
		groupID = "test-group"
	)

	// The cluster is process-wide; keep this test serial and use minimal worker pools.
	t.Setenv("IM_ACTOR_CALLBACK_WORKERS", "1")
	t.Setenv("IM_ACTOR_EXECUTOR_WORKERS", "1")
	if err := bases.InitImServer(nil); err != nil {
		t.Fatalf("init local actor cluster: %v", err)
	}
	received := make(chan *pbobjs.RpcMessageWraper, 1)
	bases.GetCluster().RegisterActor(method, func() actorsystem.IUntypedActor {
		return &groupDelMembersActorStub{received: received}
	})

	engine := gin.New()
	engine.POST(path, func(ctx *gin.Context) {
		// Supply the context from successful Server API authentication without an app lookup.
		ctx.Set(string(bases.CtxKey_AppKey), appKey)
		ctx.Set(string(bases.CtxKey_Session), session)
	}, GroupDelMembers)
	request := httptest.NewRequest(http.MethodPost, path,
		strings.NewReader(`{"group_id":"test-group","member_ids":["test-member-1","test-member-2"]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("HTTP status = %d, want %d; body: %s", response.Code, http.StatusOK, response.Body.String())
	}
	var result struct {
		Code *int `json:"code"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode HTTP response: %v", err)
	}
	if result.Code == nil || *result.Code != int(errs.IMErrorCode_SUCCESS) {
		t.Fatalf("expected API success code; body: %s", response.Body.String())
	}

	var rpc *pbobjs.RpcMessageWraper
	select {
	case rpc = <-received:
	default:
		t.Fatal("group member removal actor did not receive an RPC")
	}
	if rpc.GetMethod() != method || rpc.GetTargetId() != groupID || rpc.GetAppKey() != appKey ||
		rpc.GetSession() != session || !rpc.GetIsFromApi() {
		t.Errorf("unexpected RPC envelope: method=%q target=%q app=%q session=%q fromAPI=%t",
			rpc.GetMethod(), rpc.GetTargetId(), rpc.GetAppKey(), rpc.GetSession(), rpc.GetIsFromApi())
	}
	var payload pbobjs.GroupMembersReq
	if err := proto.Unmarshal(rpc.GetAppDataBytes(), &payload); err != nil {
		t.Fatalf("decode actor request: %v", err)
	}
	want := &pbobjs.GroupMembersReq{GroupId: groupID, MemberIds: []string{"test-member-1", "test-member-2"}}
	if !proto.Equal(&payload, want) {
		t.Errorf("actor request = %v, want %v", &payload, want)
	}
}

type groupDelMembersActorStub struct {
	actorsystem.UntypedActor
	received chan<- *pbobjs.RpcMessageWraper
}

func (*groupDelMembersActorStub) CreateInputObj() proto.Message {
	return &pbobjs.RpcMessageWraper{}
}

func (a *groupDelMembersActorStub) OnReceive(_ context.Context, input proto.Message) {
	a.received <- input.(*pbobjs.RpcMessageWraper)
	// An empty GroupMembersResp encodes as zero bytes.
	a.GetSender().TellAndNoSender(&pbobjs.RpcMessageWraper{ResultCode: int32(errs.IMErrorCode_SUCCESS)})
}
