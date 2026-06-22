package services

import (
	"errors"
	"im-server/services/commonservices"
	"im-server/services/connectmanager/server/imcontext"
	"sync"
	"testing"
)

type fakeWsContext struct {
	attachment imcontext.Attachment
	active     bool
}

func newFakeWsContext(appkey, userid, session string) *fakeWsContext {
	ctx := &fakeWsContext{active: true}
	imcontext.SetContextAttr(ctx, imcontext.StateKey_Appkey, appkey)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_UserID, userid)
	imcontext.SetContextAttr(ctx, imcontext.StateKey_ConnectSession, session)
	return ctx
}

func (ctx *fakeWsContext) Write(message interface{}) {}

func (ctx *fakeWsContext) Close(err error) {
	ctx.active = false
}

func (ctx *fakeWsContext) Attachment() imcontext.Attachment {
	return ctx.attachment
}

func (ctx *fakeWsContext) SetAttachment(attachment imcontext.Attachment) {
	ctx.attachment = attachment
}

func (ctx *fakeWsContext) IsActive() bool {
	return ctx.active
}

func (ctx *fakeWsContext) RemoteAddr() string {
	return "127.0.0.1:10000"
}

func TestAtomicConnectCountsAddSessions(t *testing.T) {
	resetOnlineConnectStateForTest()

	ctx1 := newFakeWsContext("app1", "u1", "s1")
	ctx2 := newFakeWsContext("app1", "u1", "s2")
	identifier := getUserIdentifier("app1", "u1")
	storeOnlineSession("s1", ctx1)
	addUserSessionLocked(identifier, "app1", "s1", ctx1)

	metrics := GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 1, SessionConnectCount: 1}) {
		t.Fatalf("metrics after first add = %+v", metrics)
	}

	storeOnlineSession("s2", ctx2)
	addUserSessionLocked(identifier, "app1", "s2", ctx2)
	metrics = GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 2, SessionConnectCount: 2}) {
		t.Fatalf("metrics after second add = %+v", metrics)
	}
	if got := GetConnectCountByUser("app1", "u1"); got != 2 {
		t.Fatalf("GetConnectCountByUser = %d, want 2", got)
	}
}

func TestAtomicConnectCountsDuplicateAdd(t *testing.T) {
	resetOnlineConnectStateForTest()

	ctx1 := newFakeWsContext("app1", "u1", "s1")
	ctx2 := newFakeWsContext("app1", "u1", "s1")
	identifier := getUserIdentifier("app1", "u1")
	storeOnlineSession("s1", ctx1)
	addUserSessionLocked(identifier, "app1", "s1", ctx1)
	storeOnlineSession("s1", ctx2)
	addUserSessionLocked(identifier, "app1", "s1", ctx2)

	metrics := GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 1, SessionConnectCount: 1}) {
		t.Fatalf("metrics after duplicate add = %+v", metrics)
	}
	if obj, ok := OnlineSessionConnectMap.Load("s1"); !ok || obj != ctx2 {
		t.Fatal("duplicate add should refresh session context without incrementing counters")
	}
}

func TestAtomicConnectCountsRemoveSessions(t *testing.T) {
	resetOnlineConnectStateForTest()

	ctx1 := newFakeWsContext("app1", "u1", "s1")
	ctx2 := newFakeWsContext("app1", "u1", "s2")
	identifier := getUserIdentifier("app1", "u1")
	storeOnlineSession("s1", ctx1)
	userSessionMap := addUserSessionLocked(identifier, "app1", "s1", ctx1)
	storeOnlineSession("s2", ctx2)
	addUserSessionLocked(identifier, "app1", "s2", ctx2)

	deleteOnlineSession("s1")
	if !removeUserSessionLocked(identifier, "app1", "s1", userSessionMap) {
		t.Fatal("first remove should report a user-session removal")
	}
	metrics := GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 1, SessionConnectCount: 1}) {
		t.Fatalf("metrics after one remove = %+v", metrics)
	}

	deleteOnlineSession("s1")
	if removeUserSessionLocked(identifier, "app1", "s1", userSessionMap) {
		t.Fatal("repeated remove should not report a user-session removal")
	}
	metrics = GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 1, SessionConnectCount: 1}) {
		t.Fatalf("metrics after repeated remove = %+v", metrics)
	}

	deleteOnlineSession("s2")
	if !removeUserSessionLocked(identifier, "app1", "s2", userSessionMap) {
		t.Fatal("last remove should report a user-session removal")
	}
	metrics = GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{}) {
		t.Fatalf("metrics after all removes = %+v", metrics)
	}
}

func TestAtomicConnectCountsKickRemovalPath(t *testing.T) {
	resetOnlineConnectStateForTest()

	ctx1 := newFakeWsContext("app1", "u1", "s1")
	ctx2 := newFakeWsContext("app1", "u1", "s2")
	identifier := getUserIdentifier("app1", "u1")
	storeOnlineSession("s1", ctx1)
	userSessionMap := addUserSessionLocked(identifier, "app1", "s1", ctx1)
	storeOnlineSession("s2", ctx2)
	addUserSessionLocked(identifier, "app1", "s2", ctx2)

	if !removeUserSessionLocked(identifier, "app1", "s2", userSessionMap) {
		t.Fatal("kick remove should update user counters")
	}
	kicked, ok := loadAndDeleteOnlineSession("s2")
	if !ok || kicked != ctx2 {
		t.Fatal("kick remove should delete existing session")
	}
	ctx2.Close(errors.New("kick off"))

	metrics := GetOnlineConnectMetrics()
	if metrics != (commonservices.ClientConnectMetrics{OnlineUserCount: 1, UserConnectCount: 1, SessionConnectCount: 1}) {
		t.Fatalf("metrics after kick remove = %+v", metrics)
	}
}

func TestForeachAppConnectCount(t *testing.T) {
	resetOnlineConnectStateForTest()

	app1Ctx := newFakeWsContext("app1", "u1", "s1")
	app2Ctx := newFakeWsContext("app2", "u2", "s2")
	addUserSessionLocked(getUserIdentifier("app1", "u1"), "app1", "s1", app1Ctx)
	addUserSessionLocked(getUserIdentifier("app2", "u2"), "app2", "s2", app2Ctx)
	addUserSessionLocked(getUserIdentifier("app2", "u3"), "app2", "s3", app2Ctx)

	counts := map[string]int64{}
	foreachAppConnectCount(func(appkey string, count int64) {
		counts[appkey] = count
	})

	if counts["app1"] != 1 {
		t.Fatalf("app1 count = %d, want 1", counts["app1"])
	}
	if counts["app2"] != 2 {
		t.Fatalf("app2 count = %d, want 2", counts["app2"])
	}

	obj, _ := OnlineUserConnectMap.Load(getUserIdentifier("app1", "u1"))
	removeUserSessionLocked(getUserIdentifier("app1", "u1"), "app1", "s1", obj.(map[string]imcontext.WsHandleContext))
	counts = map[string]int64{}
	foreachAppConnectCount(func(appkey string, count int64) {
		counts[appkey] = count
	})
	if _, ok := counts["app1"]; ok {
		t.Fatal("zero app count should not be reported")
	}
}

func TestResetOnlineConnectStateForTest(t *testing.T) {
	resetOnlineConnectStateForTest()
	OnlineUserConnectMap.Store("u", map[string]imcontext.WsHandleContext{})
	OnlineSessionConnectMap.Store("s", newFakeWsContext("app1", "u", "s"))
	appConnectCountMap.Store("app1", &sync.Map{})

	resetOnlineConnectStateForTest()
	if metrics := GetOnlineConnectMetrics(); metrics != (commonservices.ClientConnectMetrics{}) {
		t.Fatalf("metrics after reset = %+v", metrics)
	}
	if _, ok := OnlineUserConnectMap.Load("u"); ok {
		t.Fatal("OnlineUserConnectMap not reset")
	}
	if _, ok := OnlineSessionConnectMap.Load("s"); ok {
		t.Fatal("OnlineSessionConnectMap not reset")
	}
}
