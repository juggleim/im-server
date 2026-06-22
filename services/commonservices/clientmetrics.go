package commonservices

type ClientConnectMetrics struct {
	OnlineUserCount     int64 `json:"online_user_count"`
	UserConnectCount    int64 `json:"user_connect_count"`
	SessionConnectCount int64 `json:"session_connect_count"`
}

var clientConnectMetricsProvider func() ClientConnectMetrics

func RegisterClientConnectMetricsProvider(provider func() ClientConnectMetrics) {
	clientConnectMetricsProvider = provider
}

func GetClientConnectMetrics() ClientConnectMetrics {
	if clientConnectMetricsProvider == nil {
		return ClientConnectMetrics{}
	}
	return clientConnectMetricsProvider()
}
