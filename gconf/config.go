package gconf

import (
	"encoding/json"
	"io/ioutil"
)

var path = "./gconf/app.json"

var c *Config

type Config struct {
	BaseConf
	LogConf   *LogConf  `json:"log"`
	CacheConf CacheConf `json:"cache"`
	DBConf    DBConf    `json:"mysql"`
	RdbConf   RdbConf   `json:"redis"`
	MdbConf   MdbConf   `json:"mongo"`
	KafkaConf KafkaConf `json:"kafka_conf"`
	TraceConf TraceConf `json:"zipkin"`
	DingConf  DingConf  `json:"ding"`
}

type (
	BaseConf struct {
		AppName    string `json:"app_name"`    //app名称
		HttpListen string `json:"http_listen"` //http监听端口
		RunMode    string `json:"run_mode"`    //运行模式
		AppVersion string `json:"version"`     //app版本号
		AppUrl     string `json:"app_url"`     //当前路由
		GinLog     string `json:"gin_log"`     //gin_log日志
		ParamsLog  bool   `json:"params_log"`  //是否开启请求参数和返回参数打印
	}
	LogConf struct {
		Path       string `json:"path"`        //文件路径
		MaxSize    int    `json:"max_size"`    //日志文件最大MB
		MaxAge     int    `json:"max_age"`     //保留旧文件的最大天数
		MaxBackups int    `json:"max_backups"` //保留旧文件的最大个数
	}
	CacheConf struct {
	}

	DBConf struct {
		Open            bool   `json:"open"`            //是否开启
		MaxIdleConn     int    `json:"maxIdleConn"`     //空闲连接数
		MaxOpenConn     int    `json:"maxOpenConn"`     //最大连接数
		ConnMaxLifetime int    `json:"connMaxLifetime"` //连接时长
		Dsn             string `json:"dsn"`             //dsn
	}
	RdbConf struct {
		Open        bool   `json:"open"`        //是否开启
		DB          int    `json:"db"`          //默认连接库
		PoolSize    int    `json:"poolSize"`    //连接数量
		MaxRetries  int    `json:"maxRetries"`  //最大重试次数
		IdleTimeout int    `json:"idleTimeout"` //空闲链接超时时间(单位：time.Second)
		Addr        string `json:"addr"`        //DSN
		Pwd         string `json:"pwd"`         //密码
	}
	MdbConf struct {
		Open        bool   `json:"open"`        //是否开启
		DSN         string `json:"dsn"`         //dsn
		MinPoolSize uint64 `json:"minPoolSize"` //连接池最小连接数
		MaxPoolSize uint64 `json:"maxPoolSize"` //连接池最大连接数
	}
	KafkaConf struct {
		Open      bool     `json:"open"`
		Brokers   []string `json:"brokers"`
		CertFile  string   `json:"cert_file"`
		KeyFile   string   `json:"key_file"`
		CaFile    string   `json:"ca_file"`
		VerifySsl bool     `json:"verify_ssl"`
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

	//判断有没有日志配置如果没有赋初始值
	if c.LogConf == nil {
		c.LogConf = &LogConf{
			Path:       "/home/logs/app/" + c.BaseConf.AppName + "/app.log",
			MaxSize:    32,
			MaxAge:     90,
			MaxBackups: 300,
		}
	}

	return nil
}

func GetBaseConf() BaseConf { return c.BaseConf }
func GetLogConf() LogConf   { return *c.LogConf }
func GetDBConf() DBConf     { return c.DBConf }
func GetRdbConf() RdbConf   { return c.RdbConf }
func GetMdbConf() MdbConf   { return c.MdbConf }
func GetAppName() string    { return c.BaseConf.AppName }
func GetHttpListen() string { return c.BaseConf.HttpListen }
func GetRunMode() string    { return c.BaseConf.RunMode }
func GetAppVersion() string { return c.BaseConf.AppVersion }
