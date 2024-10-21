package services

import (
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/connectmanager/server/codec"
)

func PreProcessRtcCreate(msg *codec.QueryMsgBody) {
	if msg.Topic == "rtc_create" {
		var room pbobjs.RtcRoom
		err := tools.PbUnMarshal(msg.Data, &room)
		if err == nil {
			if room.RoomId == "" {
				room.RoomId = tools.GenerateUUIDShort22()
				msg.TargetId = room.RoomId
				bs, err := tools.PbMarshal(&room)
				if err == nil {
					msg.Data = bs
				}
			} else {
				msg.TargetId = room.RoomId
			}
		}
	}
}
