package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type Config struct {
	Open            bool   `toml:"open"`
	Driver          string `toml:"driver"`
	Dsn             string `toml:"dsn"`
	MaxIdleConn     int    `toml:"maxIdleConn"`
	MaxOpenConn     int    `toml:"maxOpenConn"`
	ConnMaxLifetime int    `toml:"connMaxLifetime"`
}

// db is sql *db
var DB *gorm.DB

func Init(config *Config) {
	if config.Open {
		mysqlConn(config)
	}
}

// init mysql pool
func mysqlConn(conf *Config) {
	var err error
	DB, err = gorm.Open(mysql.Open(conf.Dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[store_db] open conn,err=%s\n", err)
		panic(err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("[store_db] get DB,err=%s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(conf.ConnMaxLifetime))
	log.Printf("[store_db] mysql conn success")
}

func pgsqlConn(conf *Config) {}

func sqlLiteConn(conf *Config) {}

func sqlServerConn(conf *Config) {}
