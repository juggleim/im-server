package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/services/historymsg/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MsgExSetDao struct {
	AppKey      string    `bson:"app_key"`
	MsgId       string    `bson:"msg_id"`
	Key         string    `bson:"key"`
	Item        string    `bson:"item"`
	CreatedTime time.Time `bson:"created_time"`

	AddTime time.Time `bson:"add_time"`
}

func (ext *MsgExSetDao) TableName() string {
	return "msgexsets"
}

func (ext *MsgExSetDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(ext.TableName())
}

func (ext *MsgExSetDao) IndexCreator() func(colName string) {
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
				{
					Keys: bson.M{"item": 1},
				},
			})
			collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys: bson.D{
					{"app_key", 1},
					{"msg_id", 1},
					{"key", 1},
					{"item", 1},
				},
				Options: options.Index().SetUnique(true),
			})
		}
	}
}

func (ext *MsgExSetDao) Create(item models.MsgExSet) error {
	add := MsgExSetDao{
		AppKey:      item.AppKey,
		MsgId:       item.MsgId,
		Key:         item.Key,
		Item:        item.Item,
		CreatedTime: time.UnixMilli(item.CreatedTime),

		AddTime: time.Now(),
	}
	collection := ext.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (ext *MsgExSetDao) Delete(appkey, msgId, key, item string) error {
	collection := ext.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "msg_id": msgId, "key": key, "item": item}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (ext *MsgExSetDao) QryExtsByMsgIds(appkey string, msgIds []string) ([]*models.MsgExSet, error) {
	collection := ext.getCollection()
	retItems := []*models.MsgExSet{}
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
		var item MsgExSetDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.MsgExSet{
				AppKey:      appkey,
				MsgId:       item.MsgId,
				Key:         item.Key,
				Item:        item.Item,
				CreatedTime: item.CreatedTime.UnixMilli(),
			})
		}
	}
	return retItems, nil
}
