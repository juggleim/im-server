package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/tools"
	"im-server/services/historymsg/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MsgExtDao struct {
	AppKey      string    `bson:"app_key"`
	MsgId       string    `bson:"msg_id"`
	Key         string    `bson:"key"`
	Value       string    `bson:"value"`
	CreatedTime time.Time `bson:"created_time"`

	AddTime time.Time `bson:"add_time"`
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
			collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys: bson.D{
					{"app_key", 1},
					{"msg_id", 1},
					{"key", 1},
				},
				Options: options.Index().SetUnique(true),
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

func (msg *MsgExtDao) Delete(appkey, msgId, key string) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "msg_id": msgId, "key": key}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (msg *MsgExtDao) QryExtsByMsgIds(appkey string, msgIds []string) ([]*models.MsgExt, error) {
	collection := msg.getCollection()
	retItems := []*models.MsgExt{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "msg_id": bson.M{"$in": msgIds}}
	cur, err := collection.Find(context.TODO(), filter)
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var item MsgExtDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.MsgExt{
				AppKey:      appkey,
				MsgId:       item.MsgId,
				Key:         item.Key,
				Value:       item.Value,
				CreatedTime: item.CreatedTime.UnixMilli(),
			})
		}
	}
	return retItems, nil
}
