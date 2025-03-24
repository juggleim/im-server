package models

import "im-server/commons/pbdefines/pbobjs"

type HisMsg struct {
	ConverId    string
	SenderId    string
	ReceiverId  string
	ChannelType pbobjs.ChannelType
	MsgType     string
	MsgId       string
	SendTime    int64
	MsgSeqNo    int64
	MsgBody     []byte
	AppKey      string
	IsExt       int
	IsExset     int
	MsgExt      []byte
	MsgExset    []byte
	IsDelete    int
}

type GroupHisMsg struct {
	HisMsg
	MemberCount int
	ReadCount   int
}

type PrivateHisMsg struct {
	HisMsg
	IsRead int
}

type SystemHisMsg struct {
	HisMsg
	IsRead int
}

type ConverItem struct {
	ConverId string
	MsgId    string
}

type IGroupHisMsgStorage interface {
	SaveGroupHisMsg(msg GroupHisMsg) error
	//QryLatestMsgSeqNo(appkey, converId string) int64
	QryLatestMsg(appkey, converId string) (*GroupHisMsg, error)
	QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*GroupHisMsg, error)
	QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*GroupHisMsg, error)
	UpdateMsgBody(appkey, converId, msgId, msgType string, msgBody []byte) error
	FindById(appkey, converId, msgId string) (*GroupHisMsg, error)
	FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*GroupHisMsg, error)
	FindByConvers(appkey string, convers []ConverItem) ([]*GroupHisMsg, error)
	DelMsgs(appkey, converId string, msgIds []string) error
	UpdateMsgExtState(appkey, converId, msgId string, isExt int) error
	UpdateMsgExt(appkey, converId, msgId string, ext []byte) error
	UpdateMsgExsetState(appkey, converId, msgId string, isExset int) error
	UpdateMsgExset(appkey, converId, msgId string, ext []byte) error
	DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error

	UpdateReadCount(appkey, converId, msgId string, readCount int) error
}

type IPrivateHisMsgStorage interface {
	SavePrivateHisMsg(msg PrivateHisMsg) error
	//QryLatestMsg(appkey, converId string) *PrivateHisMsg
	QryLatestMsg(appkey, converId string) (*PrivateHisMsg, error)
	QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*PrivateHisMsg, error)
	QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*PrivateHisMsg, error)
	UpdateMsgBody(appkey, converId, msgId, msgType string, msgBody []byte) error
	FindById(appkey, converId, msgId string) (*PrivateHisMsg, error)
	FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*PrivateHisMsg, error)
	FindByConvers(appkey string, convers []ConverItem) ([]*PrivateHisMsg, error)
	DelMsgs(appkey, converId string, msgIds []string) error
	UpdateMsgExtState(appkey, converId, msgId string, isExt int) error
	UpdateMsgExt(appkey, converId, msgId string, ext []byte) error
	UpdateMsgExsetState(appkey, converId, msgId string, isExset int) error
	UpdateMsgExset(appkey, converId, msgId string, ext []byte) error
	DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error

	MarkReadByMsgIds(appkey, converId string, msgIds []string) error
	MarkReadByScope(appkey, converId string, start, end int64) error
}

type ISystemHisMsgStorage interface {
	SaveSystemHisMsg(msg SystemHisMsg) error
	// QryLatestMsgSeqNo(appkey, converId string) int64
	QryLatestMsg(appkey, converId string) (*SystemHisMsg, error)
	QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*SystemHisMsg, error)
	FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*SystemHisMsg, error)
}
type BrdCastHisMsg struct {
	ConverId    string
	SenderId    string
	ChannelType pbobjs.ChannelType
	MsgType     string
	MsgId       string
	SendTime    int64
	MsgSeqNo    int64
	MsgBody     []byte
	AppKey      string
}
type IBrdCastHisMsgStorage interface {
	SaveBrdCastHisMsg(msg BrdCastHisMsg) error
	// QryLatestMsgSeqNo(appkey, converId string) int64
	QryLatestMsg(appkey, converId string) (*BrdCastHisMsg, error)
	QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*BrdCastHisMsg, error)
	FindById(appkey, conver_id, msgId string) (*BrdCastHisMsg, error)
	FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*BrdCastHisMsg, error)
}
type GrpCastHisMsg struct {
	ConverId    string
	SenderId    string
	ReceiverId  string
	ChannelType pbobjs.ChannelType
	MsgType     string
	MsgId       string
	SendTime    int64
	MsgSeqNo    int64
	MsgBody     []byte
	AppKey      string
}
type IGrpCastHisMsgStorage interface {
	SaveGrpCastHisMsg(msg GrpCastHisMsg) error
	// QryLatestMsgSeqNo(appkey, converId string) int64
	QryLatestMsg(appkey, converId string) (*GrpCastHisMsg, error)
	QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*GrpCastHisMsg, error)
	FindById(appkey, conver_id, msgId string) (*GrpCastHisMsg, error)
	FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*GrpCastHisMsg, error)
}

type MergedMsg struct {
	ParentMsgId string
	FromId      string
	TargetId    string
	ChannelType pbobjs.ChannelType
	MsgId       string
	MsgTime     int64
	MsgBody     []byte
	AppKey      string
}

type IMergedMsgStorage interface {
	SaveMergedMsg(item MergedMsg) error
	BatchSaveMergedMsgs(items []MergedMsg) error
	QryMergedMsgs(appkey, parentMsgId string, startTime int64, count int32, isPositiveOrder bool) ([]*MergedMsg, error)
}

type HisMsgConverCleanTime struct {
	ConverId    string
	ChannelType pbobjs.ChannelType
	CleanTime   int64
	AppKey      string
}

type IHisMsgConverCleanTimeStorage interface {
	UpsertDestroyTime(item HisMsgConverCleanTime) error
	FindOne(appkey, converId string, channelType pbobjs.ChannelType) (*HisMsgConverCleanTime, error)
}

type HisMsgUserCleanTime struct {
	UserId      string
	TargetId    string
	ChannelType pbobjs.ChannelType
	CleanTime   int64
	AppKey      string
}

type IHisMsgUserCleanTimeStorage interface {
	UpsertCleanTime(item HisMsgUserCleanTime) error
	FindOne(appkey, userId, targetId string, channelType pbobjs.ChannelType) (*HisMsgUserCleanTime, error)
}

type GroupDelHisMsg struct {
	UserId   string
	TargetId string
	MsgId    string
	MsgTime  int64
	MsgSeq   int64
	AppKey   string
}

type IGroupDelHisMsgStorage interface {
	Create(item GroupDelHisMsg) error
	BatchCreate(items []GroupDelHisMsg) error
	QryDelHisMsgs(appkey, userId, targetId string, startTime int64, count int32, isPositive bool) ([]*GroupDelHisMsg, error)
}

type PrivateDelHisMsg struct {
	UserId   string
	TargetId string
	MsgId    string
	MsgTime  int64
	MsgSeq   int64
	AppKey   string
}

type IPrivateDelHisMsgStorage interface {
	Create(item PrivateDelHisMsg) error
	BatchCreate(items []PrivateDelHisMsg) error
	QryDelHisMsgs(appkey, userId, targetId string, startTime int64, count int32, isPositive bool) ([]*PrivateDelHisMsg, error)
}

type ReadInfo struct {
	AppKey      string
	MsgId       string
	ChannelType pbobjs.ChannelType
	GroupId     string
	MemberId    string
	CreatedTime int64
}

type IReadInfoStorage interface {
	Create(item ReadInfo) error
	BatchCreate(items []ReadInfo) error
	QryReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string, startId, limit int64) ([]*ReadInfo, error)
	CountReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string) int32
	CheckMsgsRead(appkey, groupId, memberId string, channelType pbobjs.ChannelType, msgIds []string) (map[string]bool, error)
}

type MsgExt struct {
	AppKey      string
	MsgId       string
	Key         string
	Value       string
	UserId      string
	CreatedTime int64
}

type IMsgExtStorage interface {
	Upsert(item MsgExt) error
	Delete(appkey, msgId, key string) error
	QryExtsByMsgIds(appkey string, msgIds []string) ([]*MsgExt, error)
}

type MsgExSet struct {
	AppKey      string
	MsgId       string
	Key         string
	Item        string
	UserId      string
	CreatedTime int64
}

type IMsgExSetStorage interface {
	Create(item MsgExSet) error
	Delete(appkey, msgId, key, item string) error
	QryExtsByMsgIds(appkey string, msgIds []string) ([]*MsgExSet, error)
}
