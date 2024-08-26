package services

import (
	"encoding/json"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/commonservices/dbs"
)

type EventSubConfigReq struct {
	AppKey         string                            `json:"app_key"`
	EventSubConfig *commonservices.EventSubConfigObj `json:"event_sub_config"`
	EventSubSwitch *commonservices.EventSubSwitchObj `json:"event_sub_switch"`
}
type EventSubConfigResp struct {
	AppKey         string                            `json:"app_key"`
	EventSubConfig *commonservices.EventSubConfigObj `json:"event_sub_config"`
	EventSubSwitch []*EventSubSwitchModel            `json:"event_sub_switch"`
}
type EventSubSwitchModel struct {
	Name  string        `json:"name"`
	Items []*ConfigItem `json:"items"`
}

func SetEventSubConfig(req *EventSubConfigReq) AdminErrorCode {
	dao := dbs.AppExtDao{}
	dao.CreateOrUpdate(req.AppKey, "event_sub_config", tools.ToJson(req.EventSubConfig))
	dao.CreateOrUpdate(req.AppKey, "event_sub_switch", tools.ToJson(req.EventSubSwitch))
	return AdminErrorCode_Success
}

func GetEventSubConfig(appkey string) (AdminErrorCode, *EventSubConfigResp) {
	dao := dbs.AppExtDao{}
	ret := &EventSubConfigResp{
		AppKey:         appkey,
		EventSubConfig: &commonservices.EventSubConfigObj{},
		EventSubSwitch: []*EventSubSwitchModel{},
	}
	appExt, err := dao.Find(appkey, "event_sub_config")
	if err == nil && appExt.AppItemValue != "" {
		json.Unmarshal([]byte(appExt.AppItemValue), ret.EventSubConfig)
	}
	appExt, err = dao.Find(appkey, "event_sub_switch")
	switchMap := map[string]int{}
	if err == nil && appExt.AppItemValue != "" {
		json.Unmarshal([]byte(appExt.AppItemValue), &switchMap)
	}
	//消息订阅
	subSwitchModel := &EventSubSwitchModel{
		Name:  "消息订阅",
		Items: []*ConfigItem{},
	}
	subSwitchModel.Items = append(subSwitchModel.Items, &ConfigItem{
		Key:   "private_msg_sub_switch",
		Value: getSwitchValue("private_msg_sub_switch", switchMap),
		Name:  "单聊消息订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &ConfigItem{
		Key:   "group_msg_sub_switch",
		Value: getSwitchValue("group_msg_sub_switch", switchMap),
		Name:  "群聊消息订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &ConfigItem{
		Key:   "chatroom_msg_sub_switch",
		Value: getSwitchValue("chatroom_msg_sub_switch", switchMap),
		Name:  "聊天室消息订阅",
	})
	ret.EventSubSwitch = append(ret.EventSubSwitch, subSwitchModel)

	//在线状态订阅
	subSwitchModel = &EventSubSwitchModel{
		Name:  "在线状态",
		Items: []*ConfigItem{},
	}
	subSwitchModel.Items = append(subSwitchModel.Items, &ConfigItem{
		Key:   "online_sub_switch",
		Value: getSwitchValue("online_sub_switch", switchMap),
		Name:  "上线状态订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &ConfigItem{
		Key:   "offline_sub_switch",
		Value: getSwitchValue("offline_sub_switch", switchMap),
		Name:  "离线状态订阅",
	})
	ret.EventSubSwitch = append(ret.EventSubSwitch, subSwitchModel)

	return AdminErrorCode_Success, ret
}
func getSwitchValue(key string, switchMap map[string]int) int {
	if val, ok := switchMap[key]; ok {
		return val
	} else {
		return 0
	}
}
