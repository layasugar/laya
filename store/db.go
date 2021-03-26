package store

import (
	"github.com/layatips/laya/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// db is sql *db
var DB *gorm.DB

func InitDB() {
	runMode := config.GetRunMode()
	c := config.GetDBConf()
	if c.Open {
		mysqlConn(c.MaxIdleConn, c.MaxOpenConn, c.ConnMaxLifetime, c.Dsn, runMode)
	}
}

// init mysql pool
func mysqlConn(maxIdleConn, maxOpenConn, connMaxLifetime int, dsn, runMode string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[store_db] open connDB,err=%s\n", err)
		panic(err)
	}

	if runMode == "debug" {
		DB = DB.Debug()
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("[store_db] get DB,err=%s\n", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(connMaxLifetime))

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("[store_db] mysql connDB err:%s", err.Error())
		panic(err)
	}
	log.Printf("[store_db] mysql connDB success")
}
