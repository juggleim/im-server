package models

import "im-server/services/commonservices"

type BotInfo struct {
	BotId       string                      `json:"bot_id"`
	Nickname    string                      `json:"nickname"`
	Portrait    string                      `json:"portrait"`
	BotConf     *commonservices.BotConf     `json:"bot_conf,omitempty"`
	BotSettings *commonservices.BotSettings `json:"bot_settings,omitempty"`
	ExtFields   map[string]string           `json:"ext_fields"`
	UpdatedTime int64                       `json:"updated_time"`
}
