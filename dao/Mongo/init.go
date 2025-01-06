package Mongo

import (
	"Fire/pkg/util/log"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoDB(MongoDBAddr, MongoDBPort string) {
	// 设置mongoDB客户端连接信息
	clientOptions := options.Client().ApplyURI("mongodb://" + MongoDBAddr + ":" + MongoDBPort)
	var err error
	MongoClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.LogrusObj.Info(err)
	}
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.LogrusObj.Info(err)
	}
	//log.LogrusObj.Info("Mongo Connect")
}

func GetClient() *mongo.Client {
	return MongoClient
}
