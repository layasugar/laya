package gconf

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var mc map[string]json.RawMessage

type Config struct {
	App     App     `json:"app"`
	DBConf  DBConf  `json:"mysql"`
	RdbConf RdbConf `json:"redis"`
	MdbConf MdbConf `json:"mongo"`
}

type (
	App struct {
		Name       string `json:"name"`        //app名称
		Mode       string `json:"mode"`        //app运行环境
		HttpListen string `json:"http_listen"` //http监听端口
		RunMode    string `json:"run_mode"`    //运行模式
		Version    string `json:"version"`     //app版本号
		Url        string `json:"url"`         //当前路由
		ParamLog   bool   `json:"param_log"`   //是否开启请求参数和返回参数打印
		LogPath    string `json:"log_path"`    //日志路径"/home/log/app"
		Pprof      bool   `json:"pprof"`       //是否开启pprof
	}
	DBConf struct {
		Dsn         string `json:"dsn"`           //dsn
		MaxIdleConn int    `json:"max_idle_conn"` //空闲连接数
		MaxOpenConn int    `json:"max_open_conn"` //最大连接数
		MaxLifeTime int    `json:"max_life_time"` //连接可重用的最大时间
		MaxIdleTime int    `json:"max_idle_time"` //在关闭连接之前,连接可能处于空闲状态的最大时间
	}
	RdbConf struct {
		DB          int    `json:"db"`           //默认连接库
		Pwd         string `json:"pwd"`          //密码
		Addr        string `json:"addr"`         //DSN
		PoolSize    int    `json:"pool_size"`    //连接池数量
		MaxRetries  int    `json:"max_retries"`  //最大重试次数
		IdleTimeout int    `json:"idle_timeout"` //空闲链接超时时间(单位：time.Second)
	}
	MdbConf struct {
		DSN         string `json:"dsn"`           //dsn
		MinPoolSize uint64 `json:"min_pool_size"` //连接池最小连接数
		MaxPoolSize uint64 `json:"max_pool_size"` //连接池最大连接数
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

func GetBaseConf() (*App, error) {
	var bcf = App{}
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
func GetConf(k string) (json.RawMessage, error) {
	raw, ok := mc[k]
	if !ok {
		return nil, Nil
	}
	return raw, nil
}
