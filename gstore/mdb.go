package gstore

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

// 初始化mongodb
func InitMdb(minPoolSize, maxPoolSize uint64, dSN string) *mongo.Client {
	return connMdb(minPoolSize, maxPoolSize, dSN)
}

func connMdb(minPoolSize, maxPoolSize uint64, dsn string) *mongo.Client {
	var err error
	MdbOptions := options.Client().
		ApplyURI(dsn).
		SetMaxPoolSize(minPoolSize).
		SetMinPoolSize(maxPoolSize)
	Mdb, err := mongo.NewClient(MdbOptions)
	if err != nil {
		log.Printf("[gstore_mongodb] mongo fail, err=%s", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Mdb.Connect(ctx)
	if err != nil {
		log.Printf("[gstore_mongodb] mongo fail, err=%s", err)
		panic(err)
	}
	log.Printf("[gstore_mongodb] mongo success")

	return Mdb
}

func MdbSurvive(mdb *mongo.Client) error {
	return mdb.Ping(context.Background(), readpref.Primary())
}
