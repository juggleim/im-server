package serversdk

const (
	DefaultApiUrl string = ""
)

type JuggleIMSdk struct {
	Appkey string
	Secret string
	ApiUrl string
}

func NewJuggleIMSdk(appkey, secret, apiUrl string) *JuggleIMSdk {
	return &JuggleIMSdk{
		Appkey: appkey,
		Secret: secret,
		ApiUrl: apiUrl,
	}
}
