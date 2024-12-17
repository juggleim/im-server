package models

type FriendRel struct {
	ID       int64
	AppKey   string
	UserId   string
	FriendId string
	OrderTag string
}

type IFriendRelStorage interface {
	Upsert(item FriendRel) error
	BatchUpsert(items []FriendRel) error
	QueryFriendRels(appkey, userId string, startId, limit int64) ([]*FriendRel, error)
	BatchDelete(appkey, userId string, friendIds []string) error
	QueryFriendRelsByFriendIds(appkey, userId string, friendIds []string) ([]*FriendRel, error)
}
