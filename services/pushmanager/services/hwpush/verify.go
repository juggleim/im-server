package hwpush

import (
	"errors"
	"regexp"
)

var (
	ttlPattern   = regexp.MustCompile("\\d+|\\d+[sS]|\\d+.\\d{1,9}|\\d+.\\d{1,9}[sS]")
	colorPattern = regexp.MustCompile("^#[0-9a-fA-F]{6}$")
)

func ValidateMessage(message *Message) error {
	if message == nil {
		return errors.New("message must not be null")
	}

	// validate field target, one of Token, Topic and Condition must be invoked
	if err := validateFieldTarget(message.Token, message.Topic, message.Condition); err != nil {
		return err
	}

	// validate android config
	if err := validateAndroidConfig(message.Android); err != nil {
		return err
	}

	// validate web common config
	if err := validateWebPushConfig(message.WebPush); err != nil {
		return err
	}
	return nil
}

func validateFieldTarget(token []string, strings ...string) error {
	count := 0
	if token != nil {
		count++
	}

	for _, s := range strings {
		if s != "" {
			count++
		}
	}

	if count == 1 {
		return nil
	}
	return errors.New("token, topics or condition must be choice one")
}

func validateWebPushConfig(webPushConfig *WebPushConfig) error {
	if webPushConfig == nil {
		return nil
	}

	if err := validateWebPushHeaders(webPushConfig.Headers); err != nil {
		return err
	}

	return validateWebPushNotification(webPushConfig.Notification)
}

func validateWebPushHeaders(headers *WebPushHeaders) error {
	if headers == nil {
		return nil
	}

	if headers.TTL != "" && !ttlPattern.MatchString(headers.TTL) {
		return errors.New("malformed ttl")
	}

	if headers.Urgency != "" &&
		headers.Urgency != UrgencyHigh &&
		headers.Urgency != UrgencyNormal &&
		headers.Urgency != UrgencyLow &&
		headers.Urgency != UrgencyVeryLow {
		return errors.New("priority must be 'high', 'normal', 'low' or 'very-low'")
	}
	return nil
}

func validateWebPushNotification(notification *WebPushNotification) error {
	if notification == nil {
		return nil
	}

	if err := validateWebPushAction(notification.Actions); err != nil {
		return err
	}

	if err := validateWebPushDirection(notification.Dir); err != nil {
		return err
	}
	return nil
}

func validateWebPushAction(actions []*WebPushAction) error {
	if actions == nil {
		return nil
	}

	for _, action := range actions {
		if action.Action == "" {
			return errors.New("web common action can't be empty")
		}
	}
	return nil
}

func validateWebPushDirection(dir string) error {
	if dir != DirAuto && dir != DirLtr && dir != DirRtl {
		return errors.New("web common dir must be 'auto', 'ltr', 'rtl'")
	}
	return nil
}

func validateAndroidConfig(androidConfig *AndroidConfig) error {
	if androidConfig == nil {
		return nil
	}

	if androidConfig.CollapseKey < -1 || androidConfig.CollapseKey > 100 {
		return errors.New("collapse_key must be in interval [-1 - 100]")
	}

	if androidConfig.Urgency != "" &&
		androidConfig.Urgency != DeliveryPriorityHigh &&
		androidConfig.Urgency != DeliveryPriorityNormal {
		return errors.New("delivery_priority must be 'HIGH' or 'NORMAL'")
	}

	if androidConfig.TTL != "" && !ttlPattern.MatchString(androidConfig.TTL) {
		return errors.New("malformed ttl")
	}

	if androidConfig.FastAppTarget != 0 &&
		androidConfig.FastAppTarget != FastAppTargetDevelop &&
		androidConfig.FastAppTarget != FastAppTargetProduct {
		return errors.New("invalid fast_app_target")
	}

	// validate android notification
	return validateAndroidNotification(androidConfig.Notification)
}

func validateAndroidNotification(notification *AndroidNotification) error {
	if notification == nil {
		return nil
	}

	if notification.Sound == "" && !notification.DefaultSound {
		return errors.New("sound must not be empty when default_sound is false")
	}

	if err := validateAndroidNotifyStyle(notification); err != nil {
		return err
	}

	if err := validateAndroidNotifyPriority(notification); err != nil {
		return err
	}

	if err := validateVibrateTimings(notification); err != nil {
		return err
	}

	if err := validateVisibility(notification); err != nil {
		return err
	}

	if err := validateLightSetting(notification); err != nil {
		return err
	}

	if notification.Color != "" && !colorPattern.MatchString(notification.Color) {
		return errors.New("color must be in the form #RRGGBB")
	}

	// validate click action
	return validateClickAction(notification.ClickAction)
}

func validateAndroidNotifyStyle(notification *AndroidNotification) error {
	switch notification.Style {
	case StyleBigText:
		if notification.BigTitle == "" {
			return errors.New("big_title must not be empty when style is 1")
		}

		if notification.BigBody == "" {
			return errors.New("big_body must not be empty when style is 1")
		}

	}
	return nil
}

func validateAndroidNotifyPriority(notification *AndroidNotification) error {
	if notification.Importance != "" &&
		notification.Importance != NotificationPriorityHigh &&
		notification.Importance != NotificationPriorityDefault &&
		notification.Importance != NotificationPriorityLow {
		return errors.New("Importance must be 'HIGH', 'NORMAL' or 'LOW'")
	}
	return nil
}

func validateVibrateTimings(notification *AndroidNotification) error {
	if notification.VibrateConfig != nil {
		if len(notification.VibrateConfig) > 10 {
			return errors.New("vibrate_timings can't be more than 10 elements")
		}
		for _, vibrateTiming := range notification.VibrateConfig {
			if !ttlPattern.MatchString(vibrateTiming) {
				return errors.New("malformed vibrate_timings")
			}
		}
	}
	return nil
}

func validateVisibility(notification *AndroidNotification) error {
	if notification.Visibility == "" {
		notification.Visibility = VisibilityPrivate
		return nil
	}
	if notification.Visibility != VisibilityUnspecified && notification.Visibility != VisibilityPrivate &&
		notification.Visibility != VisibilityPublic && notification.Visibility != VisibilitySecret {
		return errors.New("visibility must be VISIBILITY_UNSPECIFIED, PRIVATE, PUBLIC or SECRET")
	}
	return nil
}

func validateLightSetting(notification *AndroidNotification) error {
	if notification.LightSettings == nil {
		return nil
	}

	if notification.LightSettings.Color == nil {
		return errors.New("light_settings.color can't be nil")
	}

	if notification.LightSettings.LightOnDuration == "" ||
		!ttlPattern.MatchString(notification.LightSettings.LightOnDuration) {
		return errors.New("light_settings.light_on_duration is empty or malformed")
	}

	if notification.LightSettings.LightOffDuration == "" ||
		!ttlPattern.MatchString(notification.LightSettings.LightOffDuration) {
		return errors.New("light_settings.light_off_duration is empty or malformed")
	}
	return nil
}

func validateClickAction(clickAction *ClickAction) error {
	if clickAction == nil {
		return errors.New("click_action object must not be null")
	}

	switch clickAction.Type {
	case TypeIntentOrAction:
		if clickAction.Intent == "" && clickAction.Action == "" {
			return errors.New("at least one of intent and action is not empty when type is 1")
		}
	case TypeUrl:
		if clickAction.Url == "" {
			return errors.New("url must not be empty when type is 2")
		}
	case TypeApp:
	case TypeRichResource:
		if clickAction.RichResource == "" {
			return errors.New("rich_resource must not be empty when type is 4")
		}
	default:
		return errors.New("type must be in the interval [1 - 4]")
	}
	return nil
}
