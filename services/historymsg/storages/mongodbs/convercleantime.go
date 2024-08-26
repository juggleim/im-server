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

type HisMsgConverCleanTimeDao struct {
	ConverId    string `bson:"conver_id"`
	ChannelType int    `bson:"channel_type"`
	CleanTime   int64  `bson:"clean_time"`
	AppKey      string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *HisMsgConverCleanTimeDao) TableName() string {
	return "convercleantimes"
}

func (msg *HisMsgConverCleanTimeDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *HisMsgConverCleanTimeDao) IndexCreator() func(colName string) {
	return func(colName string) {
		collection := mongocommons.GetCollection(colName)
		if collection != nil {
			collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
				{
					Keys: bson.M{"app_key": 1},
				},
				{
					Keys: bson.M{"conver_id": 1},
				},
				{
					Keys: bson.M{"channel_type": 1},
				},
			})
		}
	}
}

func (msg *HisMsgConverCleanTimeDao) UpsertDestroyTime(item models.HisMsgConverCleanTime) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": item.AppKey, "conver_id": item.ConverId, "channel_type": item.ChannelType}
	update := bson.M{
		"$set": bson.M{
			"app_key":      item.AppKey,
			"conver_id":    item.ConverId,
			"channel_type": item.ChannelType,
			"clean_time":   item.CleanTime,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, &options.UpdateOptions{Upsert: tools.BoolPtr(true)})
	return err
}

func (msg *HisMsgConverCleanTimeDao) FindOne(appkey, converId string, channelType pbobjs.ChannelType) (*models.HisMsgConverCleanTime, error) {
	collection := msg.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "channel_type": channelType}
	result := collection.FindOne(context.TODO(), filter)
	var item HisMsgConverCleanTimeDao
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}
	return &models.HisMsgConverCleanTime{
		ConverId:    item.ConverId,
		ChannelType: pbobjs.ChannelType(item.ChannelType),
		CleanTime:   item.CleanTime,
		AppKey:      item.AppKey,
	}, nil
}
