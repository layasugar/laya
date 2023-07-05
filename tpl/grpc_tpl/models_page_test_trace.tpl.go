package grpc_tpl

const ModelsPageTestTraceTpl = `package test

import (
	"github.com/layasugar/laya"

	"{{.goModName}}/models/data/test"
	"{{.goModName}}/pb"
)

type (
	Req struct {
		Kind uint8 {{.tagName}}json:"kind"{{.tagName}} // 1表示发起http请求, 2表示发起rpc请求
	}

	Rsp struct {
		Code string {{.tagName}}json:"code"{{.tagName}}
	}
)

func RpcTraceTest(ctx *laya.GrpcContext, pm *pb.GrpcTraceTestReq) (*Rsp, error) {
	var res Rsp
	switch pm.Kind {
	case 1:
		d, err := test.RpcToHttpTraceTest(ctx)
		if err != nil {
			return nil, err
		}

		res.Code = d.Code
	case 2:
		d, err := test.RpcToRpcTraceTest(ctx)
		if err != nil {
			return nil, err
		}

		res.Code = d.Code
	}

	return &res, nil
}
`
