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
	"time"
)

func applyUserConverTagDeltas(ctx context.Context, userConverTags *UserConverTags, tags []models.UserConverTag) (changed []models.UserConverTag, code errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	maxUserTagLimit := 100
	if appinfo, ok := commonservices.GetAppInfo(appkey); ok {
		maxUserTagLimit = appinfo.MaxUserConverTags
	}
	changed, backup, code := userConverTags.AddTagsWithBackup(tags, maxUserTagLimit)
	if code != errs.IMErrorCode_SUCCESS {
		return nil, code
	}
	if len(changed) == 0 {
		return nil, errs.IMErrorCode_SUCCESS
	}
	if err := storages.NewUserConverTagStorage().BatchUpsert(changed); err != nil {
		logs.WithContext(ctx).Errorf("upsert user conver tags fail. err:%v", err)
		userConverTags.RollbackTags(backup)
		return nil, errs.IMErrorCode_CONVER_ADDTAGFAIL
	}
	return changed, errs.IMErrorCode_SUCCESS
}

func CreateUserConverTags(ctx context.Context, req *pbobjs.UserConverTags) (errs.IMErrorCode, *pbobjs.UserConverTags) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	userConverTags := GetUserConverTags(appkey, userId)
	curr := time.Now().UnixMilli()
	tags := make([]models.UserConverTag, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tags = append(tags, models.UserConverTag{
			AppKey:      appkey,
			UserId:      userId,
			Tag:         tag.Tag,
			TagName:     tag.TagName,
			TagOrder:    int(tag.TagOrder),
			CreatedTime: curr,
		})
	}
	changed, code := applyUserConverTagDeltas(ctx, userConverTags, tags)
	if code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	resp := &pbobjs.UserConverTags{
		Tags: []*pbobjs.ConverTag{},
	}
	if len(changed) > 0 {
		cmdmsg := &CmdMsg_CreateConverTags{
			Tags: []*ConverTag{},
		}
		for _, tag := range changed {
			cmdmsg.Tags = append(cmdmsg.Tags, &ConverTag{
				Tag:         tag.Tag,
				TagName:     tag.TagName,
				TagOrder:    int32(tag.TagOrder),
				TagType:     int32(pbobjs.ConverTagType_UserConverTag),
				CreatedTime: tag.CreatedTime,
				IsAdd:       tag.IsAdd,
			})
			resp.Tags = append(resp.Tags, &pbobjs.ConverTag{
				Tag:   tag.Tag,
				IsAdd: tag.IsAdd,
			})
		}
		// ntf other device
		flag := msgdefines.SetCmdMsg(0)
		bs, _ := json.Marshal(cmdmsg)
		commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
			MsgType:    msgdefines.CmdMsgType_CreateConverTags,
			MsgContent: bs,
			Flags:      flag,
		})
	}
	return errs.IMErrorCode_SUCCESS, resp
}

// createUserConverTagIfMissing persists a new user tag and updates cache, without CreateConverTags multi-device sync.
func createUserConverTagIfMissing(ctx context.Context, appkey, userId string, userConverTags *UserConverTags, req *pbobjs.TagConvers) errs.IMErrorCode {
	_, code := applyUserConverTagDeltas(ctx, userConverTags, []models.UserConverTag{
		{
			AppKey:      appkey,
			UserId:      userId,
			Tag:         req.Tag,
			TagName:     req.TagName,
			CreatedTime: time.Now().UnixMilli(),
		},
	})
	return code
}

func TagAddConvers(ctx context.Context, req *pbobjs.TagConvers) errs.IMErrorCode {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	cmdmsg := &CmdMsg_TagConvers{
		Tag: req.Tag,
	}
	userConverTags := GetUserConverTags(appkey, userId)
	if !userConverTags.ContainsTag(req.Tag) {
		if code := createUserConverTagIfMissing(ctx, appkey, userId, userConverTags, req); code != errs.IMErrorCode_SUCCESS {
			return code
		}
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
		err := relStorage.BatchCreate(rels)
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
			// userConvers.PersistConver(conver.TargetId, conver.SubChannel, conver.ChannelType)
			c := userConvers.QryConver(conver.TargetId, conver.SubChannel, conver.ChannelType)
			userConvers.PersistConverV2(c)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

type CmdMsg_CreateConverTags struct {
	Tags []*ConverTag `json:"tags"`
}

var CmdMsgType_TagAddConvers string = msgdefines.CmdMsgType_TagAddConvers
var CmdMsgType_TagDelConvers string = msgdefines.CmdMsgType_TagDelConvers

type CmdMsg_TagConvers struct {
	Tag     string          `json:"tag"`
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
			// userConvers.PersistConver(conver.TargetId, conver.SubChannel, conver.ChannelType)
			c := userConvers.QryConver(conver.TargetId, conver.SubChannel, conver.ChannelType)
			userConvers.PersistConverV2(c)
		}
	}
	return errs.IMErrorCode_SUCCESS
}

func QryUserConverTags(ctx context.Context) (*pbobjs.UserConverTags, errs.IMErrorCode) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)

	userConverTags := GetUserConverTags(appkey, userId)
	tags := userConverTags.QryTags()
	ret := &pbobjs.UserConverTags{
		Tags: []*pbobjs.ConverTag{},
	}
	for _, tag := range tags {
		ret.Tags = append(ret.Tags, &pbobjs.ConverTag{
			Tag:      tag.Tag,
			TagName:  tag.TagName,
			TagType:  pbobjs.ConverTagType_UserConverTag,
			TagOrder: int32(tag.TagOrder),
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
	userConverTags := GetUserConverTags(appkey, userId)
	for _, tag := range req.Tags {
		if userConverTags.ContainsTag(tag.Tag) {
			if userConverTags.DelTag(tag.Tag) {
				storage.Delete(appkey, userId, tag.Tag)
				err := relStorage.DeleteByTag(appkey, userId, tag.Tag)
				if err != nil {
					logs.WithContext(ctx).Error(err.Error())
				}
				cmdmsg.Tags = append(cmdmsg.Tags, &ConverTag{
					Tag: tag.Tag,
				})
			}
		}
	}
	// ntf other device
	if len(cmdmsg.Tags) > 0 {
		flag := msgdefines.SetCmdMsg(0)
		bs, _ := json.Marshal(cmdmsg)
		commonservices.AsyncPrivateMsg(ctx, userId, userId, &pbobjs.UpMsg{
			MsgType:    CmdMsgType_DelConverTags,
			MsgContent: bs,
			Flags:      flag,
		})
	}
	return errs.IMErrorCode_SUCCESS
}

var CmdMsgType_DelConverTags string = msgdefines.CmdMsgType_DelConverTags

type CmdMsg_DelUserConverTags struct {
	Tags []*ConverTag `json:"tags"`
}
type ConverTag struct {
	Tag         string `json:"tag"`
	TagName     string `json:"tag_name,omitempty"`
	TagOrder    int32  `json:"tag_order,omitempty"`
	TagType     int32  `json:"tag_type,omitempty"`
	CreatedTime int64  `json:"created_time,omitempty"`
	IsAdd       bool   `json:"is_add,omitempty"`
}
