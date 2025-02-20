package services

import (
	"context"
	"encoding/json"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/commonservices/logs"
	"im-server/services/commonservices/msgdefines"
	"im-server/services/conversation/storages"
	"im-server/services/conversation/storages/models"
)

var (
	SystemConverTag_All       string = "jg_all"
	SystemConverTag_Unread    string = "jg_unread"
	SystemConverTag_MentionMe string = "jg_mentionme"
	SystemConverTag_Private   string = "jg_private"
	SystemConverTag_Group     string = "jg_group"
)

var systemConverTags []*pbobjs.ConverTag

func init() {
	systemConverTags = append(systemConverTags, &pbobjs.ConverTag{
		Tag:     SystemConverTag_All,
		TagType: pbobjs.ConverTagType_SystemConverTag,
	})
	systemConverTags = append(systemConverTags, &pbobjs.ConverTag{
		Tag:     SystemConverTag_Unread,
		TagType: pbobjs.ConverTagType_SystemConverTag,
	})
	systemConverTags = append(systemConverTags, &pbobjs.ConverTag{
		Tag:     SystemConverTag_MentionMe,
		TagType: pbobjs.ConverTagType_SystemConverTag,
	})
	systemConverTags = append(systemConverTags, &pbobjs.ConverTag{
		Tag:     SystemConverTag_Private,
		TagType: pbobjs.ConverTagType_SystemConverTag,
	})
	systemConverTags = append(systemConverTags, &pbobjs.ConverTag{
		Tag:     SystemConverTag_Group,
		TagType: pbobjs.ConverTagType_SystemConverTag,
	})
}

func TagAddConvers(ctx context.Context, req *pbobjs.TagConvers) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewUserConverTagStorage()
	cmdmsg := &CmdMsg_TagConvers{
		Tag:     req.Tag,
		TagName: req.TagName,
	}
	err := storage.Upsert(models.UserConverTag{
		AppKey:  appkey,
		UserId:  userId,
		Tag:     req.Tag,
		TagName: req.TagName,
	})
	if err != nil {
		logs.WithContext(ctx).Errorf("create user conver tag fail. err:%v", err)
		return errs.IMErrorCode_CONVER_ADDTAGFAIL
	}
	if len(req.Convers) > 0 {
		cmdmsg.Convers = []*SimpleConver{}
		relStorage := storages.NewConverTagRelStorage()
		rels := []models.ConverTagRel{}
		for _, conver := range req.Convers {
			rels = append(rels, models.ConverTagRel{
				UserId:      userId,
				Tag:         req.Tag,
				TargetId:    conver.TargetId,
				ChannelType: conver.ChannelType,
				AppKey:      appkey,
			})
			cmdmsg.Convers = append(cmdmsg.Convers, &SimpleConver{
				TargetId:    conver.TargetId,
				ChannelType: int32(conver.ChannelType),
			})
		}
		err = relStorage.BatchCreate(rels)
		if err != nil {
			logs.WithContext(ctx).Errorf("tag add convers fail. err:%v", err)
			return errs.IMErrorCode_CONVER_TAGADDCONVERFAIL
		}
	}
	// ntf other device
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(cmdmsg)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    CmdMsgType_TagAddConvers,
		MsgContent: bs,
		Flags:      flag,
	})
	//cache
	userConvers := getUserConvers(appkey, userId)
	affected := userConvers.TagAddConvers(req.Tag, req.Convers)
	if affected {
		for _, conver := range req.Convers {
			userConvers.PersistConver(conver.TargetId, conver.ChannelType)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

var CmdMsgType_TagAddConvers string = msgdefines.CmdMsgType_TagAddConvers
var CmdMsgType_TagDelConvers string = msgdefines.CmdMsgType_TagDelConvers

type CmdMsg_TagConvers struct {
	Tag     string          `json:"tag"`
	TagName string          `json:"tag_name,omitempty"`
	Convers []*SimpleConver `json:"convers,omitempty"`
}

type SimpleConver struct {
	TargetId    string `json:"target_id"`
	ChannelType int32  `json:"channel_type"`
}

func TagDelConvers(ctx context.Context, req *pbobjs.TagConvers) errs.IMErrorCode {
	if len(req.Convers) <= 0 || req.Tag == "" {
		return errs.IMErrorCode_SUCCESS
	}
	cmdmsg := &CmdMsg_TagConvers{
		Tag:     req.Tag,
		Convers: []*SimpleConver{},
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewConverTagRelStorage()
	convers := []models.TargetConver{}
	for _, conver := range req.Convers {
		convers = append(convers, models.TargetConver{
			TargetId:    conver.TargetId,
			ChannelType: conver.ChannelType,
		})
		cmdmsg.Convers = append(cmdmsg.Convers, &SimpleConver{
			TargetId:    conver.TargetId,
			ChannelType: int32(conver.ChannelType),
		})
	}
	err := storage.BatchDelete(appkey, userId, req.Tag, convers)
	if err != nil {
		logs.WithContext(ctx).Errorf("err:%v", err)
	}
	// ntf other device
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(cmdmsg)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    CmdMsgType_TagDelConvers,
		MsgContent: bs,
		Flags:      flag,
	})
	//cache
	userConvers := getUserConvers(appkey, userId)
	affected := userConvers.TagDelConvers(req.Tag, req.Convers)
	if affected {
		for _, conver := range req.Convers {
			userConvers.PersistConver(conver.TargetId, conver.ChannelType)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryUserConverTags(ctx context.Context) (*pbobjs.UserConverTags, errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewUserConverTagStorage()
	tags, err := storage.QryTags(appkey, userId)
	if err != nil {
		logs.WithContext(ctx).Errorf("err:%v", err)
	}
	ret := &pbobjs.UserConverTags{
		Tags: []*pbobjs.ConverTag{},
	}
	//add system conver tags
	ret.Tags = append(ret.Tags, systemConverTags...)
	for _, tag := range tags {
		ret.Tags = append(ret.Tags, &pbobjs.ConverTag{
			Tag:     tag.Tag,
			TagName: tag.TagName,
		})
	}
	return ret, errs.IMErrorCode_SUCCESS
}

func DelUserConverTags(ctx context.Context, req *pbobjs.UserConverTags) errs.IMErrorCode {
	if len(req.Tags) <= 0 {
		return errs.IMErrorCode_SUCCESS
	}
	cmdmsg := &CmdMsg_DelUserConverTags{
		Tags: []*ConverTag{},
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewUserConverTagStorage()
	relStorage := storages.NewConverTagRelStorage()
	for _, tag := range req.Tags {
		storage.Delete(appkey, userId, tag.Tag)
		relStorage.DeleteByTag(appkey, userId, tag.Tag)
		cmdmsg.Tags = append(cmdmsg.Tags, &ConverTag{
			Tag: tag.Tag,
		})
	}
	// ntf other device
	flag := msgdefines.SetCmdMsg(0)
	bs, _ := json.Marshal(cmdmsg)
	commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
		MsgType:    CmdMsgType_DelConverTags,
		MsgContent: bs,
		Flags:      flag,
	})
	return errs.IMErrorCode_SUCCESS
}

var CmdMsgType_DelConverTags string = msgdefines.CmdMsgType_DelConverTags

type CmdMsg_DelUserConverTags struct {
	Tags []*ConverTag `json:"tags"`
}
type ConverTag struct {
	Tag string `json:"tag"`
}
