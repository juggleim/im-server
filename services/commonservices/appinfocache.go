package commonservices

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"time"

	"im-server/commons/caches"
	"im-server/commons/tools"
	"im-server/services/commonservices/dbs"
	"im-server/services/commonservices/sms"
	"im-server/services/commonservices/transengines"

	"github.com/kataras/i18n"
)

var appInfoCache *caches.LruCache
var appLocks *tools.SegmentatedLocks

type AppInfo struct {
	AppKey       string    `default:"-"`
	AppSecret    string    `default:"-"`
	AppSecureKey string    `default:"-"`
	AppStatus    int       `default:"-"`
	CreatedTime  time.Time `default:"-"`

	EventSubConfigObj  *EventSubConfigObj `default:"-"`
	EventSubSwitchObj  *EventSubSwitchObj `default:"-"`
	SecurityDomainsObj *SecurityDomains   `default:"-"`
	ZegoConfigObj      *ZegoConfigObj     `default:"-"`

	TokenEffectiveMinute  int    `default:"0"`
	OfflineMsgSaveTime    int    `default:"1440"`
	OfflineCmdMsgSaveTime int    `default:"10080"`
	CloseOfflineMsg       bool   `default:"false"`
	IsSubApiMsg           bool   `default:"false"`
	IsOpenPush            bool   `default:"true"`
	PushLanguage          string `default:"en_US"`
	KickMode              int    `default:"0"`
	OpenVisualLog         bool   `default:"false"`
	RecordGlobalConvers   bool   `default:"false"`
	OpenSensitive         bool   `default:"false"`
	RecordMsgLogs         bool   `default:"false"`

	//group
	IsHideMsgBeforeJoinGroup bool `default:"false"`
	HideGrpMsg               bool `default:"false"`
	NotCheckGrpMember        bool `default:"false"`
	MaxGrpMemberCount        int  `default:"10000"`
	GrpMsgThreshold          int  `default:"10000"`
	MsgThreshold             int  `default:"3000"`
	OpenGrpSnapshot          bool `default:"false"`
	BigGrpThreshold          int  `default:"1000"`
	ClosePushGrpThreshold    int  `default:"0"`

	EventSubConfig  string `default:""`
	EventSubSwitch  string `default:""`
	SecurityDomains string `default:""`
	ZegoConfig      string `default:""`
	TransEngineConf string `default:""`
	SmsEngineConf   string `default:""`

	// TestItem  string
	// TestInt   int
	// TestBool  bool  `default:"true"`
	// TestInt64 int64 `default:"10"`

	//other configure
	MsgTransConfs *MsgTransConfs            `default:"-"`
	TransEngine   transengines.ITransEngine `default:"-"`
	SmsEngine     sms.ISmsEngine            `default:"-"`
	I18nKeys      *i18n.I18n                `default:"-"`
}

var notExistAppInfo *AppInfo

func init() {
	appLocks = tools.NewSegmentatedLocks(64)
	notExistAppInfo = &AppInfo{}

	appInfoCache = caches.NewLruCache("appinfo_cache", 10000, nil)
	appInfoCache.AddTimeoutAfterRead(5 * time.Minute)
	appInfoCache.AddTimeoutAfterCreate(10 * time.Minute)
	appInfoCache.SetValueCreator(func(key interface{}) interface{} {
		appTable := dbs.AppInfoDao{}
		app := appTable.FindByAppkey(key.(string))
		if app != nil {
			appInfo := &AppInfo{
				AppKey:       app.AppKey,
				AppSecret:    app.AppSecret,
				AppSecureKey: app.AppSecureKey,
				AppStatus:    app.AppStatus,
				CreatedTime:  app.CreatedTime,
			}

			appExtTable := dbs.AppExtDao{}
			appExtList := appExtTable.FindListByAppkey(key.(string))
			extMap := make(map[string]string)
			if len(appExtList) > 0 {
				for _, appExt := range appExtList {
					extMap[strings.ToLower(appExt.AppItemKey)] = appExt.AppItemValue
				}
			}
			FillObjField(appInfo, extMap)

			//event subscription config
			if appInfo.EventSubConfigObj == nil && appInfo.EventSubConfig != "" {
				eventSubConfig := &EventSubConfigObj{}
				err := json.Unmarshal([]byte(appInfo.EventSubConfig), eventSubConfig)
				if err == nil {
					appInfo.EventSubConfigObj = eventSubConfig
				}
			}
			if appInfo.EventSubSwitchObj == nil && appInfo.EventSubSwitch != "" {
				eventSubSwitch := &EventSubSwitchObj{}
				err := json.Unmarshal([]byte(appInfo.EventSubSwitch), eventSubSwitch)
				if err == nil {
					appInfo.EventSubSwitchObj = eventSubSwitch
				}
			}
			if appInfo.SecurityDomainsObj == nil && appInfo.SecurityDomains != "" {
				domains := &SecurityDomains{}
				err := json.Unmarshal([]byte(appInfo.SecurityDomains), domains)
				if err == nil {
					appInfo.SecurityDomainsObj = domains
				}
			}
			if appInfo.ZegoConfigObj == nil && appInfo.ZegoConfig != "" {
				zegoConfig := &ZegoConfigObj{}
				err := json.Unmarshal([]byte(appInfo.ZegoConfig), zegoConfig)
				if err == nil {
					appInfo.ZegoConfigObj = zegoConfig
				}
			}
			return appInfo
		}
		return notExistAppInfo
	})
}

func FillObjField(obj interface{}, valMap map[string]string) {
	FillObjFieldWithIgnore(obj, valMap, false)
}

func FillObjFieldWithIgnore(obj interface{}, valMap map[string]string, ignoreDefault bool) {
	objVal := reflect.ValueOf(obj).Elem()
	for i := 0; i < objVal.NumField(); i++ {
		fieldName := objVal.Type().Field(i).Name
		if !strings.HasPrefix(fieldName, "HasField_") {
			fieldType := objVal.Type().Field(i).Type
			fieldTag := objVal.Type().Field(i).Tag
			defaultStr := strings.TrimSpace(fieldTag.Get("default"))
			if !ignoreDefault && defaultStr != "" && defaultStr != "-" {
				setFieldValue(objVal.FieldByName(fieldName), fieldType, defaultStr)
			}
			lowerFieldName := tools.CamelToSnake(fieldName) //strings.ToLower(fieldName)
			if mapVal, ok := valMap[lowerFieldName]; ok {
				afterTrimMapVal := strings.TrimSpace(mapVal)
				setFieldValue(objVal.FieldByName(fieldName), fieldType, afterTrimMapVal)

				//Handle HasField_
				hasField := objVal.FieldByName("HasField_" + fieldName)
				if hasField.IsValid() {
					hasField.SetBool(true)
				}
			}
		}
	}
}

func setFieldValue(field reflect.Value, typ reflect.Type, val string) {
	typeStr := typ.String()
	if typeStr == "string" {
		field.Set(reflect.ValueOf(val))
	} else {
		if val != "" {
			if typeStr == "int" {
				intVal, err := strconv.Atoi(val)
				if err == nil {
					field.Set(reflect.ValueOf(intVal))
				}
			} else if typeStr == "int64" {
				int64Val, err := strconv.ParseInt(val, 10, 64)
				if err == nil {
					field.Set(reflect.ValueOf(int64Val))
				}
			} else if typeStr == "bool" {
				boolVal, err := strconv.ParseBool(val)
				if err == nil {
					field.Set(reflect.ValueOf(boolVal))
				}
			}
		}
	}
}

func GetAppInfo(appkey string) (*AppInfo, bool) {
	val, ok := appInfoCache.GetByCreator(appkey, nil)
	if ok {
		info := val.(*AppInfo)
		if info == notExistAppInfo {
			return nil, false
		} else {
			if info.AppStatus != 0 {
				return info, false
			}
			return info, true
		}
	} else {
		return nil, false
	}
}

type EventSubConfigObj struct {
	EventSubUrl  string `json:"event_sub_url"`
	EventSubAuth string `json:"event_sub_auth"`
}

type EventSubSwitchObj struct {
	PrivateMsgSubSwitch  int `json:"private_msg_sub_switch"`
	GroupMsgSubSwitch    int `json:"group_msg_sub_switch"`
	ChatroomMsgSubSwitch int `json:"chatroom_msg_sub_switch"`
	OnlineSubSwitch      int `json:"online_sub_switch"`
	OfflineSubSwitch     int `json:"offline_sub_switch"`
}

type SecurityDomains struct {
	Domains []string `json:"domains"`
}

type ZegoConfigObj struct {
	AppId  int64  `json:"app_id"`
	Secret string `json:"secret"`
}
