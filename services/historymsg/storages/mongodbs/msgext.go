package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MsgExtDao struct {
	AppKey string `bson:"app_key"`
	MsgId  string `bson:"msg_id"`
	Key    string `bson:"key"`
	Value  string `bson:"value"`
}

func (ext *MsgExtDao) TableName() string {
	return "msgexts"
}

func (msg *MsgExtDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *MsgExtDao) IndexCreator() func(colName string) {
	return func(colName string) {

		collection := mongocommons.GetCollection(colName)
		if collection != nil {
			collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
				{
					Keys: bson.M{"app_key": 1},
				},
				{
					Keys: bson.M{"msg_id": 1},
				},
				{
					Keys: bson.M{"key": 1},
				},
			})
		}
	}
}

func (msg *MsgExtDao) Upsert(item models.MsgExt) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": item.AppKey, "msg_id": item.MsgId, "key": item.Key}
	update := bson.M{
		"$set": bson.M{
			"app_key": item.AppKey,
			"msg_id":  item.MsgId,
			"key":     item.Key,
			"value":   item.Value,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update, &options.UpdateOptions{Upsert: tools.BoolPtr(true)})
	return err
}
