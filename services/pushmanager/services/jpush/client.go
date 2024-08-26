package jpush

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/levigross/grequests"
)

var (
	PushHost   = "https://api.jpush.cn/v3"
	DeviceHost = "https://device.jpush.cn/v3"
)

type JpushClient struct {
	*grequests.Session
	host         string
	appKey       string
	masterSecret string
}

func NewJpushClient(appKey, appSecret string) *JpushClient {
	return &JpushClient{
		host:         PushHost,
		appKey:       appKey,
		masterSecret: appSecret,
		Session: grequests.NewSession(&grequests.RequestOptions{
			UserAgent: "go-jpush/v0.1.2",
			Auth:      []string{appKey, appSecret},
			Headers: map[string]string{
				"Accept": "application/json",
			},
		}),
	}
}

func New(appKey, appSecret string) *JpushClient {
	return NewJpushClient(appKey, appSecret)
}

func (j *JpushClient) Url(path string) string {
	return j.host + path
}

// PushCid cid 是用于防止 api 调用端重试造成服务端的重复推送而定义的一个推送参数。
// 用户使用一个 cid 推送后，再次使用相同的 cid 进行推送，则会直接返回第一次成功推送的结果，不会再次进行推送。
// count 不传则默认为1。范围为[1, 1000]
// type: 取值：push(默认), schedule
func (j *JpushClient) PushCid(count int, typ string) (cidList []string, err error) {
	if count == 0 {
		count = 1
	}
	if typ == "" {
		typ = "push"
	}
	params := map[string]string{
		"count": fmt.Sprint(count),
		"type":  typ,
	}
	var out struct {
		CidList []string `json:"cidlist"`
	}
	err = j.Do(http.MethodGet, "/push/cid", params, &out)
	return out.CidList, err
}

// PushValidate 该 API 只用于验证推送调用是否能够成功，与推送 API 的区别在于：不向用户发送任何消息
func (j *JpushClient) PushValidate(push *Payload) (msgId string, err error) {
	var out struct {
		SendNo string `json:"sendno"`
		MsgId  string `json:"msg_id"`
	}
	err = j.Do(http.MethodPost, "/push/validate", push, &out)
	return out.MsgId, err
}

// Push 向某单个设备或者某设备列表推送一条通知、或者消息。
// 推送的内容只能是 JSON 表示的一个推送对象。
func (j *JpushClient) Push(push *Payload) (msgId string, err error) {
	var out struct {
		SendNo string `json:"sendno"`
		MsgId  string `json:"msg_id"`
	}
	err = j.Do(http.MethodPost, "/push", push, &out)
	return out.MsgId, err
}

// ScheduleCreate 创建计划任务
func (j *JpushClient) ScheduleCreate(push *SchedulePayload) (scheduleId string, err error) {
	var out struct {
		ScheduleId string `json:"schedule_id"`
		Name       string `json:"name"`
	}
	err = j.Do(http.MethodPost, "/schedules", push, &out)
	if err != nil {
		return "", err
	}
	return out.ScheduleId, nil
}

// ScheduleUpdate 修改指定的计划任务
// 更新操作可为 "name"，"enabled"、"trigger"或"push" 四项中的一项或多项。
// 不支持部分更新, 需要更新一整块
func (j *JpushClient) ScheduleUpdate(push *SchedulePayload) (out *SchedulePayload, err error) {
	err = j.Do(http.MethodPut, "/schedules", push, &out)
	return
}

// ScheduleDelete 删除制定的计划任务
func (j *JpushClient) ScheduleDelete(scheduleId string) error {
	return j.Do(http.MethodDelete, "/schedules/"+scheduleId, nil, nil)
}

// ScheduleList 获取有效的计划任务列表
func (j *JpushClient) ScheduleList(pageNo int) (list []SchedulePayload, err error) {
	var out struct {
		TotalCount int               `json:"total_count"`
		TotalPages int               `json:"total_pages"`
		Page       int               `json:"page"`
		Schedules  []SchedulePayload `json:"schedules"`
	}
	err = j.Do(http.MethodGet, "/schedules?page="+fmt.Sprint(pageNo), nil, &out)
	return out.Schedules, err
}

// ScheduleGet 获取指定的计划任务
func (j *JpushClient) ScheduleGet(scheduleId string) (schedule SchedulePayload, err error) {
	err = j.Do(http.MethodGet, "/schedules/"+scheduleId, nil, &schedule)
	return
}

// ScheduleGetMsgIds 获取定时任务对应的所有 msg_id
func (j *JpushClient) ScheduleGetMsgIds(scheduleId string) (msgIds []string, err error) {
	var out struct {
		MsgIds []struct {
			MsgId string `json:"msg_id"`
		} `json:"msgids"`
	}
	err = j.Do(http.MethodGet, "/schedules/"+scheduleId+"/msg_ids", nil, &out)
	if err != nil {
		return nil, err
	}
	for _, item := range out.MsgIds {
		if item.MsgId != "" {
			msgIds = append(msgIds, item.MsgId)
		}
	}
	return
}

type DeviceInfo struct {
	Tags   []string `json:"tags"`
	Alias  string   `json:"alias"`
	Mobile string   `json:"mobile"`
}

// DeviceGet 查询设备的别名与标签
func (j *JpushClient) DeviceGet(regId string) (info DeviceInfo, err error) {
	err = j.Do(http.MethodGet, DeviceHost+"/devices/"+regId, nil, &info)
	return
}

type TagSet struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// DeviceUpdateSet tags: 支持add, remove 或者空字符串。当tags参数为空字符串的时候，表示清空所有的 tags；add/remove 下是增加或删除指定的 tag；
// 一次 add/remove tag 的上限均为 100 个，且总长度均不能超过 1000 字节。
// 可以多次调用 API 设置，一个注册 id tag 上限为1000个，应用 tag 总数没有限制
type DeviceUpdateSet struct {
	Tags   interface{} `json:"tags,omitempty"`
	Alias  string      `json:"alias,omitempty"`
	Mobile string      `json:"mobile,omitempty"`
}

// DeviceSet 设置设备的别名与标签
func (j *JpushClient) DeviceSet(regId string, setInfo *DeviceUpdateSet) error {
	return j.Do(http.MethodPost, DeviceHost+"/devices/"+regId, setInfo, nil)
}

// AliasGet 查询别名
func (j *JpushClient) AliasGet(alias string) (regIds []string, err error) {
	var out struct {
		RegistrationIds []string `json:"registration_ids"`
	}
	err = j.Do(http.MethodGet, "/aliases/"+alias, nil, &out)
	if err != nil {
		return nil, err
	}
	return out.RegistrationIds, nil
}

// AliasDelete 删除别名
func (j *JpushClient) AliasDelete(alias string) error {
	return j.Do(http.MethodDelete, "/aliases/"+alias, nil, nil)
}

// TagList 查询标签列表
func (j *JpushClient) TagList() (tags []string, err error) {
	var out struct {
		Tags []string `json:"tags"`
	}
	err = j.Do(http.MethodGet, "/tags/", nil, &out)
	if err != nil {
		return nil, err
	}
	return out.Tags, nil
}

// IsTag 判断设备与标签绑定关系
func (j *JpushClient) IsTag(regId, tagId string) (ok bool, err error) {
	path := fmt.Sprintf("/tags/%s/registration_ids/%s", regId, tagId)
	var out struct {
		Result bool `json:"result"`
	}
	err = j.Do(http.MethodGet, path, nil, &out)
	return out.Result, err
}

type RegistrationIdSet struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

type TagUpdateSet struct {
	RegistrationIds RegistrationIdSet `json:"registration_ids"`
}

// TagUpdate 更新标签
func (j *JpushClient) TagUpdate(tag string, set *TagUpdateSet) (err error) {
	return j.Do(http.MethodPost, "/tags/"+tag, set, nil)
}

// TagDelete 删除标签
func (j *JpushClient) TagDelete(tag string) error {
	return j.Do(http.MethodDelete, "/tags/"+tag, nil, nil)
}

func (j *JpushClient) Do(method, path string, inp, out interface{}) error {
	var resp *grequests.Response
	var err error

	url := j.Url(path)
	if strings.HasPrefix(path, "http") {
		url = path
	}
	if method == http.MethodGet {
		var params = make(map[string]string)
		if val, ok := inp.(map[string]string); ok {
			params = val
		}

		resp, err = j.Get(url, &grequests.RequestOptions{
			Params: params,
		})
	} else if method == http.MethodPost {
		if val, ok := inp.(map[string]string); ok {
			resp, err = j.Post(url, &grequests.RequestOptions{
				Data: val,
			})
		} else if inp != nil {
			resp, err = j.Post(url, &grequests.RequestOptions{
				JSON: inp,
			})
		} else {
			resp, err = j.Post(url, nil)
		}
	} else if method == http.MethodPut {
		if val, ok := inp.(map[string]string); ok {
			resp, err = j.Put(url, &grequests.RequestOptions{
				Data: val,
			})
		} else if inp != nil {
			resp, err = j.Put(url, &grequests.RequestOptions{
				JSON: inp,
			})
		} else {
			resp, err = j.Put(url, nil)
		}
	} else if method == http.MethodDelete {
		if val, ok := inp.(map[string]string); ok {
			resp, err = j.Delete(url, &grequests.RequestOptions{
				Data: val,
			})
		} else if inp != nil {
			resp, err = j.Delete(url, &grequests.RequestOptions{
				JSON: inp,
			})
		} else {
			resp, err = j.Delete(url, nil)
		}
	}
	if err != nil {
		return err
	}

	// println(resp.String(), resp.StatusCode)

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		if strings.Contains(resp.String(), "error") {
			var er struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			err = resp.JSON(&er)
			if err != nil {
				return err
			}
			return fmt.Errorf("%d - %s", er.Error.Code, er.Error.Message)
		}
		return errors.New(resp.RawResponse.Status)
	}
	if out != nil {
		return resp.JSON(&out)
	}
	return nil
}
