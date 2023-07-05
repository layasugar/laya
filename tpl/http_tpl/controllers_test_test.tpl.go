package http_tpl

const ControllersTestTestTpl = `package test

import (
	"errors"
	"github.com/layasugar/laya"

	"{{.goModName}}/models/page/test"
)

// HttpTraceTest 测试http请求和链路追踪(http_to_http http_to_grpc)
func (ctrl *controller) HttpTraceTest(ctx *laya.WebContext) {
	// 参数绑定
	var pm test.Req
	err := ctx.ShouldBindJSON(&pm)
	if err != nil {
		ctrl.Fail(ctx, err)
		return
	}

	// 参数校验
	if pm.Kind != 1 && pm.Kind != 2 {
		ctrl.Fail(ctx, errors.New("kind 只能是1,2"))
	}

	// 业务处理
	res, err := test.HttpTraceTest(ctx, pm)
	if err != nil {
		ctrl.Fail(ctx, err)
		return
	}

	// 响应
	ctrl.Suc(ctx, res)
}
`
