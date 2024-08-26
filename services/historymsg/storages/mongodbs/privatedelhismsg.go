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

type PrivateDelHisMsgDao struct {
	UserId   string `bson:"user_id"`
	TargetId string `bson:"target_id"`
	MsgId    string `bson:"msg_id"`
	MsgTime  int64  `bson:"msg_time"`
	MsgSeq   int64  `bson:"msg_seq"`
	AppKey   string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *PrivateDelHisMsgDao) TableName() string {
	return "p_delhismsgs"
}

func (msg *PrivateDelHisMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *PrivateDelHisMsgDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"msg_time": -1},
				},
			})
		}
	}
}

func (msg *PrivateDelHisMsgDao) Create(item models.PrivateDelHisMsg) error {
	add := PrivateDelHisMsgDao{
		UserId:   item.UserId,
		TargetId: item.TargetId,
		MsgId:    item.MsgId,
		MsgTime:  item.MsgTime,
		MsgSeq:   item.MsgSeq,
		AppKey:   item.AppKey,

		AddTime: time.Now(),
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}
func (msg *PrivateDelHisMsgDao) BatchCreate(items []models.PrivateDelHisMsg) error {
	if len(items) <= 0 {
		return errors.New("no data need to insert")
	}
	adds := []interface{}{}
	for _, item := range items {
		adds = append(adds, PrivateDelHisMsgDao{
			UserId:   item.UserId,
			TargetId: item.TargetId,
			MsgId:    item.MsgId,
			MsgTime:  item.MsgTime,
			MsgSeq:   item.MsgSeq,
			AppKey:   item.AppKey,

			AddTime: time.Now(),
		})
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertMany(context.TODO(), adds)
	return err
}
func (msg *PrivateDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId string, startTime int64, count int32, isPositive bool) ([]*models.PrivateDelHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.PrivateDelHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "user_id": userId, "target_id": targetId}
	dbSort := -1
	if isPositive {
		dbSort = 1
		filter["msg_time"] = bson.M{
			"$gt": startTime,
		}
	} else {
		filter["msg_time"] = bson.M{
			"$lt": startTime,
		}
	}
	cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"msg_time", dbSort}}), options.Find().SetLimit(int64(count)))
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var item PrivateDelHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.PrivateDelHisMsg{
				UserId:   item.UserId,
				TargetId: item.TargetId,
				MsgId:    item.MsgId,
				MsgTime:  item.MsgTime,
				MsgSeq:   item.MsgSeq,
				AppKey:   item.AppKey,
			})
		}
	}
	return retItems, nil
}
