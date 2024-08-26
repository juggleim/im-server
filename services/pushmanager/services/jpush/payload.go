package jpush

import (
	"encoding/json"
)

type Payload struct {
	Platform     *Platform     `json:"platform"`
	Audience     *Audience     `json:"audience"`
	Notification *Notification `json:"notification,omitempty"`
	Message      *Message      `json:"message,omitempty"`
	SmsMessage   *SmsMessage   `json:"sms_message,omitempty"`
	Options      *Options      `json:"options,omitempty"`
	Cid          string        `json:"cid,omitempty"`
}

func (p *Payload) MarshalJSON() ([]byte, error) {
	payload := struct {
		Platform     interface{}   `json:"platform"`
		Audience     interface{}   `json:"audience"`
		Notification *Notification `json:"notification,omitempty"`
		Message      *Message      `json:"message,omitempty"`
		SmsMessage   *SmsMessage   `json:"sms_message,omitempty"`
		Options      *Options      `json:"options,omitempty"`
		Cid          string        `json:"cid,omitempty"`
	}{
		Platform:     p.Platform.Interface(),
		Audience:     p.Audience.Interface(),
		Notification: p.Notification,
		Message:      p.Message,
		SmsMessage:   p.SmsMessage,
		Options:      p.Options,
		Cid:          p.Cid,
	}
	return json.Marshal(payload)
}

type Single struct {
	Time string `json:"time"`
}

const (
	Day   = "day"
	Week  = "week"
	Month = "month"
)

const (
	WeekMonday    = "MON"
	WeekTuesday   = "TUE"
	WeekWednesday = "WED"
	WeekThursday  = "THU"
	WeekFriday    = "FRI"
	WeekSaturday  = "SAT"
	WeekSunday    = "SUN"
)

type Periodical struct {
	Start     string   `json:"start"`
	End       string   `json:"end"`
	Time      string   `json:"time"`
	TimeUnit  string   `json:"time_unit"`
	Frequency int      `json:"frequency"`
	Point     []string `json:"point,omitempty"`
}

type Trigger struct {
	Single     *Single     `json:"single,omitempty"`
	Periodical *Periodical `json:"periodical,omitempty"`
}

type SchedulePayload struct {
	ScheduleId string   `json:"schedule_id,omitempty"`
	Name       string   `json:"name,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
	Trigger    *Trigger `json:"trigger,omitempty"`
	Push       *Payload `json:"push,omitempty"`
	Cid        string   `json:"cid,omitempty"`
}
