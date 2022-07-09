package laya

import (
	"flag"
	"log"

	"github.com/layasugar/laya/gcal"
	"github.com/layasugar/laya/gcnf"
)

const (
	servicesConfKey = "services"
)

func init() {
	var f string
	flag.StringVar(&f, "config", "", "set a config file")
	flag.Parse()

	// 初始化配置
	err := gcnf.InitConfig(f)
	if err != nil {
		panic(err)
	}

	// 初始化调用gcal
	var services = gcnf.GetConfigMap(servicesConfKey)
	if len(services) > 0 {
		err := gcal.LoadService(services)
		if err != nil {
			log.Printf("[app] init load services error: %s", err.Error())
		}
	}
}
