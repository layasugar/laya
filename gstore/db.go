package gstore

import (
	"github.com/layatips/laya/gconf"
	"github.com/layatips/laya/genv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

// init db
func InitDB(cf *gconf.DBConf) *gorm.DB {
	return mysqlConn(cf.MaxIdleConn, cf.MaxOpenConn, cf.ConnMaxLifetime, cf.Dsn, genv.RunMode())
}

// init mysql pool
func mysqlConn(maxIdleConn, maxOpenConn, connMaxLifetime int, dsn, runMode string) *gorm.DB {
	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[store_db] open connDB,err=%s\n", err)
		panic(err)
	}

	if runMode == "debug" {
		DB = DB.Debug()
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("[store_db] get Mdb,err=%s\n", err)
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
	return DB
}

// mysql 存活检测
func DBSurvive(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}
