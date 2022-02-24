package gconf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	cTimer            = 5 * time.Second   // 配置重载时间, 配置文件更新5s后重载配置
	httpListenKey     = "app.http_listen" // http_listen
	pbRpcListenKey    = "app.pbrpc_liten" // rpc_listen
	defaultConfigFile = "conf/app.toml"   // 固定配置文件
)

var V = viper.New()
var configChargeHandleFunc []func()
var t *time.Timer

// InitConfig 初始化配置信息
func InitConfig(file string) error {
	var f string
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if file == "" {
		f = pwd + "/" + defaultConfigFile
	} else {
		f = pwd + "/" + file
	}

	V.SetConfigFile(f)
	err = V.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	V.WatchConfig()
	return nil
}

// OnConfigCharge 注册配置改变后的处理
func OnConfigCharge() {
	var f = func(in fsnotify.Event) {
		if t == nil {
			t = time.NewTimer(cTimer)
		} else {
			t.Reset(cTimer)
		}

		go func() {
			<-t.C
			// 只处理写入事件
			if in.Op&fsnotify.Write == fsnotify.Write {
				for _, item := range configChargeHandleFunc {
					item()
				}
			}
		}()
	}
	V.OnConfigChange(f)
}

// RegisterConfigCharge 可以在程序启动前注册多个配置变化函数
func RegisterConfigCharge(f ...func()) {
	configChargeHandleFunc = append(configChargeHandleFunc, f...)
}

// LoadErrMsg 根据code加载提示信息
func LoadErrMsg(code uint32) string {
	key := fmt.Sprintf("err_code.%d", code)
	s := V.GetString(key)
	return s
}

func GetConfigMap(key string) []map[string]interface{} {
	var configMaps []map[string]interface{}
	b := V.Get(key)
	switch b.(type) {
	case []interface{}:
		si := b.([]interface{})
		for _, item := range si {
			if sim, ok := item.(map[string]interface{}); ok {
				configMaps = append(configMaps, sim)
			}
		}
	}
	return configMaps
}
