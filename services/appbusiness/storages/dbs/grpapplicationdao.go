package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/services/appbusiness/storages/models"
	"time"
)

type GrpApplicationDao struct {
	ID          int64  `gorm:"primary_key"`
	GroupId     string `gorm:"group_id"`
	ApplyType   int    `gorm:"apply_type"`
	SponsorId   string `gorm:"sponsor_id"`
	RecipientId string `gorm:"recipient_id"`
	InviterId   string `gorm:"inviter_id"`
	OperatorId  string `gorm:"operator_id"`
	ApplyTime   int64  `gorm:"apply_time"`
	Status      int    `gorm:"status"`
	AppKey      string `gorm:"app_key"`
}

func (apply GrpApplicationDao) TableName() string {
	return "grpapplications"
}

func (apply GrpApplicationDao) InviteUpsert(item models.GrpApplication) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,apply_type,group_id,sponsor_id,recipient_id,inviter_id,operator_id,apply_time,status)VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE apply_time=VALUES(apply_time),status=VALUES(status),inviter_id=VALUES(inviter_id),operator_id=VALUES(operator_id)", apply.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, models.GrpApplicationType_Invite, item.GroupId, "", item.RecipientId, item.InviterId, item.OperatorId, item.ApplyTime, item.Status).Error
}

func (apply GrpApplicationDao) ApplyUpsert(item models.GrpApplication) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,apply_type,group_id,sponsor_id,recipient_id,inviter_id,operator_id,apply_time,status)VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE apply_time=VALUES(apply_time),status=VALUES(status),inviter_id=VALUES(inviter_id),operator_id=VALUES(operator_id)", apply.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, models.GrpApplicationType_Apply, item.GroupId, item.SponsorId, "", item.InviterId, item.OperatorId, item.ApplyTime, item.Status).Error
}

func (apply GrpApplicationDao) QueryMyGrpApplications(appkey, sponsorId string, startTime, count int64, isPositive bool) ([]*models.GrpApplication, error) {
	var items []*GrpApplicationDao
	params := []interface{}{}
	condition := "app_key=? and apply_type=? and sponsor_id=?"
	params = append(params, appkey)
	params = append(params, models.GrpApplicationType_Apply)
	params = append(params, sponsorId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.GrpApplication{}
	for _, app := range items {
		ret = append(ret, &models.GrpApplication{
			GroupId:     app.GroupId,
			ApplyType:   models.GrpApplicationType(app.ApplyType),
			SponsorId:   app.SponsorId,
			RecipientId: app.RecipientId,
			InviterId:   app.InviterId,
			OperatorId:  app.OperatorId,
			Status:      models.GrpApplicationStatus(app.Status),
			ApplyTime:   app.ApplyTime,
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}

func (apply GrpApplicationDao) QueryMyPendingGrpInvitations(appkey, recipientId string, startTime, count int64, isPositive bool) ([]*models.GrpApplication, error) {
	var items []*GrpApplicationDao
	params := []interface{}{}
	condition := "app_key=? and apply_type=? and recipient_id=?"
	params = append(params, appkey)
	params = append(params, models.GrpApplicationType_Invite)
	params = append(params, recipientId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.GrpApplication{}
	for _, app := range items {
		ret = append(ret, &models.GrpApplication{
			GroupId:     app.GroupId,
			ApplyType:   models.GrpApplicationType(app.ApplyType),
			SponsorId:   app.SponsorId,
			RecipientId: app.RecipientId,
			InviterId:   app.InviterId,
			OperatorId:  app.OperatorId,
			Status:      models.GrpApplicationStatus(app.Status),
			ApplyTime:   app.ApplyTime,
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}

func (apply GrpApplicationDao) QueryGrpInvitations(appkey, groupId string, startTime, count int64, isPositive bool) ([]*models.GrpApplication, error) {
	var items []*GrpApplicationDao
	params := []interface{}{}
	condition := "app_key=? and apply_type=? and group_id=?"
	params = append(params, appkey)
	params = append(params, models.GrpApplicationType_Invite)
	params = append(params, groupId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.GrpApplication{}
	for _, app := range items {
		ret = append(ret, &models.GrpApplication{
			GroupId:     app.GroupId,
			ApplyType:   models.GrpApplicationType(app.ApplyType),
			SponsorId:   app.SponsorId,
			RecipientId: app.RecipientId,
			InviterId:   app.InviterId,
			OperatorId:  app.OperatorId,
			Status:      models.GrpApplicationStatus(app.Status),
			ApplyTime:   app.ApplyTime,
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}

func (apply GrpApplicationDao) QueryGrpPendingApplications(appkey, groupId string, startTime, count int64, isPositive bool) ([]*models.GrpApplication, error) {
	var items []*GrpApplicationDao
	params := []interface{}{}
	condition := "app_key=? and apply_type=? and group_id=?"
	params = append(params, appkey)
	params = append(params, models.GrpApplicationType_Apply)
	params = append(params, groupId)
	orderStr := "apply_time desc"
	if isPositive {
		orderStr = "apply_time asc"
		condition = condition + " and apply_time>?"
	} else {
		if startTime <= 0 {
			startTime = time.Now().UnixMilli()
		}
		condition = condition + " and apply_time<?"
	}
	params = append(params, startTime)
	err := dbcommons.GetDb().Where(condition, params...).Order(orderStr).Limit(count).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.GrpApplication{}
	for _, app := range items {
		ret = append(ret, &models.GrpApplication{
			GroupId:     app.GroupId,
			ApplyType:   models.GrpApplicationType(app.ApplyType),
			SponsorId:   app.SponsorId,
			RecipientId: app.RecipientId,
			InviterId:   app.InviterId,
			OperatorId:  app.OperatorId,
			Status:      models.GrpApplicationStatus(app.Status),
			ApplyTime:   app.ApplyTime,
			AppKey:      app.AppKey,
		})
	}
	return ret, nil
}
