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
		interceptor *interceptors.MsgInterceptor
		want        bool
	}{
		{name: "API callback disabled", isFromAPI: true, interceptor: customCallback, want: true},
		{name: "API callback enabled", isFromAPI: true, isSubAPIMsg: true, interceptor: customCallback, want: false},
		{name: "client message", interceptor: customCallback, want: false},
		{name: "API content review", isFromAPI: true, interceptor: contentReview, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			appInfo := &AppInfo{IsSubApiMsg: test.isSubAPIMsg}
			if got := shouldSkipApiMsgCallback(test.isFromAPI, appInfo, test.interceptor); got != test.want {
				t.Fatalf("got %v, want %v", got, test.want)
			}
		})
	}
}
