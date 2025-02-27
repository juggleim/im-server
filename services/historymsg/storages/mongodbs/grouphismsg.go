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

type GroupHisMsgDao struct {
	// ID          primitive.ObjectID `bson:"_id"`
	ConverId    string   `bson:"conver_id"`
	SenderId    string   `bson:"sender_id"`
	ReceiverId  string   `bson:"receiver_id"`
	ChannelType int      `bson:"channel_type"`
	MsgType     string   `bson:"msg_type"`
	MsgId       string   `bson:"msg_id"`
	SendTime    int64    `bson:"send_time"`
	MsgSeqNo    int64    `bson:"msg_seq_no"`
	MsgBody     []byte   `bson:"msg_body"`
	AppKey      string   `bson:"app_key"`
	IsExt       int      `bson:"is_ext"`
	IsExset     int      `bson:"is_exset"`
	MsgExt      []byte   `bson:"msg_ext"`
	MsgExset    []byte   `bson:"msg_exset"`
	MemberCount int      `bson:"member_count"`
	ReadCount   int      `bson:"read_count"`
	IsDelete    int      `bson:"is_delete"`
	DelUserIds  []string `bson:"del_user_ids"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *GroupHisMsgDao) TableName() string {
	return "g_hismsgs"
}

func (msg *GroupHisMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *GroupHisMsgDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"sender_id": 1},
				},
				{
					Keys: bson.M{"del_user_ids": 1},
				},
			})
		}
	}
}

func (msg *GroupHisMsgDao) SaveGroupHisMsg(item models.GroupHisMsg) error {
	add := GroupHisMsgDao{
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
		MemberCount: item.MemberCount,
		ReadCount:   item.ReadCount,
		IsDelete:    item.IsDelete,

		AddTime: time.Now(),
	}
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	_, err := collection.InsertOne(context.TODO(), add)
	return err
}

func (msg *GroupHisMsgDao) QryLatestMsgSeqNo(appkey, converId string) int64 {
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

func (msg *GroupHisMsgDao) QryLatestMsg(appkey, converId string) (*models.GroupHisMsg, error) {
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "conver_id": converId}
		result := collection.FindOne(context.TODO(), filter, options.FindOne().SetSort(bson.D{{"send_time", -1}}))
		var item GroupHisMsgDao
		err := result.Decode(&item)
		if err == nil {
			return dbMsg2GrpMsg(&item), nil
		} else {
			return nil, err
		}
	}
	return nil, errors.New("no mongo client")
}

func (msg *GroupHisMsgDao) QryHisMsgsExcludeDel(appkey, converId, userId, targetId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string) ([]*models.GroupHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.GroupHisMsg{}
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
		var item GroupHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2GrpMsg(&item))
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *GroupHisMsgDao) QryHisMsgs(appkey, converId string, startTime int64, count int32, isPositiveOrder bool, cleanTime int64, msgTypes []string, excludeMsgIds []string) ([]*models.GroupHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.GroupHisMsg{}
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
		var item GroupHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2GrpMsg(&item))
		}
	}
	if !isPositiveOrder {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *GroupHisMsgDao) UpdateMsgBody(appkey, converId, msgId, msgType string, msgBody []byte) error {
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

func (msg *GroupHisMsgDao) UpdateReadCount(appkey, converId, msgId string, readCount int) error {
	collection := msg.getCollection()
	if collection == nil {
		return errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	update := bson.M{
		"$set": bson.M{
			"read_count": readCount,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (msg *GroupHisMsgDao) FindById(appkey, converId, msgId string) (*models.GroupHisMsg, error) {
	collection := msg.getCollection()
	if collection == nil {
		return nil, errors.New("no mongo client")
	}
	filter := bson.M{"app_key": appkey, "conver_id": converId, "msg_id": msgId}
	result := collection.FindOne(context.TODO(), filter)
	var item GroupHisMsgDao
	err := result.Decode(&item)
	if err != nil {
		return nil, err
	}
	return dbMsg2GrpMsg(&item), nil
}

func (msg *GroupHisMsgDao) FindByIds(appkey, converId string, msgIds []string, cleanTime int64) ([]*models.GroupHisMsg, error) {
	collection := msg.getCollection()
	retItems := []*models.GroupHisMsg{}
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
		var item GroupHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2GrpMsg(&item))
		}
	}
	return retItems, nil
}

func (msg *GroupHisMsgDao) FindByConvers(appkey string, convers []models.ConverItem) ([]*models.GroupHisMsg, error) {
	retItems := []*models.GroupHisMsg{}
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
		var item GroupHisMsgDao
		err = cur.Decode(&item)
		if err == nil {
			retItems = append(retItems, dbMsg2GrpMsg(&item))
		}
	}
	return retItems, nil
}

func (msg *GroupHisMsgDao) DelMsgs(appkey, converId string, msgIds []string) error {
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

func (msg *GroupHisMsgDao) UpdateMsgExtState(appkey, converId, msgId string, isExt int) error {
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

func (msg *GroupHisMsgDao) UpdateMsgExt(appkey, converId, msgId string, ext []byte) error {
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

func (msg *GroupHisMsgDao) UpdateMsgExsetState(appkey, converId, msgId string, isExset int) error {
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

func (msg *GroupHisMsgDao) UpdateMsgExset(appkey, converId, msgId string, ext []byte) error {
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
func (msg *GroupHisMsgDao) DelSomeoneMsgsBaseTime(appkey, converId string, cleanTime int64, senderId string) error {
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

func dbMsg2GrpMsg(dbMsg *GroupHisMsgDao) *models.GroupHisMsg {
	return &models.GroupHisMsg{
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
		MemberCount: dbMsg.MemberCount,
		ReadCount:   dbMsg.ReadCount,
	}
}
