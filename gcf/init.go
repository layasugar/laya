package gcf

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

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

	// 加载环境变量
	viper.AutomaticEnv()

	return nil
}
