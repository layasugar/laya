package server_tpl

const MainTpl = `package main

import (
	"github.com/layasugar/laya"
	"time"

	"{{.goModName}}/controllers/test"
)

// defaultAppSetup 初始化基本服务器
func defaultAppSetup() *laya.App {
	app := laya.DefaultApp()

	return app
}

func main() {
	app := defaultAppSetup()

	// 模拟2次任务调度
	for i := 1; i < 3; i++ {
		// 生成一个ctx全局传递
		ctx := app.NewContext("", "xiaosss")

		test.Ctrl.Task(ctx, uint8(i))
	}

	time.Sleep(time.Minute)
}
`
