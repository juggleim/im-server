package models

type FriendApplicationStatus int

var (
	FriendApplicationStatus_Apply   FriendApplicationStatus = 0
	FriendApplicationStatus_Agree   FriendApplicationStatus = 1
	FriendApplicationStatus_Decline FriendApplicationStatus = 2
	FriendApplicationStatus_Expired FriendApplicationStatus = 3
)

type FriendApplication struct {
	ID          int64
	RecipientId string
	SponsorId   string
	ApplyTime   int64
	Status      FriendApplicationStatus
	AppKey      string
}

type IFriendApplicationStorage interface {
	Upsert(item FriendApplication) error
	QueryPendingApplications(appkey, recipientId string, startTime, count int64, isPositive bool) ([]*FriendApplication, error)
	QueryMyApplications(appkey, sponsorId string, startTime, count int64, isPositive bool) ([]*FriendApplication, error)
	QueryApplications(appkey, userId string, startTime, count int64, isPositive bool) ([]*FriendApplication, error)
	UpdateStatus(appkey, sponsorId, recipientId string, status FriendApplicationStatus) error
}
