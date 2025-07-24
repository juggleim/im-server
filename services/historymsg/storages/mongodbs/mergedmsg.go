package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MergedMsgDao struct {
	// ID          primitive.ObjectID `bson:"_id"`
	ParentMsgId string    `bson:"parent_msg_id"`
	FromId      string    `bson:"from_id"`
	TargetId    string    `bson:"target_id"`
	ChannelType int       `bson:"channel_type"`
	SubChannel  string    `bson:"sub_channel"`
	MsgId       string    `bson:"msg_id"`
	MsgTime     int64     `bson:"msg_time"`
	MsgBody     []byte    `bson:"msg_body"`
	AppKey      string    `bson:"app_key"`
	AddTime     time.Time `bson:"add_time"`
}

func (msg *MergedMsgDao) TableName() string {
	return "mergedmsgs"
}

func (msg *MergedMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *MergedMsgDao) IndexCreator() func(colName string) {
	return func(colName string) {
		collection := mongocommons.GetCollection(colName)
		if collection != nil {
			collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
				{
					Keys: bson.M{"app_key": 1},
				},
				{
					Keys: bson.M{"parent_msg_id": 1},
				},
				{
					Keys: bson.M{"send_time": 1},
				},
			})
		}
	}
}

func (msg *MergedMsgDao) SaveMergedMsg(item models.MergedMsg) error {
	add := MergedMsgDao{
		ParentMsgId: item.ParentMsgId,
		FromId:      item.FromId,
		TargetId:    item.TargetId,
		ChannelType: int(item.ChannelType),
		SubChannel:  item.SubChannel,
		MsgId:       item.MsgId,
		MsgTime:     item.MsgTime,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		AddTime:     time.Now(),
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (msg *MergedMsgDao) BatchSaveMergedMsgs(items []models.MergedMsg) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	adds := []interface{}{}
	for _, item := range items {
		adds = append(adds, MergedMsgDao{
			ParentMsgId: item.ParentMsgId,
			FromId:      item.FromId,
			TargetId:    item.TargetId,
			ChannelType: int(item.ChannelType),
			SubChannel:  item.SubChannel,
			MsgId:       item.MsgId,
			MsgTime:     item.MsgTime,
			MsgBody:     item.MsgBody,
			AppKey:      item.AppKey,
			AddTime:     time.Now(),
		})
	}
	_, err := collection.InsertMany(context.TODO(), adds)
	return err
}

func (msg *MergedMsgDao) QryMergedMsgs(appkey, parentMsgId string, startTime int64, count int32, isPositiveOrder bool) ([]*models.MergedMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.MergedMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "parent_msg_id": parentMsgId}
	dbSort := -1
	if isPositiveOrder {
		dbSort = 1
		filter["send_time"] = bson.M{
			"$gt": startTime,
		}
	} else {
		filter["send_time"] = bson.M{
			"$lt": startTime,
		}
	}
	cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"send_time", dbSort}}), options.Find().SetLimit(int64(count)))
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var item MergedMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.MergedMsg{
				ParentMsgId: item.ParentMsgId,
				FromId:      item.FromId,
				TargetId:    item.TargetId,
				ChannelType: pbobjs.ChannelType(item.ChannelType),
				SubChannel:  item.SubChannel,
				MsgId:       item.MsgId,
				MsgTime:     item.MsgTime,
				MsgBody:     item.MsgBody,
				AppKey:      item.AppKey,
			})
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].MsgTime < retItems[j].MsgTime
		})
	}
	return retItems, nil
}
