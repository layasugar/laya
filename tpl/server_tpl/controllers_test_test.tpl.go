package server_tpl

const ControllersTestTestTpl = `package test

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/models/page/test"
)

// Task 测试任务处理
func (ctrl *controller) Task(ctx *laya.Context, kind uint8) {
	// 日志记录
	ctx.InfoF("测试基础任务")

	// 业务处理
	res, err := test.TaskTest(ctx, test.Req{
		Kind: kind,
	})
	if err != nil {
		ctx.ErrorF("任务错误%v", err)
		return
	}

	ctx.InfoF("任务完成%v", res)

	// 调用ctx.Done, 结束任务
	ctx.Done()
}
`
