package mongodb

import (
	"context"
	"github.com/BurntSushi/toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mdb *mongo.Client

var path = "./config/mongo"

// 初始化mongodb
func Init() {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Printf("[store_mongodb] parse db config %s failed,err= %s\n", path, err)
		return
	}
	if config.Open {
		conn(&config)
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
		log.Printf("[store_mongodb] open conn,err= %s\n", err)
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = Mdb.Connect(ctx)

}
