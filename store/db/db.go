package db

import (
	"github.com/LaYa-op/laya/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// db is sql *db
var DB *gorm.DB

func Init(configPath ...string) {
	path := ""
	Configs := conf.ListFiles(path)
	if len(configPath) == 1 {
		path = configPath[0]
	}
	var config Config
	for _, name := range Configs {
		err := conf.ReadFile(name, &config)
		if err != nil {
			log.Printf("[store_db] parse db conf %s failed,err= %s\n", name, err)
			continue
		}
	}
	if config.Open {
		mysqlDB(&config)
	}
}

// init mysql pool
func mysqlDB(conf *Config) {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[store_db] open conn,err= %s\n", err)
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("[store_db] get DB,err= %s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(conf.ConnMaxLifetime))
}

func pgsqlDB(conf *Config) {}

func sqlLiteDB(conf *Config) {}

func sqlServerDB(conf *Config) {}
