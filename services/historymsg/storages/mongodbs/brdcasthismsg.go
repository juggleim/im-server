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

type BrdCastHisMsgDao struct {
	// ID          primitive.ObjectID `bson:"_id"`
	ConverId    string `bson:"conver_id"`
	SenderId    string `bson:"sender_id"`
	ChannelType int    `bson:"channel_type"`
	MsgType     string `bson:"msg_type"`
	MsgId       string `bson:"msg_id"`
	SendTime    int64  `bson:"send_time"`
	MsgSeqNo    int64  `bson:"msg_seq_no"`
	MsgBody     []byte `bson:"msg_body"`
	AppKey      string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *BrdCastHisMsgDao) TableName() string {
	return "bc_hismsgs"
}

func (msg *BrdCastHisMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *BrdCastHisMsgDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"send_time": -1},
				},
				{
					Keys: bson.M{"msg_type": 1},
				},
				{
					Keys: bson.M{"msg_id": 1},
				},
			})
		}
	}
}

func (msg *BrdCastHisMsgDao) SaveBrdCastHisMsg(item models.BrdCastHisMsg) error {
	add := BrdCastHisMsgDao{
		ConverId:    item.ConverId,
		SenderId:    item.SenderId,
		ChannelType: int(item.ChannelType),
		MsgType:     item.MsgType,
		MsgId:       item.MsgId,
		SendTime:    item.SendTime,
		MsgSeqNo:    item.MsgSeqNo,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,

		AddTime: time.Now(),
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (msg *BrdCastHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "conver_id": converId}
		result := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"msg_seq_no": 1}), options.FindOne().SetSort(bson.D{{"send_time", -1}}))
		var item BrdCastHisMsgDao
		err := result.Decode(&item)
		if err == nil {
			return item.MsgSeqNo
		}
	}
	return 0
}

func (msg *BrdCastHisMsgDao) QryLatestMsg(appkey, converId string) (*models.BrdCastHisMsg, error) {
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "conver_id": converId}
		result := collection.FindOne(context.TODO(), filter, options.FindOne().SetSort(bson.D{{"send_time", -1}}))
		var item BrdCastHisMsgDao
		err := result.Decode(&item)
		if err == nil {
			return dbMsg2BrdMsg(&item), nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("no mongo client")
}

func (msg *BrdCastHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.BrdCastHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.BrdCastHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId}
	dbSort := -1
	if isPositiveOrder {
		dbSort = 1
		begin := startTime
		if begin < cleanTime {
			begin = cleanTime
		}
		filter["send_time"] = bson.M{
			"$gt": begin,
		}
	} else {
		filter["send_time"] = bson.M{
			"$lt": startTime,
			"$gt": cleanTime,
		}
	}
	if len(msgTypes) > 0 {
		filter["msg_type"] = bson.M{"$in": msgTypes}
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
		var item BrdCastHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2BrdMsg(&item))
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *BrdCastHisMsgDao) FindById(appkey, converId, msgId string) (*models.BrdCastHisMsg, error) {
	collection := msg.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	result := collection.FindOne(context.TODO(), filter)
	var item BrdCastHisMsgDao
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}
	return dbMsg2BrdMsg(&item), nil
}

func (msg *BrdCastHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.BrdCastHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.BrdCastHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId,
		"send_time": bson.M{
			"$gt": cleanTime,
		},
		"msg_id": bson.M{
			"$in": msgIds,
		},
	}

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
		var item BrdCastHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2BrdMsg(&item))
		}
	}
	return retItems, nil
}

func dbMsg2BrdMsg(dbMsg *BrdCastHisMsgDao) *models.BrdCastHisMsg {
	return &models.BrdCastHisMsg{
		ConverId:    dbMsg.ConverId,
		SenderId:    dbMsg.SenderId,
		ChannelType: pbobjs.ChannelType(dbMsg.ChannelType),
		MsgType:     dbMsg.MsgType,
		MsgId:       dbMsg.MsgId,
		SendTime:    dbMsg.SendTime,
		MsgSeqNo:    dbMsg.MsgSeqNo,
		MsgBody:     dbMsg.MsgBody,
		AppKey:      dbMsg.AppKey,
	}
}
