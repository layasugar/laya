package gstore

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// init db
func InitDB(maxIdleConn, maxOpenConn, connMaxLifetime int, dsn string, logger logger.Interface) *gorm.DB {
	return mysqlConn(maxIdleConn, maxOpenConn, connMaxLifetime, dsn, logger)
}

// init mysql pool
func mysqlConn(maxIdleConn, maxOpenConn, connMaxLifetime int, dsn string, logger logger.Interface) *gorm.DB {
	var err error
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		log.Printf("[gstore_db] mysql fail, err=%s", err)
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("[gstore_db] mysql fail, err=%s", err)
		panic(err)
	}
	sqlDB.SetMaxIdleConns(maxIdleConn)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * time.Duration(connMaxLifetime))

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("[gstore_db] mysql fail, err:%s", err.Error())
		panic(err)
	}
	log.Printf("[gstore_db] mysql success")
	return DB
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
