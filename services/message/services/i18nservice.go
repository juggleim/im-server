package services

import (
	"embed"
	"im-server/commons/logs"

	"github.com/kataras/i18n"
)

var (
	PlaceholderKey_Text  string = "placeholder_text"
	PlaceholderKey_Image string = "placeholder_image"
	PlaceholderKey_Voice string = "placeholder_voice"
	PlaceholderKey_File  string = "placeholder_file"
	PlaceholderKey_Video string = "placeholder_video"
	PlaceholderKey_Merge string = "placeholder_merge"
)

//go:embed locales/*
var i18nFs embed.FS

var i18nClient *i18n.I18n

func init() {
	loader, err := i18n.FS(i18nFs, "./locales/*/*.yml")
	if err == nil {
		client, err := i18n.New(loader, "en_US", "zh_CN")
		if err == nil {
			i18nClient = client
		} else {
			logs.Error("failed to create i18n client")
		}
	} else {
		logs.Error("failed to load i18n files")
	}
}

func GetI18nStr(language, key, defaultStr string) string {
	if i18nClient != nil {
		return i18nClient.Tr(language, key)
	}
	return defaultStr
}
