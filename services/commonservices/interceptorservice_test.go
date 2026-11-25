package commonservices

import (
	"testing"

	"im-server/services/commonservices/interceptors"
)

func TestShouldSkipApiMsgCallback(t *testing.T) {
	customCallback := &interceptors.MsgInterceptor{Interceptor: &interceptors.CustomInterceptor{}}
	contentReview := &interceptors.MsgInterceptor{Interceptor: &interceptors.BdInterceptor{}}

	tests := []struct {
		name        string
		isFromAPI   bool
		isSubAPIMsg bool
		senderID    string
		whitelist   string
		interceptor *interceptors.MsgInterceptor
		want        bool
	}{
		{name: "API callback disabled", isFromAPI: true, interceptor: customCallback, want: true},
		{name: "API callback enabled without whitelist", isFromAPI: true, isSubAPIMsg: true, senderID: "user1", interceptor: customCallback, want: false},
		{name: "API sender in JSON whitelist", isFromAPI: true, isSubAPIMsg: true, senderID: "user1", whitelist: `["user1","user2"]`, interceptor: customCallback, want: false},
		{name: "API sender outside JSON whitelist", isFromAPI: true, isSubAPIMsg: true, senderID: "user3", whitelist: `["user1","user2"]`, interceptor: customCallback, want: true},
		{name: "API empty JSON whitelist denies all", isFromAPI: true, isSubAPIMsg: true, senderID: "user1", whitelist: `[]`, interceptor: customCallback, want: true},
		{name: "API sender in comma-separated whitelist", isFromAPI: true, isSubAPIMsg: true, senderID: "user2", whitelist: "user1, user2", interceptor: customCallback, want: false},
		{name: "API sender outside comma-separated whitelist", isFromAPI: true, isSubAPIMsg: true, senderID: "user3", whitelist: "user1, user2", interceptor: customCallback, want: true},
		{name: "client message", interceptor: customCallback, want: false},
		{name: "API content review", isFromAPI: true, interceptor: contentReview, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appInfo := &AppInfo{
				IsSubApiMsg:                    test.isSubAPIMsg,
				ApiMsgPreCallbackUserWhitelist: test.whitelist,
			}
			appInfo.loadApiMsgPreCallbackUserWhitelist()
			if got := shouldSkipApiMsgCallback(test.isFromAPI, test.senderID, appInfo, test.interceptor); got != test.want {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}
