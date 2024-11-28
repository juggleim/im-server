package transengines

var (
	DefaultTransEngine ITransEngine = &NilTransEngine{}
)

type ITransEngine interface {
	Translate(content string, targetLanguages []string) map[string]string
}

type NilTransEngine struct{}

func (engine *NilTransEngine) Translate(content string, targetLanguages []string) map[string]string {
	return map[string]string{}
}
