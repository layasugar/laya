package gconf

import (
	"encoding/json"
	"github.com/layasugar/laya/gstore"
	"io/ioutil"
	"log"
)

var mc map[string]json.RawMessage

type Config struct {
	BaseConf  BaseConf  `json:"app"`
	DBConf    DBConf    `json:"mysql"`
	RdbConf   RdbConf   `json:"redis"`
	MdbConf   MdbConf   `json:"mongo"`
	EsConf    EsConf    `json:"es"`
	KafkaConf KafkaConf `json:"kafka"`
	TraceConf TraceConf `json:"zipkin"`
	DingConf  DingConf  `json:"ding"`
}

type (
	BaseConf struct {
		AppName    string `json:"app_name"`    //app名称
		AppMode    string `json:"app_mode"`    //app运行环境
		HttpListen string `json:"http_listen"` //http监听端口
		RunMode    string `json:"run_mode"`    //运行模式
		AppVersion string `json:"version"`     //app版本号
		AppUrl     string `json:"app_url"`     //当前路由
		ParamLog   bool   `json:"param_log"`   //是否开启请求参数和返回参数打印
		LogPath    string `json:"log_path"`    //日志路径"/home/log/app"
		Pprof      bool   `json:"pprof"`       //是否开启pprof
	}
	DBConf struct {
		Dsn string `json:"dsn"` //dsn
		gstore.DbPoolCfg
	}
	RdbConf struct {
		DB          int    `json:"db"`          //默认连接库
		PoolSize    int    `json:"poolSize"`    //连接数量
		MaxRetries  int    `json:"maxRetries"`  //最大重试次数
		IdleTimeout int    `json:"idleTimeout"` //空闲链接超时时间(单位：time.Second)
		Addr        string `json:"addr"`        //DSN
		Pwd         string `json:"pwd"`         //密码
	}
	MdbConf struct {
		DSN         string `json:"dsn"`         //dsn
		MinPoolSize uint64 `json:"minPoolSize"` //连接池最小连接数
		MaxPoolSize uint64 `json:"maxPoolSize"` //连接池最大连接数
	}
	EsConf struct {
		Addr []string `json:"addr"`
		User string   `json:"user"`
		Pwd  string   `json:"pwd"`
	}
	KafkaConf struct {
		Brokers      []string `json:"brokers"`
		Topic        string   `json:"topic"`
		Group        string   `json:"group"`
		User         string   `json:"user"`
		Pwd          string   `json:"pwd"`
		CertFile     string   `json:"cert_file"`
		KeyFile      string   `json:"key_file"`
		CaFile       string   `json:"ca_file"`
		KafkaVersion string   `json:"kafka_version"`
		Scram        string   `json:"scram"`
		VerifySsl    bool     `json:"verify_ssl"`
	}
	TraceConf struct {
		ZipkinAddr string `json:"zipkin_addr"` //zipkin地址
		Mod        uint64 `json:"mod"`         //采样率,0==不进行链路追踪，1==全量。值越大，采样率越低，对性能影响越小
	}
	DingConf struct {
		RobotKey  string `json:"robot_key"`
		RobotHost string `json:"robot_host"`
	}
)

func InitConfig(cfp string) error {
	r, err := ioutil.ReadFile(cfp)
	if err != nil {
		panic(err)
	}
	mc = make(map[string]json.RawMessage)
	err = json.Unmarshal(r, &mc)
	if err != nil {
		log.Printf("%s file load err,err is %s\n", cfp, err.Error())
		panic(err)
	}
	return nil
}

func GetBaseConf() (*BaseConf, error) {
	var bcf = BaseConf{}
	raw, ok := mc["app"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &bcf)
	if err != nil {
		return nil, err
	}
	return &bcf, nil
}
func GetTraceConf() (*TraceConf, error) {
	var tc = TraceConf{}
	raw, ok := mc["zipkin"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &tc)
	if err != nil {
		return nil, err
	}
	return &tc, nil
}
func GetDingConf() (*DingConf, error) {
	var dc = DingConf{}
	raw, ok := mc["ding"]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &dc)
	if err != nil {
		return nil, err
	}
	return &dc, nil
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
func GetEsConf(k string) (*EsConf, error) {
	var esc = EsConf{}
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &esc)
	if err != nil {
		return nil, err
	}
	return &esc, nil
}
func GetKafkaConf(k string) (*KafkaConf, error) {
	var kc = KafkaConf{}
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	err := json.Unmarshal(raw, &kc)
	if err != nil {
		return nil, err
	}
	return &kc, nil
}
func GetConf(k string) (json.RawMessage, error) {
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	return raw, nil
}
