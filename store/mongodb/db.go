package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mdb *mongo.Client

type Config struct {
	Open        bool   `toml:"open"`
	DSN         string `toml:"dsn"`
	maxIdleConn uint64 `toml:"maxIdleConn"`
	maxOpenConn uint64 `toml:"maxOpenConn"`
}

// 初始化mongodb
func Init(config *Config) {
	if config.Open {
		conn(config)
	}
}

func conn(conf *Config) {
	var err error

	MdbOptions := options.Client().
		ApplyURI(conf.DSN).
		SetMaxPoolSize(conf.maxIdleConn).
		SetMinPoolSize(conf.maxIdleConn)
	Mdb, err = mongo.NewClient(MdbOptions)
	if err != nil {
		log.Printf("[store_mongodb] open conn,err=%s\n", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Mdb.Connect(ctx)

}
