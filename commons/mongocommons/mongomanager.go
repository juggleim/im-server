package mongocommons

import (
	"context"
	"fmt"
	"im-server/commons/configures"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

var MongoDbName string

type MongoError struct {
	Msg string
}

func (e *MongoError) Error() string {
	return e.Msg
}

var indexCreatorMap map[string]func(collectionName string)

func InitMongodb() error {
	indexCreatorMap = make(map[string]func(collectionName string))
	mongoUrl := fmt.Sprintf("mongodb://%s", configures.Config.MongoDb.Address) // "mongodb://127.0.0.1:27017"
	MongoDbName = configures.Config.MongoDb.DbName                             //"msg_db"

	var err error
	clientOptions := options.Client().ApplyURI(mongoUrl).SetConnectTimeout(5 * time.Second).SetMaxPoolSize(32)
	//clientOptions.Monitor = otelmongo.NewMonitor()

	// 连接到MongoDB
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func GetMongoClient() *mongo.Client {
	if client == nil {
		InitMongodb()
	}
	return client
}

func GetCollection(collectionName string) *mongo.Collection {
	mongoClient := GetMongoClient()
	if mongoClient != nil {
		return mongoClient.Database(MongoDbName).Collection(collectionName)
	}
	return nil
}

func GetMongoDatabase() *mongo.Database {
	mongoClient := GetMongoClient()
	if mongoClient != nil {
		return mongoClient.Database(MongoDbName)
	}
	return nil
}

func RegIndexCreator(name string, creator func(colName string)) {
	indexCreatorMap[name] = creator
}

func Register(creator CollectionCreator) {
	RegIndexCreator(creator.TableName(), creator.IndexCreator())
}

type CollectionCreator interface {
	TableName() string
	IndexCreator() func(colName string)
}

func InitMongoCollections() error {
	mongoClient := GetMongoClient()
	existColNames, err := mongoClient.Database(MongoDbName).ListCollectionNames(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	existNameMap := make(map[string]bool)
	for _, colName := range existColNames {
		existNameMap[colName] = true
	}

	for colName, creator := range indexCreatorMap {
		if _, exist := existNameMap[colName]; !exist { //initial
			err := mongoClient.Database(MongoDbName).CreateCollection(context.Background(), colName)
			if err != nil {
				fmt.Println("create collection[", colName, "] failed")
				continue
			} else {
				fmt.Println("create collection[", colName, "] success")
			}
			if creator != nil {
				fmt.Println("create index for collection [", colName, "]")
				creator(colName)
			}
		}
	}
	return nil
}
