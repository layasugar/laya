package gcf

import "time"

const (
	cTimer            = 5 * time.Second // 配置重载时间, 配置文件更新5s后执行重载函数
	defaultConfigFile = "conf/app.toml" // 固定配置文件
)

var configChargeHandleFunc []func()
var t *time.Timer
