package mongodb

import (
	"context"
	"github.com/LaYa-op/laya/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var Mdb *mongo.Client

// 初始化mongodb
func Init() {
	path := ""
	Configs := config.ListFiles(path)

	log.Printf("[store_mongodb] DB_INIT with %d cluster\n", len(Configs))
	var config Config
	for _, name := range Configs {
		err := config.ReadFile(name, &config)
		if err != nil {
			log.Printf("[store_mongodb] parse db config %s failed,err= %s\n", name, err)
			continue
		}
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
