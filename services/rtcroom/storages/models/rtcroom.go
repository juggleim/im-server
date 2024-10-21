package models

import "im-server/commons/pbdefines/pbobjs"

type RtcRoom struct {
	RoomId   string
	RoomType pbobjs.RtcRoomType
	OwnerId  string
	AppKey   string
}

type IRtcRoomStorage interface {
	Create(item RtcRoom) error
	FindById(appkey, roomId string) (*RtcRoom, error)
	Delete(appkey, chatId string) error
}

type RtcRoomMember struct {
	ID           int64
	RoomId       string
	MemberId     string
	DeviceId     string
	RtcState     pbobjs.RtcState
	InviterId    string
	CameraEnable int32
	MicEnable    int32
	CallTime     int64
	ConnectTime  int64
	HangupTime   int64
	AppKey       string

	LatestPingTime int64
}

type IRtcRoomMemberStorage interface {
	Upsert(item RtcRoomMember) error
	Delete(appkey, roomId, memberId string) error
	DeleteByRoomId(appkey, roomId string) error
	QueryMembers(appkey, roomId string, startId, limit int64) ([]*RtcRoomMember, error)
	QueryRoomsByMember(appkey, memberId string, limit int64) ([]*RtcRoomMember, error)
}
