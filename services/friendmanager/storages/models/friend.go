package models

type FriendRel struct {
	ID          int64
	AppKey      string
	UserId      string
	FriendId    string
	DisplayName string
	OrderTag    string
	CreatedTime int64
	UpdatedTime int64
}

type IFriendRelStorage interface {
	Upsert(item FriendRel) error
	BatchUpsert(items []FriendRel) error
	QueryFriendRels(appkey, userId string, startId, limit int64, isPositive bool) ([]*FriendRel, error)
	QueryFriendRelsWithPage(appkey, userId string, orderTag string, page, size int64) ([]*FriendRel, error)
	BatchDelete(appkey, userId string, friendIds []string) error
	GetFriendRel(appkey, userId, friendId string) (*FriendRel, error)
	QueryFriendRelsByFriendIds(appkey, userId string, friendIds []string) ([]*FriendRel, error)
	UpdateOrderTag(appkey, friendId string, orderTag string) error
	UpdateDisplayName(appkey, userId, friendId, displayName string) error
}
