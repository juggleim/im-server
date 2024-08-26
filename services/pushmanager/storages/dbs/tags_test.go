package dbs

import (
	"context"
	"github.com/jinzhu/gorm"
	"im-server/services/pushmanager/storages/models"
	"log"
	"testing"
	"time"
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
	db, err = gorm.Open("mysql", sqlDsn)

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
		return
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "" + defaultTableName
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(500)
	db.DB().SetConnMaxLifetime(time.Second * 9) // mysql连接默认10s断开
	return
}
