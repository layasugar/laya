package http_tpl

const ControllersTestBaseTpl = `package test

import (
	"{{.goModName}}/controllers"
)

var Ctrl = &controller{}

type controller struct {
	controllers.BaseCtrl
}
`
