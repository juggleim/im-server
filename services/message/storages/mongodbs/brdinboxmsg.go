package mongodbs

import (
	"context"
	"im-server/commons/configures"
	"im-server/commons/mongocommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/message/storages/models"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BrdInboxMsgDao struct {
	SenderId    string `bson:"sender_id"`
	SendTime    int64  `bson:"send_time"`
	MsgId       string `bson:"msg_id"`
	ChannelType int    `bson:"channel_type"`
	MsgBody     []byte `bson:"msg_body"`
	AppKey      string `bson:"app_key"`

	AddTime time.Time `bson:"add_time"`
}

func (msg *BrdInboxMsgDao) TableName() string {
	return "brdinboxmsgs"
}

func (msg *BrdInboxMsgDao) getCollection() *mongo.Collection {
	return mongocommons.GetCollection(msg.TableName())
}

func (msg *BrdInboxMsgDao) IndexCreator() func(colName string) {
	return func(colName string) {
		collection := mongocommons.GetCollection(colName)
		if collection != nil {
			collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
				{
					Keys: bson.M{"app_key": 1},
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
				Options: options.Index().SetExpireAfterSeconds(int32(configures.CmdMsgExpired / 1000)),
			})
		}
	}
}

func (msg *BrdInboxMsgDao) SaveMsg(item models.BrdInboxMsgMsg) error {
	collection := msg.getCollection()
	if collection != nil {
		_, err := collection.InsertOne(context.TODO(), BrdInboxMsgDao{
			SenderId:    item.SenderId,
			SendTime:    item.SendTime,
			MsgId:       item.MsgId,
			ChannelType: int(item.ChannelType),
			MsgBody:     item.MsgBody,
			AppKey:      item.AppKey,
			AddTime:     time.Now(),
		}, options.InsertOne().SetBypassDocumentValidation(true))
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

func (msg *BrdInboxMsgDao) QryMsgsBaseTime(appkey string, start int64, count int) ([]*models.BrdInboxMsgMsg, error) {
	retItems := []*models.BrdInboxMsgMsg{}
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey, "send_time": bson.M{"$gt": start}}

		cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"send_time", 1}}), options.Find().SetLimit(int64(count)))
		defer func() {
			if cur != nil {
				cur.Close(context.TODO())
			}
		}()
		if err == nil {
			for cur.Next(context.TODO()) {
				var item BrdInboxMsgDao
				err = cur.Decode(&item)
				if err == nil {
					retItems = append(retItems, &models.BrdInboxMsgMsg{
						SenderId:    item.SenderId,
						SendTime:    item.SendTime,
						MsgId:       item.MsgId,
						ChannelType: pbobjs.ChannelType(item.ChannelType),
						MsgBody:     item.MsgBody,
						AppKey:      item.AppKey,
					})
				}
			}
		}
	}
	return retItems, nil
}

func (msg *BrdInboxMsgDao) QryLatestMsg(appkey string, count int) ([]*models.BrdInboxMsgMsg, error) {
	retItems := []*models.BrdInboxMsgMsg{}
	collection := msg.getCollection()
	if collection != nil {
		filter := bson.M{"app_key": appkey}
		cur, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"send_time", -1}}), options.Find().SetLimit(int64(count)))
		defer func() {
			if cur != nil {
				cur.Close(context.TODO())
			}
		}()
		if err == nil {
			for cur.Next(context.TODO()) {
				var item BrdInboxMsgDao
				err = cur.Decode(&item)
				if err == nil {
					retItems = append(retItems, &models.BrdInboxMsgMsg{
						SenderId:    item.SenderId,
						SendTime:    item.SendTime,
						MsgId:       item.MsgId,
						ChannelType: pbobjs.ChannelType(item.ChannelType),
						MsgBody:     item.MsgBody,
						AppKey:      item.AppKey,
					})
				}
			}
		}
	}
	if len(retItems) > 0 {
		sort.Slice(retItems, func(i, j int) bool {
			return retItems[i].SendTime < retItems[j].SendTime
		})
	}
	return retItems, nil
}

func (msg *BrdInboxMsgDao) DelMsgsBaseTime(appkey string, start int64) error {
	return nil
}
