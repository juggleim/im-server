package services

import (
	"bytes"
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/logs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var tempParamsReg *regexp.Regexp
var tempCache *caches.LruCache

func init() {
	tempParamsReg, _ = regexp.Compile("{{.([^}]*)}}")
	tempCache = caches.NewLruCacheWithAddReadTimeout("temp_cache", 10000, nil, 10*time.Minute, 10*time.Minute)
}

func FetchTemplateParams(template string) []string {
	if tempParamsReg == nil {
		return []string{}
	}
	params := tempParamsReg.FindAllString(template, -1)
	for i := 0; i < len(params); i++ {
		if len(params[i]) > 5 {
			params[i] = params[i][3:]
			params[i] = params[i][:len(params[i])-2]
		}
	}
	return params
}

func TemplateI18nAssign(ctx context.Context, template, lang string) string {
	appkey := bases.GetAppKeyFromCtx(ctx)
	params := FetchTemplateParams(template)
	if len(params) > 0 {
		templateExecutor := getTextTemplate(appkey, template)
		if templateExecutor != nil {
			paramValueMap := map[string]string{}
			for _, param := range params {
				paramValueMap[param] = commonservices.GetI18nValue(appkey, lang, param, "")
			}
			buf := bytes.NewBuffer([]byte{})
			err := templateExecutor.Execute(buf, paramValueMap)
			if err == nil {
				return buf.String()
			}
		}
	}
	return template
}

func getTextTemplate(appkey, templateStr string) *template.Template {
	tempHash := tools.SHA1(templateStr)
	key := strings.Join([]string{appkey, tempHash}, "_")
	t, exist := tempCache.GetByCreator(key, func() interface{} {
		t := template.New(key)
		_, err := t.Parse(templateStr)
		if err != nil {
			logs.Errorf("failed init text template[err:%v,template:%s]", err, templateStr)
		}
		return t
	})
	if exist && t != nil {
		return t.(*template.Template)
	}
	return nil
}
