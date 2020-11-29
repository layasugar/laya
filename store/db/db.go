package db

import (
	"github.com/BurntSushi/toml"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// db is sql *db
var DB *gorm.DB

var path = "./config/db/db.toml"

func Init() {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Printf("[store_db] parse db config %s failed,err= %s\n", path, err)
		return
	}
	if config.Open {
		mysqlConn(&config)
	}
}

// init mysql pool
func mysqlConn(conf *Config) {
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

func pgsqlConn(conf *Config) {}

func sqlLiteConn(conf *Config) {}

func sqlServerConn(conf *Config) {}
