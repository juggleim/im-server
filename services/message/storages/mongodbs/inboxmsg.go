package mongodbs

import (
	"context"
	"im-server/commons/configures"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InboxMsgDao struct {
	// ID          primitive.ObjectID `bson:"_id"`
	UserId      string `bson:"user_id"`
	SendTime    int64  `bson:"send_time"`
	MsgId       string `bson:"msg_id"`
	ChannelType int    `bson:"channel_type"`
	MsgBody     []byte `bson:"msg_body"`
	AppKey      string `bson:"app_key"`
	TargetId    string `bson:"target_id"`
	MsgType     string `bson:"msg_type"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *InboxMsgDao) TableName() string {
	return "inboxmsgs"
}

func (msg *InboxMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *InboxMsgDao) IndexCreator() func(colName string) {
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
					Keys: bson.M{"send_time": 1},
				},
				{
					Keys: bson.M{"msg_id": 1},
				},
			})
			collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    bson.M{"add_time": 1},
				Options: options.Index().SetExpireAfterSeconds(int32(configures.MsgExpired / 1000)),
			})
		}
	}
}

func (msg *InboxMsgDao) SaveMsg(item models.Msg) error {
	add := InboxMsgDao{
		UserId:      item.UserId,
		SendTime:    item.SendTime,
		MsgId:       item.MsgId,
		ChannelType: int(item.ChannelType),
		MsgBody:     item.MsgBody,
		AppKey:      item.AppKey,
		TargetId:    item.TargetId,
		MsgType:     item.MsgType,
		AddTime:     time.Now(),
	}
	collection := msg.getCollection()
	if collection != nil {
		_, err := collection.InsertOne(context.TODO(), add)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {
		return &mongocommons.MongoError{
			Msg: "Failed insert",
		}
	}
}

func (msg *InboxMsgDao) UpsertMsg(item models.Msg) error {
	return msg.SaveMsg(item)
}

func (msg *InboxMsgDao) QryMsgsBaseTime(appkey, userId string, start int64, count int) ([]*models.Msg, error) {
	inboxMsgs := []*models.Msg{}
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "user_id": userId, "send_time": bson.M{"$gt": start}}

		cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"send_time", 1}}), options.Find().SetLimit(int64(count)))
		defer func() {
			if cur != nil {
				cur.Close(context.TODO())
			}
		}()
		if err == nil {
			for cur.Next(context.TODO()) {
				var item InboxMsgDao
				err = cur.Decode(&item)
				if err == nil {
					inboxMsgs = append(inboxMsgs, &models.Msg{
						UserId:      item.UserId,
						SendTime:    item.SendTime,
						MsgId:       item.MsgId,
						ChannelType: pbobjs.ChannelType(item.ChannelType),
						MsgBody:     item.MsgBody,
						AppKey:      item.AppKey,
						TargetId:    item.TargetId,
						MsgType:     item.MsgType,
					})
				}
			}
		}
	}
	return inboxMsgs, nil
}

func (msg *InboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return nil
}
