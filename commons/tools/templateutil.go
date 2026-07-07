package tools

import (
	"io"
	"regexp"

	"github.com/valyala/fasttemplate"
)

var templateParamReg = regexp.MustCompile(`\$\{[^{}]+\}`)

func ReplaceTemplateParams(template string, params map[string]string) string {
	if templateParamReg.MatchString(template) {
		return fasttemplate.ExecuteFuncString(template, "${", "}", func(w io.Writer, tag string) (int, error) {
			if value, exist := params[tag]; exist {
				return w.Write([]byte(value))
			}
			return w.Write([]byte("${" + tag + "}"))
		})
	}
	return template
}
