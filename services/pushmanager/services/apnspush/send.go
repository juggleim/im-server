package apnspush

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sideshow/apns2"
)

var (
	ErrNotFound      = errors.New("APNs configuration not found")
	ErrConfiguration = errors.New("invalid APNs configuration")
	ErrLoad          = errors.New("APNs configuration temporarily unavailable")
	ErrClosed        = errors.New("APNs manager closed")
	ErrInvalidToken  = errors.New("invalid APNs device token")
)

// Route returns a shallow notification copy, leaving the caller's notification
// unchanged. The caller must not mutate its payload while Send is in progress.
// An absent VoIP client or empty VoIP token preserves the legacy alert fallback.
// A nonempty but malformed selected token is rejected, not silently replaced.
func Route(clients *Clients, packageName string, notification *apns2.Notification, isVoip bool, ordinaryToken, voipToken string) (*apns2.Client, *apns2.Notification, error) {
	if clients == nil || notification == nil || !validTopic(packageName) {
		return nil, nil, fmt.Errorf("%w: missing client, notification or topic", ErrConfiguration)
	}
	n := *notification
	client := clients.ApnsClient
	n.Topic, n.DeviceToken, n.PushType = packageName, ordinaryToken, apns2.PushTypeAlert
	if isVoip && clients.ApnsVoipClient != nil && voipToken != "" {
		client = clients.ApnsVoipClient
		n.Topic, n.DeviceToken, n.PushType = packageName+".voip", voipToken, apns2.PushTypeVOIP
	}
	if !ValidDeviceToken(n.DeviceToken) {
		return nil, nil, ErrInvalidToken
	}
	if client == nil || client.HTTPClient == nil {
		return nil, nil, fmt.Errorf("%w: selected client unavailable", ErrConfiguration)
	}
	return client, &n, nil
}

// ValidDeviceToken checks the hex wire representation without assuming Apple's
// token byte length is fixed. It also prevents path/query injection into the SDK.
func ValidDeviceToken(value string) bool {
	if len(value) == 0 || len(value)%2 != 0 {
		return false
	}
	for _, c := range value {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

type ErrorClass string

const (
	ClassNone           ErrorClass = ""
	ClassConfiguration  ErrorClass = "config"
	ClassAuthentication ErrorClass = "auth"
	ClassDevice         ErrorClass = "device"
	ClassThrottling     ErrorClass = "throttling"
	ClassUpstream       ErrorClass = "upstream"
	ClassNetwork        ErrorClass = "network"
)

// Classify accepts the unmodified Send response/error pair. APNs rejections
// retain SDK semantics (non-nil response, nil error). No automatic retry occurs.
func Classify(response *apns2.Response, err error) ErrorClass {
	if err != nil {
		var marshalerError *json.MarshalerError
		var typeError *json.UnsupportedTypeError
		var valueError *json.UnsupportedValueError
		switch {
		case errors.Is(err, ErrInvalidToken):
			return ClassDevice
		case errors.Is(err, ErrConfiguration), errors.Is(err, ErrNotFound), errors.Is(err, ErrLoad), errors.Is(err, ErrClosed):
			return ClassConfiguration
		case errors.As(err, &marshalerError), errors.As(err, &typeError), errors.As(err, &valueError):
			return ClassConfiguration
		default:
			return ClassNetwork
		}
	}
	if response == nil {
		return ClassUpstream
	}
	if response.Sent() {
		return ClassNone
	}
	switch response.Reason {
	case apns2.ReasonBadDeviceToken, apns2.ReasonDeviceTokenNotForTopic, apns2.ReasonMissingDeviceToken, apns2.ReasonUnregistered:
		return ClassDevice
	case apns2.ReasonExpiredProviderToken, apns2.ReasonForbidden, apns2.ReasonInvalidProviderToken, apns2.ReasonMissingProviderToken, apns2.ReasonTopicDisallowed:
		return ClassAuthentication
	case apns2.ReasonTooManyProviderTokenUpdates, apns2.ReasonTooManyRequests:
		return ClassThrottling
	case apns2.ReasonIdleTimeout, apns2.ReasonInternalServerError, apns2.ReasonServiceUnavailable, apns2.ReasonShutdown:
		return ClassUpstream
	}
	switch {
	case response.StatusCode == 403:
		return ClassAuthentication
	case response.StatusCode == 410:
		return ClassDevice
	case response.StatusCode == 429:
		return ClassThrottling
	case response.StatusCode >= 500:
		return ClassUpstream
	case response.StatusCode >= 400 && response.StatusCode < 500:
		return ClassConfiguration
	default:
		return ClassUpstream
	}
}

// Diagnostic returns only static, credential-free guidance for safe logging.
func Diagnostic(response *apns2.Response, err error) string {
	if err == nil && response != nil && response.Reason == apns2.ReasonExpiredProviderToken {
		return "check node clock and provider JWT health; routine refresh belongs to the APNs SDK"
	}
	return string(Classify(response, err))
}

// SDK network errors can contain a URL with a device token. Keep errors.Is/As
// useful for callers, but never include that URL in the printable message.
type sendError struct{ cause error }

func (e *sendError) Error() string { return "APNs transport request failed" }
func (e *sendError) Unwrap() error { return e.cause }
