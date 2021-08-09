package gstore

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

const (
	defaultPoolMaxIdle     = 10                                 // 连接池空闲连接数量
	defaultPoolMaxOpen     = 100                                // 连接池最大连接数量
	defaultConnMaxLifeTime = time.Second * time.Duration(21600) // MySQL默认长连接时间为8个小时,所以我们设置连接可重用的时间为6小时最为合理
	defaultConnMaxIdleTime = time.Second * time.Duration(7200)  // 设置连接2个小时没有用到就断开连接(内存要求较高可降低该值)
)

type PoolCfg struct {
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifeTime int
	ConnMaxIdleTime int
}

// init db
func InitDB(dsn string, poolCfg *PoolCfg, logger logger.Interface) *gorm.DB {
	var err error

	Db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.Printf("[gstore_db] mysql open fail, err=%s", err)
		panic(err)
	}

	d, err := Db.DB()
	if err != nil {
		log.Printf("[gstore_db] mysql db fail, err=%s", err)
		panic(err)
	}

	setPoolMaxOpen(d, poolCfg.MaxOpenConn)
	setMaxIdleConn(d, poolCfg.MaxIdleConn)
	setConnMaxLifetime(d, poolCfg.ConnMaxLifeTime)
	setConnMaxIdleTime(d, poolCfg.ConnMaxIdleTime)

	err = d.Ping()
	if err != nil {
		log.Printf("[gstore_db] mysql ping fail, err:%s", err.Error())
		panic(err)
	}
	log.Printf("[gstore_db] mysql success")
	return Db
}

// mysql survive
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

func setPoolMaxOpen(d *sql.DB, n int) {
	if n == 0 {
		d.SetMaxOpenConns(defaultPoolMaxOpen)
	} else {
		d.SetMaxOpenConns(n)
	}
	return
}

func setMaxIdleConn(d *sql.DB, n int) {
	if n == 0 {
		d.SetMaxIdleConns(defaultPoolMaxIdle)
	} else {
		d.SetMaxIdleConns(n)
	}
	return
}

func setConnMaxLifetime(d *sql.DB, n int) {
	if n == 0 {
		d.SetConnMaxLifetime(defaultConnMaxLifeTime)
	} else {
		d.SetConnMaxLifetime(time.Second * time.Duration(n))
	}
	return
}

func setConnMaxIdleTime(d *sql.DB, n int) {
	if n == 0 {
		d.SetConnMaxIdleTime(defaultConnMaxIdleTime)
	} else {
		d.SetConnMaxIdleTime(time.Second * time.Duration(n))
	}
	return
}
