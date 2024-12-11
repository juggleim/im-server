package models

type FriendApplicationStatus int

var (
	FriendApplicationStatus_Apply   = 0
	FriendApplicationStatus_Agree   = 1
	FriendApplicationStatus_Decline = 2
	FriendApplicationStatus_Expired = 3
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
}
