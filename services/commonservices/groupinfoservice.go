package commonservices

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/caches"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

var (
	groupInfoCache *caches.LruCache
	groupInfoLocks *tools.SegmentatedLocks

	GroupField_MemberCount string = "member_count"
)

func getGrpKey(appkey, userId string) string {
	return strings.Join([]string{appkey, userId}, "_")
}
func init() {
	groupInfoCache = caches.NewLruCacheWithAddReadTimeout("grpinfo_cache", 10000, nil, 5*time.Second, 5*time.Second)
	groupInfoLocks = tools.NewSegmentatedLocks(256)
}

func GetGroupInfoFromCache(ctx context.Context, groupId string) *pbobjs.GroupInfo {
	appkey := bases.GetAppKeyFromCtx(ctx)
	key := getGrpKey(appkey, groupId)
	if val, exist := groupInfoCache.Get(key); exist {
		return val.(*pbobjs.GroupInfo)
	} else {
		l := groupInfoLocks.GetLocks(key)
		l.Lock()
		defer l.Unlock()

		if val, exist := groupInfoCache.Get(key); exist {
			return val.(*pbobjs.GroupInfo)
		} else {
			grpInfo := GetGroupInfoFromRpc(ctx, groupId)
			groupInfoCache.Add(key, grpInfo)
			return grpInfo
		}
	}
}

func GetGroupInfoFromRpc(ctx context.Context, groupId string) *pbobjs.GroupInfo {
	_, respObj, err := bases.SyncRpcCall(ctx, "qry_group_info", groupId, &pbobjs.GroupInfoReq{
		GroupId: groupId,
	}, func() proto.Message {
		return &pbobjs.GroupInfo{}
	})
	if err == nil && respObj != nil {
		return respObj.(*pbobjs.GroupInfo)
	}
	return &pbobjs.GroupInfo{
		GroupId: groupId,
	}
}
