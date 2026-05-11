package tests

import (
	"encoding/json"
	"testing"
	"time"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/simulator/examples"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 与 examples.Token1 / Token2 对应的业务用户 id（需已在环境中注册；与 simulator/main 常用 demo 一致）
const (
	statusSubUser1 = "userid1"
	statusSubUser2 = "userid2"
)

const statusSubWsAddr = "ws://127.0.0.1:9002"

func TestStatusSubscriptionOfflineThenOnlineOfflineNotify(t *testing.T) {
	statusCh := make(chan *pbobjs.DownMsg, 32)

	onMsg := func(msg *pbobjs.DownMsg) {
		if msg == nil {
			return
		}
		if msg.ChannelType != pbobjs.ChannelType_SubStatus || msg.MsgType != "jg:onlinechg" {
			return
		}
		select {
		case statusCh <- msg:
		default:
		}
	}

	client1 := wsclients.NewWsImClient(statusSubWsAddr, "appkey", examples.Token1, onMsg, examples.OnStreamMsg, examples.OnDisconnect)
	client1.DeviceId = "status-sub-test-u1"
	code1, ack1 := client1.Connect("", "")
	if code1 != utils.ClientErrorCode_Success {
		t.Skipf("无法连接 %s（code=%v），请先启动 launcher", statusSubWsAddr, code1)
		return
	}
	require.Equal(t, int32(0), ack1.Code)
	require.Equal(t, statusSubUser1, ack1.UserId, "Token1 应对应 userid1，否则请更换 Token 或常量")

	// 1) 订阅 userid2，应答中应体现离线（若无终端在线）
	subCode, list := client1.SubUsers(&pbobjs.SubUsersReq{UserIds: []string{statusSubUser2}})
	require.Equal(t, utils.ClientErrorCode_Success, subCode, "sub_users 调用失败")
	require.NotNil(t, list)
	require.Len(t, list.Items, 1)
	require.Equal(t, statusSubUser2, list.Items[0].UserId)
	require.NotNil(t, list.Items[0].OnlineStatus)
	require.Equal(t, statusSubUser2, list.Items[0].OnlineStatus.UserId)
	assert.False(t, list.Items[0].OnlineStatus.IsOnline, "userid2 未登录时应为离线")

	// 同步关系落地后再连用户 2，降低偶发竞态
	time.Sleep(300 * time.Millisecond)

	client2 := wsclients.NewWsImClient(statusSubWsAddr, "appkey", examples.Token2, examples.OnMessage, examples.OnStreamMsg, examples.OnDisconnect)
	client2.DeviceId = "status-sub-test-u2"
	code2, ack2 := client2.Connect("", "")
	require.Equal(t, utils.ClientErrorCode_Success, code2)
	require.Equal(t, int32(0), ack2.Code)
	require.Equal(t, statusSubUser2, ack2.UserId, "Token2 应对应 userid2")

	msgOnline := recvOnlineChg(t, statusCh, 15*time.Second, true)
	require.NotNil(t, msgOnline)
	require.Equal(t, statusSubUser2, msgOnline.SenderId)
	require.Equal(t, pbobjs.ChannelType_SubStatus, msgOnline.ChannelType)

	client2.Disconnect()
	time.Sleep(200 * time.Millisecond)

	msgOffline := recvOnlineChg(t, statusCh, 15*time.Second, false)
	require.NotNil(t, msgOffline)
	require.Equal(t, statusSubUser2, msgOffline.SenderId)

	_ = client1.UnSubUsers(&pbobjs.SubUsersReq{UserIds: []string{statusSubUser2}})
	client1.Disconnect()
}

type onlineChgBody struct {
	IsOnline bool `json:"is_online"`
}

func recvOnlineChg(t *testing.T, ch <-chan *pbobjs.DownMsg, timeout time.Duration, wantOnline bool) *pbobjs.DownMsg {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case msg := <-ch:
			var body onlineChgBody
			if err := json.Unmarshal(msg.MsgContent, &body); err != nil {
				continue
			}
			if body.IsOnline == wantOnline {
				return msg
			}
		case <-timer.C:
			t.Fatalf("超时未收到 jg:onlinechg is_online=%v", wantOnline)
			return nil
		}
	}
}
