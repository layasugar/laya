package http_tpl

const ModelsPageTestTraceTpl = `package test

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/models/data/test"
)

type (
	Req struct {
		Kind uint8 {{.tagName}}json:"kind"{{.tagName}} // 1表示发起http请求, 2表示发起rpc请求
	}

	Rsp struct {
		Code string {{.tagName}}json:"code"{{.tagName}}
	}
)

func HttpTraceTest(ctx *laya.WebContext, pm Req) (*Rsp, error) {
	var res Rsp
	switch pm.Kind {
	case 1:
		d, err := test.HttpToHttpTraceTest(ctx)
		if err != nil {
			return nil, err
		}

		res.Code = d.Code
	}

	return &res, nil
}
`
