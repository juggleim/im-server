package apnspush

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/sideshow/apns2"
)

func TestRouting(t *testing.T) {
	ordinary := &apns2.Client{HTTPClient: &http.Client{}}
	voip := &apns2.Client{HTTPClient: &http.Client{}}
	for _, tc := range []struct {
		name           string
		clients        *Clients
		isVoip         bool
		ordinary, voip string
		wantClient     *apns2.Client
		wantType       apns2.EPushType
		wantErr        error
	}{
		{"ordinary", &Clients{ordinary, voip}, false, "aabb", "ccdd", ordinary, apns2.PushTypeAlert, nil},
		{"voip", &Clients{ordinary, voip}, true, "aabb", "ccdd", voip, apns2.PushTypeVOIP, nil},
		{"fallback-no-client", &Clients{ordinary, nil}, true, "aabb", "ccdd", ordinary, apns2.PushTypeAlert, nil},
		{"fallback-no-token", &Clients{ordinary, voip}, true, "aabb", "", ordinary, apns2.PushTypeAlert, nil},
		{"fallback-empty-ordinary", &Clients{ordinary, voip}, true, "", "", nil, "", ErrInvalidToken},
		{"no-client", &Clients{}, false, "aabb", "", nil, "", ErrConfiguration},
		{"bad-selected-voip", &Clients{ordinary, voip}, true, "aabb", "not-hex", nil, "", ErrInvalidToken},
		{"bad-unselected-ordinary", &Clients{ordinary, voip}, true, "bad-token", "ccdd", voip, apns2.PushTypeVOIP, nil},
		{"bad-unselected-voip", &Clients{ordinary, voip}, false, "aabb", "bad-token", ordinary, apns2.PushTypeAlert, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := &apns2.Notification{Topic: "old", DeviceToken: "old", PushType: apns2.PushTypeBackground, Payload: `{}`}
			client, n, err := Route(tc.clients, "com.example", original, tc.isVoip, tc.ordinary, tc.voip)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("route error = %v", err)
			}
			if err == nil {
				if client != tc.wantClient || n.PushType != tc.wantType {
					t.Fatal("wrong route")
				}
				wantTopic, wantToken := "com.example", tc.ordinary
				if tc.wantType == apns2.PushTypeVOIP {
					wantTopic, wantToken = wantTopic+".voip", tc.voip
				}
				if n.Topic != wantTopic || n.DeviceToken != wantToken {
					t.Fatal("wrong topic/token")
				}
			}
			if original.Topic != "old" || original.DeviceToken != "old" || original.PushType != apns2.PushTypeBackground {
				t.Fatal("caller notification mutated")
			}
		})
	}
	for _, value := range []string{"", "abc", " aa ", "aa/bb", "aa?bb", "aa#bb", "aagg", "\r\n", "ＡＡ"} {
		if ValidDeviceToken(value) {
			t.Errorf("invalid token accepted: %q", value)
		}
	}
	for _, value := range []string{"aA", strings.Repeat("ab", 32), strings.Repeat("cd", 64)} {
		if !ValidDeviceToken(value) {
			t.Error("valid variable-length hex token rejected")
		}
	}
}

func TestClassification(t *testing.T) {
	for _, tc := range []struct {
		status int
		reason string
		class  ErrorClass
	}{
		{200, "", ClassNone},
		{403, apns2.ReasonExpiredProviderToken, ClassAuthentication},
		{403, apns2.ReasonInvalidProviderToken, ClassAuthentication},
		{403, apns2.ReasonMissingProviderToken, ClassAuthentication},
		{400, apns2.ReasonTopicDisallowed, ClassAuthentication},
		{400, apns2.ReasonBadDeviceToken, ClassDevice},
		{400, apns2.ReasonDeviceTokenNotForTopic, ClassDevice},
		{410, apns2.ReasonUnregistered, ClassDevice},
		{429, apns2.ReasonTooManyProviderTokenUpdates, ClassThrottling},
		{429, apns2.ReasonTooManyRequests, ClassThrottling},
		{500, apns2.ReasonInternalServerError, ClassUpstream},
		{503, apns2.ReasonShutdown, ClassUpstream},
		{400, apns2.ReasonIdleTimeout, ClassUpstream},
		{400, apns2.ReasonBadTopic, ClassConfiguration},
		{413, apns2.ReasonPayloadTooLarge, ClassConfiguration},
		{403, "new-reason", ClassAuthentication},
		{410, "new-reason", ClassDevice},
		{429, "new-reason", ClassThrottling},
		{502, "new-reason", ClassUpstream},
		{302, "new-reason", ClassUpstream},
	} {
		if got := Classify(&apns2.Response{StatusCode: tc.status, Reason: tc.reason}, nil); got != tc.class {
			t.Errorf("%s/%d: %s != %s", tc.reason, tc.status, got, tc.class)
		}
	}
	for _, tc := range []struct {
		err   error
		class ErrorClass
	}{
		{ErrNotFound, ClassConfiguration}, {ErrConfiguration, ClassConfiguration}, {ErrLoad, ClassConfiguration}, {ErrClosed, ClassConfiguration},
		{ErrInvalidToken, ClassDevice}, {context.Canceled, ClassNetwork}, {context.DeadlineExceeded, ClassNetwork}, {errors.New("dial failed"), ClassNetwork},
	} {
		if got := Classify(nil, fmt.Errorf("wrapped: %w", tc.err)); got != tc.class {
			t.Errorf("%v: %s", tc.err, got)
		}
	}
	if Classify(nil, nil) != ClassUpstream {
		t.Fatal("missing response treated as success")
	}
	if !strings.Contains(Diagnostic(&apns2.Response{StatusCode: 403, Reason: apns2.ReasonExpiredProviderToken}, nil), "clock") {
		t.Fatal("expired-token diagnostic missing clock guidance")
	}
}

func TestTimeoutEnvironment(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  time.Duration
	}{
		{"", 10 * time.Second}, {"0", 10 * time.Second}, {"-1", 10 * time.Second}, {"121", 10 * time.Second}, {"abc", 10 * time.Second}, {"1.5", 10 * time.Second}, {" 5", 10 * time.Second},
		{"999999999999999999999999999999999", 10 * time.Second}, {"1", time.Second}, {"120", 120 * time.Second},
	} {
		t.Setenv("IM_APNS_PUSH_TIMEOUT_SECONDS", tc.value)
		if got := PushTimeoutFromEnv(); got != tc.want {
			t.Errorf("%q: %v != %v", tc.value, got, tc.want)
		}
	}
}
