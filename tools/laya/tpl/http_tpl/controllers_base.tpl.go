package http_tpl

const ControllersBaseTpl =`package controllers

import (
	"{{.goModName}}/global"
)

type BaseCtrl struct {
	global.HttpResp
}
`
