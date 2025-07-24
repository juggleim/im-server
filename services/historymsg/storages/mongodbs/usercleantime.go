package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HisMsgUserCleanTimeDao struct {
	UserId      string `bson:"user_id"`
	TargetId    string `bson:"target_id"`
	ChannelType int    `bson:"channel_type"`
	SubChannel  string `bson:"sub_channel"`
	CleanTime   int64  `bson:"clean_time"`
	AppKey      string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *HisMsgUserCleanTimeDao) TableName() string {
	return "usercleantimes"
}

func (msg *HisMsgUserCleanTimeDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *HisMsgUserCleanTimeDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"target_id": 1},
				},
				{
					Keys: bson.M{"channel_type": 1},
				},
				{
					Keys: bson.M{"sub_channel": 1},
				},
			})
		}
	}
}

func (msg *HisMsgUserCleanTimeDao) UpsertCleanTime(item models.HisMsgUserCleanTime) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": item.AppKey, "user_id": item.UserId, "target_id": item.TargetId, "channel_type": item.ChannelType, "sub_channel": item.SubChannel}
	update := bson.M{
		"$set": bson.M{
			"app_key":      item.AppKey,
			"user_id":      item.UserId,
			"target_id":    item.TargetId,
			"channel_type": item.ChannelType,
			"sub_channel":  item.SubChannel,
			"clean_time":   item.CleanTime,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, &options.UpdateOptions{Upsert: tools.BoolPtr(true)})
	return err
}

func (msg *HisMsgUserCleanTimeDao) FindOne(appkey, userId, targetId, subChannel string, channelType pbobjs.ChannelType) (*models.HisMsgUserCleanTime, error) {
	collection := msg.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "user_id": userId, "target_id": targetId, "channel_type": channelType, "sub_channel": subChannel}
	result := collection.FindOne(context.TODO(), filter)
	var item HisMsgUserCleanTimeDao
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}
	return &models.HisMsgUserCleanTime{
		UserId:      item.UserId,
		TargetId:    item.TargetId,
		ChannelType: pbobjs.ChannelType(item.ChannelType),
		SubChannel:  item.SubChannel,
		CleanTime:   item.CleanTime,
		AppKey:      item.AppKey,
	}, nil
}
