package models

type GrpApplicationStatus int

var (
	GrpApplicationStatus_Apply        GrpApplicationStatus = 0
	GrpApplicationStatus_AgreeApply   GrpApplicationStatus = 1
	GrpApplicationStatus_DeclineApply GrpApplicationStatus = 2
	GrpApplicationStatus_ExpiredApply GrpApplicationStatus = 3

	GrpApplicationStatus_Invite        GrpApplicationStatus = 10
	GrpApplicationStatus_AgreeInvite   GrpApplicationStatus = 11
	GrpApplicationStatus_DeclineInvite GrpApplicationStatus = 12
	GrpApplicationStatus_ExpiredInvite GrpApplicationStatus = 13
)

type GrpApplicationType int

var (
	GrpApplicationType_Invite GrpApplicationType = 0
	GrpApplicationType_Apply  GrpApplicationType = 1
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
