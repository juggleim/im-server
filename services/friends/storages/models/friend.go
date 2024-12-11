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
}

type FriendApplication struct {
	ID        int64
	UserId    string
	SponsorId string
	ApplyTime int64
	Status    int
	AppKey    string
}
