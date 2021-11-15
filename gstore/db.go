package gstore

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

const (
	defaultPoolMaxIdle     = 10                                 // 连接池空闲连接数量
	defaultPoolMaxOpen     = 100                                // 连接池最大连接数量
	defaultConnMaxLifeTime = time.Second * time.Duration(21600) // MySQL默认长连接时间为8个小时,所以我们设置连接可重用的时间为6小时最为合理
	defaultConnMaxIdleTime = time.Second * time.Duration(600)   // 设置连接10分钟没有用到就断开连接(内存要求较高可降低该值)
)

type DbPoolCfg struct {
	MaxIdleConn int `json:"max_idle_conn"` //空闲连接数
	MaxOpenConn int `json:"max_open_conn"` //最大连接数
	MaxLifeTime int `json:"max_life_time"` //连接可重用的最大时间
	MaxIdleTime int `json:"max_idle_time"` //在关闭连接之前,连接可能处于空闲状态的最大时间
}

type dbConfig struct {
	poolCfg *DbPoolCfg
	gormCfg *gorm.Config
}

type DbConnFunc func(cfg *dbConfig)

// InitDB init db
func InitDB(dsn string, DbCfgFunc ...DbConnFunc) *gorm.DB {
	var err error
	var cfg dbConfig

	for _, f := range DbCfgFunc {
		f(&cfg)
	}

	if cfg.gormCfg == nil {
		cfg.gormCfg = &gorm.Config{}
	}

	Db, err := gorm.Open(mysql.Open(dsn), cfg.gormCfg)
	if err != nil {
		log.Printf("[gstore_db] mysql open fail, err=%s", err)
		panic(err)
	}

	cfg.setDefaultPoolConfig(Db)

	err = DbSurvive(Db)
	if err != nil {
		log.Printf("[gstore_db] mysql survive fail, err=%s", err)
		panic(err)
	}

	log.Printf("[gstore_db] mysql success")
	return Db
}

// DbSurvive mysql survive
func DbSurvive(db *gorm.DB) error {
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

// SetPoolConfig set pool config
func SetPoolConfig(cfg DbPoolCfg) DbConnFunc {
	return func(c *dbConfig) {
		c.poolCfg = &cfg
	}
}

// SetGormConfig set gorm config
func SetGormConfig(cfg *gorm.Config) DbConnFunc {
	return func(c *dbConfig) {
		c.gormCfg = cfg
	}
}

func (c *dbConfig) setDefaultPoolConfig(db *gorm.DB) {
	d, err := db.DB()
	if err != nil {
		log.Printf("[gstore_db] mysql db fail, err=%s", err)
		panic(err)
	}
	var cfg = c.poolCfg
	if cfg == nil {
		d.SetMaxOpenConns(defaultPoolMaxOpen)
		d.SetMaxIdleConns(defaultPoolMaxIdle)
		d.SetConnMaxLifetime(defaultConnMaxLifeTime)
		d.SetConnMaxIdleTime(defaultConnMaxIdleTime)
		return
	}

	if cfg.MaxOpenConn == 0 {
		d.SetMaxOpenConns(defaultPoolMaxOpen)
	} else {
		d.SetMaxOpenConns(cfg.MaxOpenConn)
	}

	if cfg.MaxIdleConn == 0 {
		d.SetMaxIdleConns(defaultPoolMaxIdle)
	} else {
		d.SetMaxIdleConns(cfg.MaxIdleConn)
	}

	if cfg.MaxLifeTime == 0 {
		d.SetConnMaxLifetime(defaultConnMaxLifeTime)
	} else {
		d.SetConnMaxLifetime(time.Second * time.Duration(cfg.MaxLifeTime))
	}

	if cfg.MaxIdleTime == 0 {
		d.SetConnMaxIdleTime(defaultConnMaxIdleTime)
	} else {
		d.SetConnMaxIdleTime(time.Second * time.Duration(cfg.MaxIdleTime))
	}
}
