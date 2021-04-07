package gconf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var path = "./gconf/app.json"

var c *Config
var mc map[string]json.RawMessage

type Config struct {
	BaseConf  BaseConf  `json:"app"`
	LogConf   *LogConf  `json:"log"`
	CacheConf MemConf   `json:"cache"`
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
	MemConf struct {
		DefaultExp int64 `json:"default_exp"`
		Cleanup    int64 `json:"cleanup"`
	}
	DBConf struct {
		Dsn             string `json:"dsn"`             //dsn
		MaxIdleConn     int    `json:"maxIdleConn"`     //空闲连接数
		MaxOpenConn     int    `json:"maxOpenConn"`     //最大连接数
		ConnMaxLifetime int    `json:"connMaxLifetime"` //连接时长
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

func InitConfig(cfp string) error {
	c = new(Config)
	r, err := ioutil.ReadFile(cfp)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(r, &c)
	if err != nil {
		log.Printf("%s file load err,err is %s\n", cfp, err.Error())
		panic(err)
	}
	mc = make(map[string]json.RawMessage)
	err = json.Unmarshal(r, &mc)
	if err != nil {
		log.Printf("%s file load err,err is %s\n", cfp, err.Error())
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

func GetBaseConf() (*BaseConf, error) {
	var bcf = BaseConf{}
	raw, ok := mc["base"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &bcf)
	if err != nil {
		return nil, err
	}
	return &bcf, nil
}
func GetLogConf() (*LogConf, error) {
	var logc = LogConf{}
	raw, ok := mc["log"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &logc)
	if err != nil {
		return nil, err
	}
	return &logc, nil
}
func GetDBConf(k string) (*DBConf, error) {
	var dbc = DBConf{}
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &dbc)
	if err != nil {
		return nil, err
	}
	return &dbc, nil
}
func GetRdbConf(k string) (*RdbConf, error) {
	var rdbc = RdbConf{}
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &rdbc)
	if err != nil {
		return nil, err
	}
	return &rdbc, nil
}
func GetMdbConf(k string) (*MdbConf, error) {
	var mdbc = MdbConf{}
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &mdbc)
	if err != nil {
		return nil, err
	}
	return &mdbc, nil
}
func GetMemConf() (*MemConf, error) {
	var memc = MemConf{}
	raw, ok := mc["mem"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &memc)
	if err != nil {
		return nil, err
	}
	return &memc, nil
}

func GetHttpListen() string { return c.BaseConf.HttpListen }
func GetRunMode() string    { return c.BaseConf.RunMode }
func GetAppVersion() string { return c.BaseConf.AppVersion }
