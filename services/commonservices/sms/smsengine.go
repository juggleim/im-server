package sms

var (
	DefaultSmsEngine ISmsEngine = &NilSmsEngine{}
)

type ISmsEngine interface {
	SmsSend(phone string, params map[string]interface{}) error
}

type NilSmsEngine struct{}

func (engine *NilSmsEngine) SmsSend(phone string, params map[string]interface{}) error {
	return nil
}
