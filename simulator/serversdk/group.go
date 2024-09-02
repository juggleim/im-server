package serversdk

import "net/http"

type GroupMembersReq struct {
	GroupId       string   `json:"group_id"`
	GroupName     string   `json:"group_name"`
	GroupPortrait string   `json:"group_portrait"`
	MemberIds     []string `json:"member_ids"`
}

func (sdk *JuggleIMSdk) CreateGroup(groupMembers GroupMembersReq) (ApiCode, string, error) {
	url := sdk.ApiUrl + "/apigateway/groups/add"
	code, traceId, err := sdk.HttpCall(http.MethodPost, url, groupMembers, nil)
	return code, traceId, err
}

func (sdk *JuggleIMSdk) GroupAddMembers(groupMembers GroupMembersReq) (ApiCode, string, error) {
	return sdk.CreateGroup(groupMembers)
}

type GroupInfo struct {
	GroupId       string            `json:"group_id"`
	GroupName     string            `json:"group_name"`
	GroupPortrait string            `json:"group_portrait"`
	IsMute        int               `json:"is_mute"`
	UpdatedTime   int64             `json:"updated_time"`
	ExtFields     map[string]string `json:"ext_fields"`
}

func (sdk *JuggleIMSdk) DissolveGroup(groupId string) (ApiCode, string, error) {
	url := sdk.ApiUrl + "/apigateway/groups/del"
	code, traceId, err := sdk.HttpCall(http.MethodPost, url, &GroupInfo{
		GroupId: groupId,
	}, nil)
	return code, traceId, err
}
