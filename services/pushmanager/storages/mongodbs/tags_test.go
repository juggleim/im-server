package mongodbs

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im-server/services/pushmanager/storages/models"
	"testing"
	"time"
)

func TestUserTagsDao_GetUserTags(t *testing.T) {
	db, err := initMongoDB("mongodb://127.0.0.1:27017")
	if err != nil {
		t.Error(err)
		return
	}
	storage := NewUserTagsDao(db.Database("im_db"))
	ctx := context.Background()
	appKey := "testappKey"
	var userIDs []string
	storage.AddUserTags(ctx, appKey, "testuser1", "tag1", "tag2", "tag3")
	storage.AddUserTags(ctx, appKey, "testuser2", "tag4", "tag5", "tag3", "tag1")

	userIDs, err = storage.GetUserWithTags(ctx, appKey, models.Condition{
		TagsAnd: []string{"tag1", "tag3"},
	}, 1, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tagsAnd:%v", userIDs)

	userIDs, err = storage.GetUserWithTags(ctx, appKey, models.Condition{
		TagsOr: []string{"tag4", "tag2"},
	}, 1, 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tagsOr:%v", userIDs)

	storage.DeleteUserTags(ctx, appKey, "testuser1", "tag1")
	t.Log(storage.GetUserTags(ctx, appKey, "testuser1"))

	storage.ClearUserTag(ctx, appKey, "testuser1")
	storage.ClearUserTag(ctx, appKey, "testuser2")
}

func initMongoDB(dsn string) (db *mongo.Client, err error) {
	clientOptions := options.Client().ApplyURI(dsn).SetConnectTimeout(5 * time.Second).SetMaxPoolSize(32)
	//clientOptions.Monitor = otelmongo.NewMonitor()

	// 连接到MongoDB
	db, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}
	// 检查连接
	err = db.Ping(context.TODO(), nil)
	if err != nil {
		return
	}
	return

}
