package models

// UserSubRel 用户状态订阅关系（对外模型，不含持久化细节）
type UserSubRel struct {
	ID                 int64
	UserId             string
	SubscriberId       string
	SubscriberDeviceId string
	CreatedTime        int64
	AppKey             string
}

type IUserSubRelStorage interface {
	BatchCreate(items []*UserSubRel) error
	QryBySubscriber(appkey, subscriberId, deviceId string, limit int) ([]*UserSubRel, error)
	// QryByUserID 按目标用户（user_id）查询订阅关系，id > afterID，按 id 升序，用于分页扫表
	QryByUserID(appkey, targetUserId string, afterID int64, limit int) ([]*UserSubRel, error)
	Delete(appkey, userId, subscriberId, deviceId string) error
	DeleteByRelIDs(relIds []int64) error
}
