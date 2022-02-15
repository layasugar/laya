package server_tpl

const ModelsDataTestTestTraceTpl = `package test

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/models/dao/cal/task_test"
)

func TaskToHttpTest(ctx *laya.Context) (*Rsp, error) {
	d, err := task_test.TaskToHttpTest(ctx)
	if err != nil {
		return nil, err
	}

	var res = Rsp{
		Code: d.Code,
	}

	return &res, nil
}

func TaskToRpcTest(ctx *laya.Context) (*Rsp, error) {
	d, err := task_test.TaskToGrpcTest(ctx)
	if err != nil {
		return nil, err
	}

	var res = Rsp{
		Code: d.Message,
	}

	return &res, nil
}
`
