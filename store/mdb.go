package store

import (
	"context"
	"github.com/layatips/laya/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mdb *mongo.Client

// 初始化mongodb
func InitMdb() {
	c := gconf.GetMdbConf()
	if c.Open {
		connMdb(c.MinPoolSize, c.MaxPoolSize, c.DSN)
	}
}

func connMdb(minPoolSize, maxPoolSize uint64, dsn string) {
	var err error

	MdbOptions := options.Client().
		ApplyURI(dsn).
		SetMaxPoolSize(minPoolSize).
		SetMinPoolSize(maxPoolSize)
	Mdb, err = mongo.NewClient(MdbOptions)
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
}
