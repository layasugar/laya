package gstore

import (
	"context"
	"github.com/layatips/laya/gconf"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

// 初始化mongodb
func InitMdb(cf *gconf.MdbConf) *mongo.Client {
	return connMdb(cf.MinPoolSize, cf.MaxPoolSize, cf.DSN)
}

func connMdb(minPoolSize, maxPoolSize uint64, dsn string) *mongo.Client {
	var err error
	MdbOptions := options.Client().
		ApplyURI(dsn).
		SetMaxPoolSize(minPoolSize).
		SetMinPoolSize(maxPoolSize)
	Mdb, err := mongo.NewClient(MdbOptions)
	if err != nil {
		log.Printf("[store_mongodb] connMdb open,err=%s\n", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Mdb.Connect(ctx)
	if err != nil {
		log.Printf("[store_mongodb] mongo connMdb error,err=%s\n", err)
		panic(err)
	}
	log.Printf("[store_mongodb] mongo connMdb success")

	return Mdb
}

func MdbSurvive(mdb *mongo.Client) error {
	return mdb.Ping(context.Background(), readpref.Primary())
}
