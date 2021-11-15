package gconf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	ctimer = 5 * time.Second
)

var C = viper.New()
var cFunc []func()
var t *time.Timer

// InitConfig 初始化配置, cfp config path配置路径或者带.的具体文件名
// tp 是配置文件类型
func InitConfig(cfp string) error {
	if len(cfp) == 0 {
		// 先加载默认的配置
		C.SetConfigFile(defaultConfigFile)
		err := C.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	} else {
		if !strings.Contains(cfp, ".") {
			panic(fmt.Errorf("Fatal error config file: %s \n", "非法路径"))
		}
		tree := strings.Split(cfp, ".")
		exp := tree[len(tree)-1:]
		if !strings.Contains(defaultConfigType, exp[0]) {
			panic(fmt.Errorf("Fatal error config file: %s \n", "非法后缀"))
		}

		C.SetConfigFile(cfp)
		err := C.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	C.WatchConfig()
	return nil
}

// OnConfigCharge 注册配置改变后的处理
func OnConfigCharge() {
	var f = func(in fsnotify.Event) {
		if t == nil {
			t = time.NewTimer(ctimer)
		} else {
			t.Reset(ctimer)
		}

		go func() {
			<-t.C
			// 只处理写入事件
			if in.Op&fsnotify.Write == fsnotify.Write {
				for _, item := range cFunc {
					item()
				}
			}
		}()
	}
	C.OnConfigChange(f)
}

// RegisterConfigCharge 可以在程序启动前注册多个配置变化函数
func RegisterConfigCharge(f ...func()) {
	cFunc = append(cFunc, f...)
}
