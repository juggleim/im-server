package models

type GrpApplicationStatus int

var (
	GrpApplicationStatus_Apply        = 0
	GrpApplicationStatus_AgreeApply   = 1
	GrpApplicationStatus_DeclineApply = 2
	GrpApplicationStatus_ExpiredApply = 3

	GrpApplicationStatus_Invite        = 10
	GrpApplicationStatus_AgreeInvite   = 11
	GrpApplicationStatus_DeclineInvite = 12
	GrpApplicationStatus_ExpiredInvite = 13
)

type GrpApplicationType int

var (
	GrpApplicationType_Invite = 0
	GrpApplicationType_Apply  = 1
)

type GrpApplication struct {
	ID          int64
	GroupId     string
	ApplyType   GrpApplicationType
	SponsorId   string
	RecipientId string
	InviterId   string
	OperatorId  string
	ApplyTime   int64
	Status      GrpApplicationStatus
	AppKey      string
}

type IGrpApplicationStorage interface {
	InviteUpsert(item GrpApplication) error
	ApplyUpsert(item GrpApplication) error
	QueryMyGrpApplications(appkey, sponsorId string, startTime, count int64, isPositive bool) ([]*GrpApplication, error)
	QueryMyPendingGrpInvitations(appkey, recipientId string, startTime, count int64, isPositive bool) ([]*GrpApplication, error)
	QueryGrpInvitations(appkey, groupId string, startTime, count int64, isPositive bool) ([]*GrpApplication, error)
	QueryGrpPendingApplications(appkey, groupId string, startTime, count int64, isPositive bool) ([]*GrpApplication, error)
}
