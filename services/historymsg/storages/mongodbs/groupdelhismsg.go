package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/services/historymsg/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupDelHisMsgDao struct {
	UserId     string `bson:"user_id"`
	TargetId   string `bson:"target_id"`
	SubChannel string `bson:"sub_channel"`
	MsgId      string `bson:"msg_id"`
	MsgTime    int64  `bson:"msg_time"`
	MsgSeq     int64  `bson:"msg_seq"`
	AppKey     string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *GroupDelHisMsgDao) TableName() string {
	return "g_delhismsgs"
}

func (msg *GroupDelHisMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection((&GroupHisMsgDao{}).TableName())
}

func (msg *GroupDelHisMsgDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"sub_channel": 1},
				},
				{
					Keys: bson.M{"msg_time": -1},
				},
			})
		}
	}
}

func (msg *GroupDelHisMsgDao) Create(item models.GroupDelHisMsg) error {
	converId := commonservices.GetConversationId(item.UserId, item.TargetId, pbobjs.ChannelType_Group)
	filter := bson.M{"app_key": item.AppKey, "conver_id": converId, "sub_channel": item.SubChannel, "msg_id": item.MsgId}
	update := bson.M{"$addToSet": bson.M{"del_user_ids": bson.M{"$each": []string{item.UserId}}}}

	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg *GroupDelHisMsgDao) BatchCreate(items []models.GroupDelHisMsg) error {
	if len(items) <= 0 {
		return errors.New("no data need to insert")
	}
	msgIds := []string{}
	appkey := items[0].AppKey
	userId := items[0].UserId
	converId := commonservices.GetConversationId(items[0].UserId, items[0].TargetId, pbobjs.ChannelType_Group)
	subChannel := items[0].SubChannel

	for _, item := range items {
		msgIds = append(msgIds, item.MsgId)
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "sub_channel": subChannel, "msg_id": bson.M{"$in": msgIds}}
	update := bson.M{
		"$addToSet": bson.M{"del_user_ids": bson.M{"$each": []string{userId}}},
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func (msg *GroupDelHisMsgDao) QryDelHisMsgs(appkey, userId, targetId, subChannel string, startTime int64, count int32, isPositive bool) ([]*models.GroupDelHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.GroupDelHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	converId := commonservices.GetConversationId(userId, targetId, pbobjs.ChannelType_Group)
	filter := bson.M{"app_key": appkey, "conver_id": converId, "sub_channel": subChannel, "del_user_ids": bson.M{"$in": []string{userId}}}
	dbSort := -1
	if isPositive {
		dbSort = 1
		filter["send_time"] = bson.M{
			"$gt": startTime,
		}
	} else {
		filter["send_time"] = bson.M{
			"$lt": startTime,
		}
	}
	projection := bson.D{{"msg_id", 1}, {"sender_id", 1}, {"receiver_id", 1}, {"send_time", 1}, {"msg_seq_no", 1}, {"app_key", 1}}
	cur, err := collection.Find(context.TODO(), filter, options.Find().SetProjection(projection), options.Find().SetSort(bson.D{{"msg_time", dbSort}}), options.Find().SetLimit(int64(count)))
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var item GroupHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.GroupDelHisMsg{
				UserId:   userId,
				TargetId: item.ReceiverId,
				MsgId:    item.MsgId,
				MsgTime:  item.SendTime,
				MsgSeq:   item.MsgSeqNo,
				AppKey:   item.AppKey,
			})
		}
	}
	return retItems, nil
}

func (msg *GroupDelHisMsgDao) QryDelHisMsgsByMsgIds(appkey, userId, targetId, subChannel string, msgIds []string) ([]*models.GroupDelHisMsg, error) {
	return nil, nil
}
