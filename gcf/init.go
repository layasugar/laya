package gcf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

const (
	cTimer            = 5 * time.Second // 配置重载时间, 配置文件更新5s后重载配置
	defaultConfigFile = "conf/app.toml" // 固定配置文件
)

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

	viper.SetConfigFile(f)
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	return nil
}
