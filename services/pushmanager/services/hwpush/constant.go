package hwpush

const (
	//the parameters of the formats below are endpoint and appId
	SendMessageFmt = "%s/v1/%s/messages:send"
)

const (
	// unspecified visibility
	VisibilityUnspecified = "VISIBILITY_UNSPECIFIED"
	// private visibility
	VisibilityPrivate = "PRIVATE"
	// public visibility
	VisibilityPublic = "PUBLIC"
	// secret visibility
	VisibilitySecret = "SECRET"
)

const (
	// high priority
	DeliveryPriorityHigh = "HIGH"
	// normal priority
	DeliveryPriorityNormal = "NORMAL"
)

const (
	// high priority
	NotificationPriorityHigh = "HIGH"
	// default priority
	NotificationPriorityDefault = "NORMAL"
	// low priority
	NotificationPriorityLow = "LOW"
)

const (
	// very low urgency
	UrgencyVeryLow = "very-low"
	// low urgency
	UrgencyLow = "low"
	// normal urgency
	UrgencyNormal = "normal"
	// high urgency
	UrgencyHigh = "high"
)

const (
	// webPush text direction auto
	DirAuto = "auto"
	// webPush text direction ltr
	DirLtr = "ltr"
	// webPush text direction rtl
	DirRtl = "rtl"
)

const (
	// success code from push server
	Success = "80000000"
	// parameter invalid code from push server
	ParameterError = "80100001"
	// token invalid code from push server
	TokenFailedErr = "80200001"
	//token timeout code from push server
	TokenTimeoutErr = "80200003"
)

const (
	StyleBigText = iota + 1
)

const (
	TypeIntentOrAction = iota + 1
	TypeUrl
	TypeApp
	TypeRichResource
)

const (
	FastAppTargetDevelop = iota + 1
	FastAppTargetProduct
)

const (
	// test user
	TargetUserTypeTest = iota + 1
	// formal user
	TargetUserTypeFormal
	// VoIP user
	TargetUserTypeVoIP
)
