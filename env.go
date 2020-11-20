package laya

// environment
var ENV string

var DelayServer string

var RedisConf struct {
	Open        bool   `json:"open"`
	DB          int    `json:"db"`
	Addr        string `json:"addr"`
	Pwd         string `json:"pwd"`
	PoolSize    int    `json:"poolSize"`
	MaxRetries  int    `json:"maxRetries"`
	IdleTimeout int    `json:"idleTimeout"`
}

var MysqlConf struct {
	Open            bool   `json:"open"`
	Dsn             string `json:"dsn"`
	MaxIdleConn     int    `json:"maxIdleConn"`
	MaxOpenConn     int    `json:"maxOpenConn"`
	ConnMaxLifetime int    `json:"connMaxLifetime"`
}
