package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/commonservices"
	"im-server/services/group/dbs"
	"time"
)

var locks *tools.SegmentatedLocks

func init() {
	locks = tools.NewSegmentatedLocks(64)
}

func QrySnapshot(ctx context.Context, req *pbobjs.QryGrpSnapshotReq) (errs.IMErrorCode, *pbobjs.GroupSnapshot) {
	if req.GroupId == "" {
		return errs.IMErrorCode_GROUP_NOSNAPSHOT, nil
	}
	appkey := bases.GetAppKeyFromCtx(ctx)
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appinfo != nil && appinfo.OpenGrpSnapshot {
		dao := dbs.GrpSnapshotDao{}
		snapshotDb, err := dao.FindNearlySnapshot(appkey, req.GroupId, req.NearlyTime)
		if err == nil && snapshotDb != nil {
			snapshot := &pbobjs.GroupSnapshot{}
			err = tools.PbUnMarshal(snapshotDb.Snapshot, snapshot)
			if err == nil {
				return errs.IMErrorCode_SUCCESS, snapshot
			}
		}
		go RegenerateGroupSnapshot(appkey, req.GroupId)
		return errs.IMErrorCode_GROUP_NOSNAPSHOT, nil
	} else {
		container, exist := GetGroupMembersFromCache(ctx, appkey, req.GroupId)
		snapshot := &pbobjs.GroupSnapshot{
			GroupId:   req.GroupId,
			MemberIds: []string{},
		}
		if exist && container != nil {
			memberMap := container.GetMemberMap()
			for memberId := range memberMap {
				snapshot.MemberIds = append(snapshot.MemberIds, memberId)
			}
		}
		return errs.IMErrorCode_SUCCESS, snapshot
	}
}
func GenerateGroupSnapshot(appkey, groupId string, memberIds []string) {
	appinfo, exist := commonservices.GetAppInfo(appkey)
	if exist && appinfo != nil && appinfo.OpenGrpSnapshot {
		innerGenerateGroupSnapshot(appkey, groupId, memberIds, time.Now().UnixMilli())
	}
}
func innerGenerateGroupSnapshot(appkey, groupId string, memberIds []string, createdTime int64) {
	dao := dbs.GrpSnapshotDao{}
	snapshot := &pbobjs.GroupSnapshot{
		GroupId:   groupId,
		MemberIds: memberIds,
	}
	bs, _ := tools.PbMarshal(snapshot)
	dao.Create(dbs.GrpSnapshotDao{
		AppKey:      appkey,
		GroupId:     groupId,
		CreatedTime: createdTime,
		Snapshot:    bs,
	})
}

func RegenerateGroupSnapshot(appkey, groupId string) {
	lock := locks.GetLocks(appkey, groupId)
	lock.Lock()
	defer lock.Unlock()

	//check
	snapshot := dbs.GrpSnapshotDao{}
	exist := snapshot.Exist(appkey, groupId)
	if exist {
		return
	}

	memberDao := dbs.GroupMemberDao{}
	var startId int64 = 0
	memberIds := []string{}
	var maxCreatedTime int64 = 0
	for i := 0; i < 3; i++ {
		members, err := memberDao.QueryMembers(appkey, groupId, startId, 1000)
		if err == nil && len(members) > 0 {
			for _, member := range members {
				memberIds = append(memberIds, member.MemberId)
				if member.ID > startId {
					startId = member.ID
				}
				createdTime := member.CreatedTime.UnixMilli()
				if createdTime > maxCreatedTime {
					maxCreatedTime = createdTime
				}
			}
		} else {
			break
		}
	}
	if len(memberIds) > 0 {
		innerGenerateGroupSnapshot(appkey, groupId, memberIds, maxCreatedTime)
	}
}
