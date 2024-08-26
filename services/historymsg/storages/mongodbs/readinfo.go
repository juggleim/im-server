package mongodbs

import (
	"context"
	"errors"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/historymsg/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReadInfoDao struct {
	AppKey      string    `bson:"app_key"`
	MsgId       string    `bson:"msg_id"`
	ChannelType int       `bson:"channel_type"`
	GroupId     string    `bson:"group_id"`
	MemberId    string    `bson:"member_id"`
	CreatedTime time.Time `bson:"created_time"`
}

func (info *ReadInfoDao) TableName() string {
	return "readinfos"
}

func (info *ReadInfoDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(info.TableName())
}

func (info *ReadInfoDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"channel_type": 1},
				},
				{
					Keys: bson.M{"group_id": 1},
				},
				{
					Keys: bson.M{"member_id": 1},
				},
			})
			collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys: bson.D{
					{"app_key", 1},
					{"msg_id", 1},
					{"channel_type", 1},
					{"group_id", 1},
					{"member_id", 1},
				},
				Options: options.Index().SetUnique(true),
			})
		}
	}
}

func (info *ReadInfoDao) Create(item models.ReadInfo) error {
	collection := info.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	var addTime time.Time
	if item.CreatedTime > 0 {
		addTime = time.UnixMilli(item.CreatedTime)
	} else {
		addTime = time.Now()
	}
	add := ReadInfoDao{
		AppKey:      item.AppKey,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		GroupId:     item.GroupId,
		MemberId:    item.MemberId,
		CreatedTime: addTime,
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (info *ReadInfoDao) BatchCreate(items []models.ReadInfo) error {
	if len(items) <= 0 {
		return errors.New("no data need to insert")
	}
	adds := []interface{}{}
	for _, item := range items {
		var addTime time.Time
		if item.CreatedTime > 0 {
			addTime = time.UnixMilli(item.CreatedTime)
		} else {
			addTime = time.Now()
		}
		adds = append(adds, ReadInfoDao{
			AppKey:      item.AppKey,
			MsgId:       item.MsgId,
			ChannelType: int(item.ChannelType),
			GroupId:     item.GroupId,
			MemberId:    item.MemberId,
			CreatedTime: addTime,
		})
	}
	collection := info.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertMany(context.TODO(), adds)
	return err
}

func (info *ReadInfoDao) QryReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string, startId, limit int64) ([]*models.ReadInfo, error) {
	collection := info.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "group_id": groupId, "channel_type": channelType, "msg_id": msgId}
	cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"created_time", 1}}), options.Find().SetLimit(limit))
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return nil, err
	}
	retItems := []*models.ReadInfo{}
	for cur.Next(context.TODO()) {
		var item ReadInfoDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, &models.ReadInfo{
				AppKey:      item.AppKey,
				MsgId:       item.MsgId,
				ChannelType: pbobjs.ChannelType(item.ChannelType),
				GroupId:     item.GroupId,
				MemberId:    item.MemberId,
				CreatedTime: item.CreatedTime.UnixMilli(),
			})
		}
	}
	return retItems, nil
}

func (info *ReadInfoDao) CountReadInfosByMsgId(appkey, groupId string, channelType pbobjs.ChannelType, msgId string) int32 {
	collection := info.getCollection()
	if collection == nil {
		return 0
	}
	filter := bson.M{"app_key": appkey, "group_id": groupId, "channel_type": channelType, "msg_id": msgId}
	count, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0
	}
	return int32(count)
}

func (info ReadInfoDao) CheckMsgsRead(appkey, groupId, memberId string, channelType pbobjs.ChannelType, msgIds []string) (map[string]bool, error) {
	ret := map[string]bool{}
	for _, msgId := range msgIds {
		ret[msgId] = false
	}
	collection := info.getCollection()
	if collection == nil {
		return ret, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "group_id": groupId, "channel_type": channelType, "member_id": memberId, "msg_id": bson.M{"$in": msgIds}}
	cur, err := collection.Find(context.TODO(), filter)
	defer func() {
		if cur != nil {
			cur.Close(context.TODO())
		}
	}()
	if err != nil {
		return ret, err
	}
	for cur.Next(context.TODO()) {
		var item ReadInfoDao
		err = cur.Decode(&item)
		if err == nil {
			if _, exist := ret[item.MsgId]; exist {
				ret[item.MsgId] = true
			}
		}
	}
	return ret, nil
}
