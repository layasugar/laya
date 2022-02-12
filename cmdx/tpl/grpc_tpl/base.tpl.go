package grpc_tpl

const BaseTpl = `package controllers

import (
	"fmt"
	"github.com/layasugar/laya"
	"{{.goModName}}/global"
	"github.com/layasugar/laya/genv"
)

// Ctrl the controllers with some useful and common function
var Ctrl = &BaseCtrl{}

type BaseCtrl struct {}

// Version version
func (ctrl *BaseCtrl) Version(ctx *laya.WebContext) {
	res := fmt.Sprintf("%s version: %s\napp_url: %s", genv.AppName(), genv.AppVersion(), genv.AppUrl())
	ctx.InfoF("测试日志%s", "hello world")
	_, _ = ctx.Writer.Write([]byte(res))
	return
}
`
