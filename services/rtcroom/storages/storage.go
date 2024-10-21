package storages

import (
	"im-server/services/rtcroom/storages/dbs"
	"im-server/services/rtcroom/storages/models"
)

func NewRtcRoomStorage() models.IRtcRoomStorage {
	return &dbs.RtcRoomDao{}
}

func NewRtcRoomMemberStorage() models.IRtcRoomMemberStorage {
	return &dbs.RtcRoomMemberDao{}
}
