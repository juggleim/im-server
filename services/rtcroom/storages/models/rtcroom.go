package models

import "im-server/commons/pbdefines/pbobjs"

type RtcRoom struct {
	RoomId       string
	RoomType     pbobjs.RtcRoomType
	RtcChannel   pbobjs.RtcChannel
	RtcMediaType pbobjs.RtcMediaType
	CreatedTime  int64
	AcceptedTime int64
	OwnerId      string
	Ext          string
	ConverId     *string
	ChannelType  pbobjs.ChannelType
	SubChannel   string
	AppKey       string
}

type IRtcRoomStorage interface {
	Create(item RtcRoom) error
	FindById(appkey, roomId string) (*RtcRoom, error)
	FindByConver(appkey, conveId string, channelType pbobjs.ChannelType, subChannel string) (*RtcRoom, error)
	Delete(appkey, roomId string) error
	UpdateAcceptedTime(appkey, roomId string, acceptedTime int64) error
}

type RtcRoomMember struct {
	ID          int64
	RoomId      string
	RoomType    pbobjs.RtcRoomType
	OwnerId     string
	MemberId    string
	DeviceId    string
	RtcState    pbobjs.RtcState
	InviterId   string
	CallTime    int64
	ConnectTime int64
	HangupTime  int64
	AppKey      string

	LatestPingTime int64
}

type IRtcRoomMemberStorage interface {
	Upsert(item RtcRoomMember) error
	Insert(item RtcRoomMember) (int64, error)
	Find(appkey, roomId, memberId string) (*RtcRoomMember, error)
	RefreshPingTime(appkey, roomId, memberId string) error
	UpdateState(appkey, roomId, memberId string, state pbobjs.RtcState, deviceId string) error
	Delete(appkey, roomId, memberId string) error
	DeleteByRoomId(appkey, roomId string) error
	DelteByRoomIdBaseTime(appkey, roomId string, baseTime int64) error
	QueryMembers(appkey, roomId string, startId, limit int64) ([]*RtcRoomMember, error)
	QueryRoomsByMember(appkey, memberId string, limit int64) ([]*RtcRoomMember, error)
}
