package sms

var (
	DefaultSmsEngine ISmsEngine = &NilSmsEngine{}
)

type ISmsEngine interface {
	SmsSend(phone, content string) error
}

type NilSmsEngine struct{}

func (engine *NilSmsEngine) SmsSend(phone, content string) error {
	return nil
}
