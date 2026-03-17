package dbs

import (
	"context"
	"log"
	"testing"
	"time"

	"im-server/services/pushmanager/storages/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserTagsDao_GetUserWithTags(t *testing.T) {
	db, err := initMysql("root:@tcp(127.0.0.1:3306)/im_tmp?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		t.Error(err)
		return
	}
	storage := NewUserTagsDao(db)
	ctx := context.Background()
	appKey := "testappKey"
	var userIDs []string
	storage.AddUserTags(ctx, appKey, "testuser1", "tag1", "tag2", "tag3")
	storage.AddUserTags(ctx, appKey, "testuser2", "tag4", "tag5", "tag3")

	userIDs, _ = storage.GetUserWithTags(ctx, appKey, models.Condition{
		TagsAnd: []string{"tag1", "tag2"},
	}, 1, 10)
	t.Logf("tagsAnd:%v", userIDs)

	userIDs, _ = storage.GetUserWithTags(ctx, appKey, models.Condition{
		TagsOr: []string{"tag4", "tag2"},
	}, 1, 10)
	t.Logf("tagsOr:%v", userIDs)

	storage.ClearUserTag(ctx, appKey, "testuser1")
	storage.ClearUserTag(ctx, appKey, "testuser2")
}

func initMysql(sqlDsn string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(sqlDsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(time.Second * 9) // mysql连接默认10s断开
	return
}
