package tests

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/apigateway/models"
	"im-server/services/commonservices/msgdefines"
	"im-server/simulator/utils"
	"im-server/simulator/wsclients"

	serversdk "github.com/juggleim/imserver-sdk-go"
	"github.com/stretchr/testify/require"
)

func TestGlobalConverTagEndToEnd(t *testing.T) {
	uniqueSuffix := strconv.FormatInt(time.Now().UnixNano(), 36)
	groupID := "sg" + uniqueSuffix
	globalTag := "sgt_" + uniqueSuffix
	nonexistentTag := "smt_" + uniqueSuffix
	tags := []string{globalTag}

	apiURL := os.Getenv("SIMULATOR_API_URL")
	if apiURL == "" {
		apiURL = ApiURL
	}
	appSecret := os.Getenv("SIMULATOR_APP_SECRET")
	if appSecret == "" {
		appSecret = AppSecret
	}
	wsAddr := os.Getenv("SIMULATOR_WS_ADDR")
	if wsAddr == "" {
		wsAddr = WsAddr
	}
	sdk := serversdk.NewJuggleIMSdk(Appkey, appSecret, apiURL)
	code, _, err := sdk.HttpCall(http.MethodPost, "/apigateway/groups/add", &models.GroupMembersReq{
		GroupId:          groupID,
		GroupName:        "simulator global conversation tag",
		MemberIds:        []string{User1},
		GlobalConverTags: &tags,
	}, nil)
	require.NoError(t, err)
	require.Equal(t, serversdk.ApiCode_Success, code)
	defer func() {
		_, _, _ = sdk.DissolveGroup(groupID)
	}()

	code, _, err = sdk.HttpCall(http.MethodPost, "/apigateway/groups/members/add", &models.GroupMembersReq{
		GroupId:   groupID,
		MemberIds: []string{User2},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, serversdk.ApiCode_Success, code)

	groupInfo := &models.GroupInfo{}
	code, _, err = sdk.HttpCall(http.MethodGet, "/apigateway/groups/info?group_id="+groupID, nil, groupInfo)
	require.NoError(t, err)
	require.Equal(t, serversdk.ApiCode_Success, code)
	require.NotNil(t, groupInfo.GlobalConverTags)
	require.ElementsMatch(t, tags, *groupInfo.GlobalConverTags)

	sender := wsclients.NewWsImClient(wsAddr, Appkey, Token1, nil, nil, nil)
	connectCode, connectAck := sender.Connect("", "")
	require.Equal(t, utils.ClientErrorCode_Success, connectCode)
	require.NotNil(t, connectAck)
	require.Equal(t, User1, connectAck.UserId)
	defer sender.Disconnect()

	receiver := wsclients.NewWsImClient(wsAddr, Appkey, Token2, nil, nil, nil)
	connectCode, connectAck = receiver.Connect("", "")
	require.Equal(t, utils.ClientErrorCode_Success, connectCode)
	require.NotNil(t, connectAck)
	require.Equal(t, User2, connectAck.UserId)
	defer receiver.Disconnect()

	flags := msgdefines.SetCountMsg(msgdefines.SetStoreMsg(0))
	sendCode, sendAck := sender.SendGroupMsg(groupID, &pbobjs.UpMsg{
		MsgType:    "sim_g_tag",
		MsgContent: []byte(fmt.Sprintf(`{"tag":%q}`, globalTag)),
		Flags:      flags,
	})
	require.Equal(t, utils.ClientErrorCode_Success, sendCode)
	require.NotNil(t, sendAck)
	require.NotEmpty(t, sendAck.MsgId)

	var receiverConver *pbobjs.Conversation
	require.Eventually(t, func() bool {
		qryCode, conver := receiver.QryConversation(&pbobjs.QryConverReq{
			TargetId:    groupID,
			ChannelType: pbobjs.ChannelType_Group,
		})
		if qryCode != utils.ClientErrorCode_Success || conver == nil {
			return false
		}
		receiverConver = conver
		return converHasGlobalTag(conver, globalTag) && conver.UnreadCount > 0
	}, 15*time.Second, 200*time.Millisecond, "receiver conversation did not contain global tag %q", globalTag)
	require.NotNil(t, receiverConver)

	require.Eventually(t, func() bool {
		qryCode, resp := receiver.QryConversations(&pbobjs.QryConversationsReq{
			Count:       1000,
			ChannelType: pbobjs.ChannelType_Group,
			Tag:         globalTag,
		})
		return qryCode == utils.ClientErrorCode_Success && responseHasConversation(resp, groupID)
	}, 15*time.Second, 200*time.Millisecond, "qry_convers did not return group %q for global tag %q", groupID, globalTag)

	qryCode, exceptResp := receiver.QryConversations(&pbobjs.QryConversationsReq{
		Count:       1000,
		ChannelType: pbobjs.ChannelType_Group,
		ExceptTag:   globalTag,
	})
	require.Equal(t, utils.ClientErrorCode_Success, qryCode)
	require.False(t, responseHasConversation(exceptResp, groupID), "qry_convers exceptTag returned the tagged group")

	groupConver := &pbobjs.SimpleConversation{
		TargetId:    groupID,
		ChannelType: pbobjs.ChannelType_Group,
	}
	matchingUnread := qryFilteredUnreadCount(t, receiver, &pbobjs.ConverFilter{
		IncludeConvers: []*pbobjs.SimpleConversation{groupConver},
		Tag:            globalTag,
	})
	require.Equal(t, receiverConver.UnreadCount, matchingUnread)
	require.Greater(t, matchingUnread, int64(0))

	excludedUnread := qryFilteredUnreadCount(t, receiver, &pbobjs.ConverFilter{
		IncludeConvers: []*pbobjs.SimpleConversation{groupConver},
		ExceptTag:      globalTag,
	})
	require.Zero(t, excludedUnread)

	nonmatchingUnread := qryFilteredUnreadCount(t, receiver, &pbobjs.ConverFilter{
		IncludeConvers: []*pbobjs.SimpleConversation{groupConver},
		Tag:            nonexistentTag,
	})
	require.Zero(t, nonmatchingUnread)
}

func converHasGlobalTag(conver *pbobjs.Conversation, tag string) bool {
	if conver == nil {
		return false
	}
	for _, converTag := range conver.ConverTags {
		if converTag.Tag == tag && converTag.TagType == pbobjs.ConverTagType_GlobalConverTag {
			return true
		}
	}
	return false
}

func responseHasConversation(resp *pbobjs.QryConversationsResp, targetID string) bool {
	if resp == nil {
		return false
	}
	for _, conver := range resp.Conversations {
		if conver.TargetId == targetID && conver.ChannelType == pbobjs.ChannelType_Group {
			return true
		}
	}
	return false
}

func qryFilteredUnreadCount(t *testing.T, client *wsclients.WsImClient, filter *pbobjs.ConverFilter) int64 {
	t.Helper()
	code, resp := client.QryTotalUnreadCount(&pbobjs.QryTotalUnreadCountReq{Filter: filter})
	require.Equal(t, utils.ClientErrorCode_Success, code)
	require.NotNil(t, resp)
	return resp.TotalCount
}
