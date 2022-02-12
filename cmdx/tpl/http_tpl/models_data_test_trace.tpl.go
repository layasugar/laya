package http_tpl

const ModelsDataTestTestTraceTpl = `package test

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/models/dao/cal/http_test"
)

func HttpToHttpTraceTest(ctx *laya.WebContext) (*Rsp, error) {
	d, err := http_test.HttpToHttpTraceTest(ctx)
	if err != nil {
		return nil, err
	}

	var res = Rsp{
		Code: d.Code,
	}

	return &res, nil
}
`
