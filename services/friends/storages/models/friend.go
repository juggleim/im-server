package models

type FriendRel struct {
	ID       int64
	AppKey   string
	UserId   string
	FriendId string
}

type IFriendRelStorage interface {
	Upsert(item FriendRel) error
	BatchUpsert(items []FriendRel) error
	QueryFriendRels(appkey, userId string, startId, limit int64) ([]*FriendRel, error)
	BatchDelete(appkey, userId string, friendIds []string) error
	QueryFriendRelsByFriendIds(appkey, userId string, friendIds []string) ([]*FriendRel, error)
}

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
