package mongodbs

import (
	"context"
	"im-server/commons/mongocommons"
	"im-server/services/pushmanager/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ models.IUserTagStorage = (*UserTagsDao)(nil)

type UserTags struct {
	AppKey  string    `bson:"app_key"`
	UserId  string    `bson:"user_id"`
	Tags    []string  `bson:"tags"`
	AddTime time.Time `bson:"add_time"`
}

type UserTagsDao struct {
	db *mongo.Database
}

func NewUserTagsDao(db *mongo.Database) *UserTagsDao {
	return &UserTagsDao{
		db: db,
	}
}

func (u *UserTagsDao) TableName() string {
	return "user_tags"
}

func (u *UserTagsDao) getCollection() *mongo.Collection {
	return u.db.Collection(u.TableName())
}

func (u *UserTagsDao) IndexCreator() func(colName string) {
	return func(colName string) {
		collection := mongocommons.GetCollection(colName)
		if collection != nil {
			collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
				{
					Keys: bson.M{"app_key": 1},
				},
				{
					Keys: bson.M{"user_id": 1},
				},
				{
					Keys: bson.M{"tags": 1},
				},
			})
		}
	}
}

func (u *UserTagsDao) AddUserTags(ctx context.Context, appKey string, userID string, tags ...string) error {
	filter := bson.M{"app_key": appKey, "user_id": userID}
	update := bson.M{"$addToSet": bson.M{"tags": bson.M{"$each": tags}}}
	_, err := u.getCollection().UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (u *UserTagsDao) DeleteUserTags(ctx context.Context, appKey string, userID string, tags ...string) error {
	filter := bson.M{"app_key": appKey, "user_id": userID}
	update := bson.M{"$pull": bson.M{"tags": bson.M{"$in": tags}}}
	_, err := u.getCollection().UpdateOne(ctx, filter, update)
	return err
}

func (u *UserTagsDao) ClearUserTag(ctx context.Context, appKey string, userID string) error {
	filter := bson.M{"app_key": appKey, "user_id": userID}
	_, err := u.getCollection().DeleteMany(ctx, filter)
	return err
}

func (u *UserTagsDao) GetUserTags(ctx context.Context, appKey string, userID string) ([]string, error) {
	filter := bson.M{"app_key": appKey, "user_id": userID}
	cursor, err := u.getCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []string
	for cursor.Next(ctx) {
		var userTag UserTags
		if err := cursor.Decode(&userTag); err != nil {
			return nil, err
		}
		tags = append(tags, userTag.Tags...)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (u *UserTagsDao) GetUserWithTags(ctx context.Context,
	appKey string, condition models.Condition, page int, perPage int) (userIDs []string, err error) {
	collection := u.getCollection()
	var filter bson.M

	if len(condition.TagsAnd) > 0 {
		filter = bson.M{"app_key": appKey, "tags": bson.M{"$all": condition.TagsAnd}}
	} else {
		filter = bson.M{"app_key": appKey, "tags": bson.M{"$in": condition.TagsOr}}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64((page - 1) * perPage))
	findOptions.SetLimit(int64(perPage))

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result struct {
			UserID string `bson:"user_id"`
		}
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, result.UserID)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
