package gconf

import (
	"encoding/json"
	"io/ioutil"
)

var path = "./gconf/app.toml"

var c *Config

type Config struct {
	BaseConf
	LogConf    LogConf    `toml:"log"`
	CacheConf  CacheConf  `toml:"cache"`
	DBConf     DBConf     `toml:"mysql"`
	RdbConf    RdbConf    `toml:"redis"`
	MdbConf    MdbConf    `toml:"mongo"`
	KafkaConf  KafkaConf  `toml:"kafka_conf"`
	I18nConfig I18nConfig `toml:"gi18n"`
	TraceConf  TraceConf  `json:"zipkin"`
	DingConf   DingConf   `json:"ding"`
}

type (
	LogConf struct {
		Open       bool   `toml:"open"`
		Driver     string `toml:"driver"`      //驱动分控制台和文件
		Path       string `toml:"path"`        //文件路径
		LogLevel   string `toml:"log_level"`   //日志等级分info,error,warn
		MaxSize    int    `toml:"max_size"`    //日志文件最大MB
		MaxAge     int    `toml:"max_age"`     //保留旧文件的最大天数
		MaxBackups int    `toml:"max_backups"` //保留旧文件的最大个数
	}
	CacheConf struct {
	}
	BaseConf struct {
		AppName    string `toml:"app_name"`    //app名称
		HttpListen string `toml:"http_listen"` //http监听端口
		RunMode    string `toml:"run_mode"`    //运行模式
		AppVersion string `toml:"version"`     //app版本号
		AppUrl     string `json:"app_url"`     //当前路由
		GinLog     string `json:"gin_log"`     //gin_log日志
		ParamsLog  bool   `json:"params_log"`  //是否开启请求参数和返回参数打印
	}
	DBConf struct {
		Open            bool   `toml:"open"`            //是否开启
		MaxIdleConn     int    `toml:"maxIdleConn"`     //空闲连接数
		MaxOpenConn     int    `toml:"maxOpenConn"`     //最大连接数
		ConnMaxLifetime int    `toml:"connMaxLifetime"` //连接时长
		Dsn             string `toml:"dsn"`             //dsn
	}
	RdbConf struct {
		Open        bool   `toml:"open"`        //是否开启
		DB          int    `toml:"db"`          //默认连接库
		PoolSize    int    `toml:"poolSize"`    //连接数量
		MaxRetries  int    `toml:"maxRetries"`  //最大重试次数
		IdleTimeout int    `toml:"idleTimeout"` //空闲链接超时时间(单位：time.Second)
		Addr        string `toml:"addr"`        //DSN
		Pwd         string `toml:"pwd"`         //密码
	}
	MdbConf struct {
		Open        bool   `toml:"open"`        //是否开启
		DSN         string `toml:"dsn"`         //dsn
		MinPoolSize uint64 `toml:"minPoolSize"` //连接池最小连接数
		MaxPoolSize uint64 `toml:"maxPoolSize"` //连接池最大连接数
	}
	KafkaConf struct {
		Open      bool     `toml:"open"`
		Brokers   []string `toml:"brokers"`
		CertFile  string   `toml:"cert_file"`
		KeyFile   string   `toml:"key_file"`
		CaFile    string   `toml:"ca_file"`
		VerifySsl bool     `toml:"verify_ssl"`
	}
	I18nConfig struct {
		Open        bool   `toml:"open"`
		DefaultLang string `toml:"defaultLang"`
		Path        string `toml:"path"`
	}
	TraceConf struct {
		Open            bool   `json:"open"`
		ServiceName     string `json:"service_name"`     //服务名
		ServiceEndpoint string `json:"service_endpoint"` //当前服务节点
		ZipkinAddr      string `json:"zipkin_addr"`      //zipkin地址
		Mod             uint64 `json:"mod"`              //采样率,0==不进行链路追踪，1==全量。值越大，采样率越低，对性能影响越小
	}
	DingConf struct {
		Open      bool   `json:"open"`
		RobotKey  string `json:"robot_key"`
		RobotHost string `json:"robot_host"`
	}
)

func InitConfig(confPath string) error {
	c = new(Config)
	r, err := ioutil.ReadFile(confPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(r, &c)
	if err != nil {
		panic(err)
	}

	return nil
}

func GetBaseConf() BaseConf     { return c.BaseConf }
func GetLogConf() LogConf       { return c.LogConf }
func GetDBConf() DBConf         { return c.DBConf }
func GetRdbConf() RdbConf       { return c.RdbConf }
func GetMdbConf() MdbConf       { return c.MdbConf }
func GetI18nConfig() I18nConfig { return c.I18nConfig }
func GetAppName() string        { return c.BaseConf.AppName }
func GetHttpListen() string     { return c.BaseConf.HttpListen }
func GetRunMode() string        { return c.BaseConf.RunMode }
func GetAppVersion() string     { return c.BaseConf.AppVersion }
