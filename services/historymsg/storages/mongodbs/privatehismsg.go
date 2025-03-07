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

type PrivateHisMsgDao struct {
	// ID          primitive.ObjectID `bson:"_id"`
	ConverId    string    `bson:"conver_id"`
	SenderId    string    `bson:"sender_id"`
	ReceiverId  string    `bson:"receiver_id"`
	MsgType     string    `bson:"msg_type"`
	ChannelType int       `bson:"channel_type"`
	SendTime    int64     `bson:"send_time"`
	MsgId       string    `bson:"msg_id"`
	MsgSeqNo    int64     `bson:"msg_seq_no"`
	MsgBody     []byte    `bson:"msg_body"`
	IsRead      int       `bson:"is_read"`
	AppKey      string    `bson:"app_key"`
	IsDelete    int       `bson:"is_delete"`
	IsExt       int       `bson:"is_ext"`
	IsExset     int       `bson:"is_exset"`
	MsgExt      []byte    `bson:"msg_ext"`
	MsgExset    []byte    `bson:"msg_exset"`
	DelUserIds  []string  `bson:"del_user_ids"`
	AddTime     time.Time `bson:"add_time"`
}

func (msg *PrivateHisMsgDao) TableName() string {
	return "p_hismsgs"
}

func (msg *PrivateHisMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *PrivateHisMsgDao) IndexCreator() func(colName string) {
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
				{
					Keys: bson.M{"msg_index": 1},
				},
				{
					Keys: bson.M{"sender_id": 1},
				},
				{
					Keys: bson.M{"del_user_ids": 1},
				},
			})
		}
	}
}

func (msg *PrivateHisMsgDao) SavePrivateHisMsg(item models.PrivateHisMsg) error {
	add := PrivateHisMsgDao{
		ConverId:    item.ConverId,
		SenderId:    item.SenderId,
		ReceiverId:  item.ReceiverId,
		ChannelType: int(item.ChannelType),
		MsgType:     item.MsgType,
		MsgId:       item.MsgId,
		SendTime:    item.SendTime,
		MsgSeqNo:    item.MsgSeqNo,
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		IsExt:       item.IsExt,
		IsExset:     item.IsExset,
		MsgExt:      item.MsgExt,
		MsgExset:    item.MsgExset,
		IsDelete:    item.IsDelete,
		IsRead:      item.IsRead,
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (msg *PrivateHisMsgDao) FindById(appkey, converId, msgId string) (*models.PrivateHisMsg, error) {
	collection := msg.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	result := collection.FindOne(context.TODO(), filter)
	var item PrivateHisMsgDao
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}
	return dbMsg2PrivateMsg(&item), nil
}

func (msg *PrivateHisMsgDao) UpdateMsgBody(appkey, converId, msgId, msgType string, msgBody []byte) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"msg_type": msgType,
			"msg_body": msgBody,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg *PrivateHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "conver_id": converId}
		result := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(bson.M{"msg_seq_no": 1}), options.FindOne().SetSort(bson.D{{"send_time", -1}}))
		var item GroupHisMsgDao
		err := result.Decode(&item)
		if err == nil {
			return item.MsgSeqNo
		}
	}
	return 0
}

func (msg *PrivateHisMsgDao) QryLatestMsg(appkey, converId string) (*models.PrivateHisMsg, error) {
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "conver_id": converId}
		result := collection.FindOne(context.TODO(), filter, options.FindOne().SetSort(bson.D{{"send_time", -1}}))
		var item PrivateHisMsgDao
		err := result.Decode(&item)
		if err == nil {
			return dbMsg2PrivateMsg(&item), nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("no mongo client")
}

func (msg *PrivateHisMsgDao) QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.PrivateHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.PrivateHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "del_user_ids": bson.M{"$nin": []string{userId}}}
	dbSort := -1
	start := startTime
	if isPositiveOrder {
		dbSort = 1
		if start < cleanTime {
			start = cleanTime
		}
		filter["send_time"] = bson.M{
			"$gt": start,
		}
	} else {
		if start <= 0 {
			start = time.Now().UnixMilli()
		}
		filter["send_time"] = bson.M{
			"$lt": start,
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
		var item PrivateHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2PrivateMsg(&item))
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *PrivateHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*models.PrivateHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.PrivateHisMsg{}
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId}

	dbSort := -1
	start := startTime
	if isPositiveOrder {
		dbSort = 1

		if start < cleanTime {
			start = cleanTime
		}
		filter["send_time"] = bson.M{
			"$gt": start,
		}
	} else {
		if start <= 0 {
			start = time.Now().UnixMilli()
		}
		if cleanTime > 0 {
			filter["send_time"] = bson.M{
				"$lt": start,
				"$gt": cleanTime,
			}
		} else {
			filter["send_time"] = bson.M{
				"$lt": start,
			}
		}
	}

	if len(msgTypes) > 0 {
		filter["msg_type"] = bson.M{"$in": msgTypes}
	}
	if len(excludeMsgIds) > 0 {
		filter["msg_id"] = bson.M{"$nin": excludeMsgIds}
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
		var item PrivateHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2PrivateMsg(&item))
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *PrivateHisMsgDao) MarkReadByMsgIds(appkey, converId string, msgIds []string) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": bson.M{"$in": msgIds}}
	update := bson.M{
		"$set": bson.M{
			"is_read": 1,
		},
	}
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func (msg *PrivateHisMsgDao) MarkReadByScope(appkey, converId string, start, end int64) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_index": bson.M{
		"$gte": start,
		"$lte": end,
	}}
	update := bson.M{
		"$set": bson.M{
			"is_read": 1,
		},
	}
	_, err := collection.UpdateMany(context.TODO(), filter, update)
	return err
}

func (msg *PrivateHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.PrivateHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.PrivateHisMsg{}
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
		var item PrivateHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2PrivateMsg(&item))
		}
	}
	return retItems, nil
}

func (msg *PrivateHisMsgDao) FindByConvers(appkey string, convers []models.ConverItem) ([]*models.PrivateHisMsg, error) {
	retItems := []*models.PrivateHisMsg{}
	length := len(convers)
	if length < 0 {
		return retItems, nil
	}
	collection := msg.getCollection()
	if collection == nil {
		return retItems, errors.New("no mongo client")
	}
	or := []bson.M{}
	for _, conver := range convers {
		or = append(or, bson.M{"conver_id": conver.ConverId, "msg_id": conver.MsgId})
	}
	filter := bson.M{
		"app_key": appkey,
		"$or":     or,
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
		var item PrivateHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2PrivateMsg(&item))
		}
	}
	return retItems, nil
}

func (msg *PrivateHisMsgDao) DelMsgs(appkey, converId string, msgIds []string) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": bson.M{"$in": msgIds}}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (msg *PrivateHisMsgDao) UpdateMsgExtState(appkey, converId, msgId string, isExt int) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"is_ext": isExt,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg *PrivateHisMsgDao) UpdateMsgExt(appkey, converId, msgId string, ext []byte) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"msg_ext": ext,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg PrivateHisMsgDao) UpdateMsgExsetState(appkey, converId, msgId string, isExset int) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"is_exset": isExset,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg PrivateHisMsgDao) UpdateMsgExset(appkey, converId, msgId string, ext []byte) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"msg_exset": ext,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

// TODO need batch delete
func (msg *PrivateHisMsgDao) DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "send_time": bson.M{"$lt": cleanTime}, "sender_id": senderId}
	_, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func dbMsg2PrivateMsg(dbMsg *PrivateHisMsgDao) *models.PrivateHisMsg {
	return &models.PrivateHisMsg{
		HisMsg: models.HisMsg{
			ConverId:    dbMsg.ConverId,
			SenderId:    dbMsg.SenderId,
			ReceiverId:  dbMsg.ReceiverId,
			ChannelType: pbobjs.ChannelType(dbMsg.ChannelType),
			MsgType:     dbMsg.MsgType,
			MsgId:       dbMsg.MsgId,
			SendTime:    dbMsg.SendTime,
			MsgSeqNo:    dbMsg.MsgSeqNo,
			MsgBody:     dbMsg.MsgBody,
			AppKey:      dbMsg.AppKey,
			IsExt:       dbMsg.IsExt,
			IsExset:     dbMsg.IsExset,
			MsgExt:      dbMsg.MsgExt,
			MsgExset:    dbMsg.MsgExset,
			IsDelete:    dbMsg.IsDelete,
		},
		IsRead: dbMsg.IsRead,
	}
}
