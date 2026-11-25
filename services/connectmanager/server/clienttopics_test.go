package server

import (
	"io"
	"testing"

	"im-server/commons/errs"
	"im-server/commons/logs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/connectmanager/server/codec"
	"im-server/services/connectmanager/server/imcontext"
)

type topicTestContext struct {
	imcontext.WsHandleContext
	attachment imcontext.Attachment
	writes     []interface{}
}

func (c *topicTestContext) Attachment() imcontext.Attachment     { return c.attachment }
func (c *topicTestContext) SetAttachment(a imcontext.Attachment) { c.attachment = a }
func (c *topicTestContext) Write(msg interface{})                { c.writes = append(c.writes, msg) }

func TestListenersRejectNonClientTopicsBeforeRouting(t *testing.T) {
	logs.InitLogs()
	logger := logs.GetInfoLogger()
	output := logger.Out
	logger.SetOutput(io.Discard)
	t.Cleanup(func() { logger.SetOutput(output) })

	type listener interface {
		PublishArrived(*codec.PublishMsgBody, int, imcontext.WsHandleContext)
		QueryArrived(*codec.QueryMsgBody, imcontext.WsHandleContext)
	}
	const index = 7
	for name, l := range map[string]listener{"user": &ImListenerImpl{}, "bot": &ImBotListenerImpl{}} {
		t.Run(name, func(t *testing.T) {
			for _, topic := range []string{"g_del_members", "unknown_client_topic"} {
				t.Run(topic, func(t *testing.T) {
					ctx := &topicTestContext{}
					// The WebSocket admin type must not bypass the client policy.
					imcontext.SetContextAttr(ctx, imcontext.StateKey_UserType, pbobjs.UserType_Admin)
					// No actor router is initialized: these requests must stop at the listener.
					l.PublishArrived(&codec.PublishMsgBody{Index: index, Topic: topic, TargetId: "test-target"}, 0, ctx)
					l.QueryArrived(&codec.QueryMsgBody{Index: index, Topic: topic, TargetId: "test-target"}, ctx)
					if len(ctx.writes) != 2 {
						t.Fatalf("got %d replies, want 2", len(ctx.writes))
					}
					pub := ctx.writes[0].(codec.IMessage).ToImWebsocketMsg().GetPubAckMsgBody()
					qry := ctx.writes[1].(codec.IMessage).ToImWebsocketMsg().GetQryAckMsgBody()
					if pub == nil || pub.Index != index || pub.Code != int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC) {
						t.Fatalf("unexpected publish ack: %v", pub)
					}
					if qry == nil || qry.Index != index || qry.Code != int32(errs.IMErrorCode_CONNECT_UNSUPPORTEDTOPIC) {
						t.Fatalf("unexpected query ack: %v", qry)
					}
				})
			}
		})
	}
}

func TestListenerClientTopicsAllowed(t *testing.T) {
	ctx := &topicTestContext{}
	if !allowClientPublish(ctx, &codec.PublishMsgBody{Topic: "g_msg"}) ||
		!allowClientQuery(ctx, &codec.QueryMsgBody{Topic: "sync_msgs"}) ||
		!allowClientQuery(ctx, &codec.QueryMsgBody{Topic: "unsub_user_status"}) || len(ctx.writes) != 0 {
		t.Fatal("supported client methods must pass without an error ack")
	}
}
