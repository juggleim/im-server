package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages"
	"im-server/services/historymsg/storages/models"
	"time"
)

func SetTopMsg(ctx context.Context, req *pbobjs.SetTopMsgReq) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewTopMsgStorage()
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	err := storage.Upsert(models.TopMsg{
		ConverId:    converId,
		ChannelType: req.ChannelType,
		MsgId:       req.MsgId,
		UserId:      userId,
		CreatedTime: time.Now(),
		AppKey:      appkey,
	})
	if err != nil {
		return errs.IMErrorCode_MSG_DEFAULT
	}
	// send cmd msg
	contentBs, _ := tools.JsonMarshal(&TopMsgCmd{
		MsgId: req.MsgId,
	})
	upMsg := &pbobjs.UpMsg{
		MsgType:    topMsgType,
		MsgContent: contentBs,
		Flags:      commonservices.SetCmdMsg(0),
	}
	if req.ChannelType == pbobjs.ChannelType_Private {
		commonservices.AsyncPrivateMsg(ctx, userId, req.TargetId, upMsg)
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		commonservices.AsyncGroupMsg(ctx, userId, req.TargetId, upMsg)
	}
	return errs.IMErrorCode_SUCCESS
}

var topMsgType string = "jg:topmsg"

type TopMsgCmd struct {
	MsgId string `json:"msg_id"`
}

func GetTopMsg(ctx context.Context, req *pbobjs.GetTopMsgReq) (errs.IMErrorCode, *pbobjs.TopMsg) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	converId := commonservices.GetConversationId(userId, req.TargetId, req.ChannelType)
	storage := storages.NewTopMsgStorage()
	msg, err := storage.FindTopMsg(appkey, converId, req.ChannelType)
	if err != nil {
		return errs.IMErrorCode_MSG_DEFAULT, nil
	}
	downMsg := &pbobjs.DownMsg{}
	if req.ChannelType == pbobjs.ChannelType_Private {
		hisStorage := storages.NewPrivateHisMsgStorage()
		hisMsg, err := hisStorage.FindById(appkey, converId, msg.MsgId)
		if err != nil || hisMsg == nil {
			return errs.IMErrorCode_MSG_DEFAULT, nil
		}
		msg := &pbobjs.DownMsg{}
		err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
		if err != nil {
			return errs.IMErrorCode_MSG_DEFAULT, nil
		}
		downMsg = msg
	} else if req.ChannelType == pbobjs.ChannelType_Group {
		hisStorage := storages.NewGroupHisMsgStorage()
		hisMsg, err := hisStorage.FindById(appkey, converId, msg.MsgId)
		if err != nil || hisMsg == nil {
			return errs.IMErrorCode_MSG_DEFAULT, nil
		}
		msg := &pbobjs.DownMsg{}
		err = tools.PbUnMarshal(hisMsg.MsgBody, msg)
		if err != nil {
			return errs.IMErrorCode_MSG_DEFAULT, nil
		}
		downMsg = msg
	} else {
		return errs.IMErrorCode_MSG_DEFAULT, nil
	}
	ret := &pbobjs.TopMsg{
		Operator:    commonservices.GetTargetDisplayUserInfo(ctx, msg.UserId),
		CreatedTime: msg.CreatedTime.UnixMilli(),
		Msg:         downMsg,
	}
	return errs.IMErrorCode_SUCCESS, ret
}
