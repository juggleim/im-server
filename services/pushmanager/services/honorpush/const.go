package honorpush

// 接口说明见：https://developer.honor.com/cn/docs/11002/reference/downlink-message
const (
	// IAMHost 获取鉴权 Token
	IAMHost = "https://iam.developer.honor.com"
	// AuthTokenPath 应用级 Access Token（client_credentials）
	AuthTokenPath = "/auth/token"

	// PushAPIHost 推送服务端点
	PushAPIHost = "https://push-api.cloud.honor.com"
)
