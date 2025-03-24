package commonservices

import (
	"context"
	"fmt"
	"im-server/commons/tools"
	"im-server/services/commonservices/logs"
	"reflect"
	"strings"
	"time"
)

type AttItemType int
type AttItemKey string
type BotConnectType int

const (
	AttItemType_Att     AttItemType = 0
	AttItemType_Setting AttItemType = 1
	AttItemType_Status  AttItemType = 2

	BotConnectType_Webhook   BotConnectType = 0
	BotConnectType_Websocket BotConnectType = 1

	//setting keys of group or group_members
	AttItemKey_HideGrpMsg        AttItemKey = "hide_grp_msg"
	AttItemKey_GrpCreator        AttItemKey = "grp_creator"
	AttItemKey_GrpAnnouncement   AttItemKey = "grp_announcement"
	AttItemKey_GrpVerifyType     AttItemKey = "grp_verify_type"
	AttItemKey_GrpAdministrators AttItemKey = "grp_administrators"
	AttItemKey_GrpDisplayName    AttItemKey = "grp_display_name"

	//setting keys of users
	AttItemKey_Language      AttItemKey = "language"
	AttItemKey_Undisturb     AttItemKey = "undisturb"
	AttItemKey_PriGlobalMute AttItemKey = "pri_global_mute"

	//setting keys of bots
	AttItemKey_Bot_Type    AttItemKey = "bot_type"
	AttItemKey_Bot_WebHook AttItemKey = "bot_webhook"
	AttItemKey_Bot_ApiKey  AttItemKey = "bot_api_key"
	AttItemKey_Bot_BotConf AttItemKey = "bot_conf"
)

type GroupSettings struct {
	HideGrpMsg          bool `default:"false"`
	HasField_HideGrpMsg bool
}

type GrpMemberSettings struct {
	HideGrpMsg          bool `default:"false"`
	HasField_HideGrpMsg bool
}

type UserSettings struct {
	Language     string `default:""`
	Undisturb    string `default:""`
	UndisturbObj *UserUndisturb
}

var GrpMemberSettingKeys map[AttItemKey]bool
var GroupSettingKeys map[AttItemKey]bool
var UserSettingKeys map[AttItemKey]bool

func init() {
	GroupSettingKeys = make(map[AttItemKey]bool)
	GroupSettingKeys[AttItemKey_HideGrpMsg] = true

	GrpMemberSettingKeys = make(map[AttItemKey]bool)
	GrpMemberSettingKeys[AttItemKey_HideGrpMsg] = true

	UserSettingKeys = make(map[AttItemKey]bool)
	UserSettingKeys[AttItemKey_Language] = true
	UserSettingKeys[AttItemKey_Undisturb] = true
}

func CheckGroupSettingKey(key string) bool {
	_, exist := GroupSettingKeys[AttItemKey(key)]
	return exist
}

func CheckGrpMemberSettingKey(key string) bool {
	_, exist := GrpMemberSettingKeys[AttItemKey(key)]
	return exist
}

func CheckUserSettingKey(key string) bool {
	_, exist := UserSettingKeys[AttItemKey(key)]
	return exist
}

func Obj2Map(obj interface{}) map[string]string {
	valMap := make(map[string]string)
	objVal := reflect.ValueOf(obj).Elem()
	for i := 0; i < objVal.NumField(); i++ {
		fieldValue := objVal.Field(i)
		fieldName := objVal.Type().Field(i).Name
		if !strings.HasPrefix(fieldName, "HasField_") {
			lowerFieldName := tools.CamelToSnake(fieldName)
			valMap[lowerFieldName] = fmt.Sprintf("%v", fieldValue)
		}
	}
	return valMap
}

type UserUndisturb struct {
	Switch   bool                 `json:"switch"`
	Timezone string               `json:"timezone"`
	Rules    []*UserUndisturbItem `json:"rules"`
}
type UserUndisturbItem struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func (user *UserUndisturb) CheckUndisturb(ctx context.Context, userId string) bool {
	if user.Switch {
		if len(user.Rules) <= 0 {
			return true
		}
		now := time.Now()
		if user.Timezone != "" {
			location, err := time.LoadLocation(user.Timezone)
			if err == nil && location != nil {
				now = now.In(location)
			} else {
				logs.WithContext(ctx).Errorf("user_id:%s\ttimezone:%s\tlocation error:%s", userId, user.Timezone, err)
			}
		}
		cur := Hhmm2Int(now.Format("15:04"))
		for _, rule := range user.Rules {
			start := Hhmm2Int(rule.Start)
			end := Hhmm2Int(rule.End)
			if end >= start {
				if cur >= start && cur <= end {
					return true
				}
			} else {
				if cur >= start || cur <= end {
					return true
				}
			}
		}
	}
	return false
}
func Hhmm2Int(str string) int64 {
	if len(str) != 5 {
		return 0
	}
	val := str[0:2] + str[3:5]

	ret, err := tools.String2Int64(val)
	if err != nil {
		return 0
	}
	return ret
}
