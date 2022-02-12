package grpc_tpl

const ModelsDaoCalRpcTestTraceTpl = `// 请求测试文件

package rpc_test

import (
	"errors"
	"github.com/layasugar/laya"
	"github.com/layasugar/laya/gcal"
	"net/http"

	"{{.goModName}}/pb"
)

type (
	CalResp struct {
		Head gcal.HTTPHead
		Body Body
	}

	Body struct {
		StatusCode uint32 {{.tagName}}json:"status_code"{{.tagName}}
		Message    string {{.tagName}}json:"message"{{.tagName}}
		Data       Data   {{.tagName}}json:"data"{{.tagName}}
		RequestID  string {{.tagName}}json:"request_id"{{.tagName}}
	}

	Data struct {
		Code string {{.tagName}}json:"code"{{.tagName}}
	}

	RpcData struct {
		Message string {{.tagName}}json:"message"{{.tagName}}
	}
)

var path = "/server-b/fast"
var serviceName1 = "http_test"
var serviceName2 = "grpc_test"

// HttpTraceTest Http测试, body是interface可以发送任何类型的数据
func HttpTraceTest(ctx *laya.GrpcContext) (*Data, error) {
	ctx.InfoF("开始请求了, %s", "aaaa")
	req := gcal.HTTPRequest{
		Method: "POST",
		Path:   path,
		Body: map[string]string{
			"data": "success",
		},
		Ctx: ctx,
		Header: map[string][]string{
			"Host": []string{"12312"},
		},
	}
	response := CalResp{}
	err := gcal.Do(serviceName1, req, &response, gcal.JSONConverter)

	// 状态码非 200
	if response.Head.StatusCode != http.StatusOK {
		return &response.Body.Data, errors.New("NETWORK_ERROR")
	}
	ctx.InfoF("结束请求了, %s", "bbbb")
	return &response.Body.Data, err
}

// RpcTraceTest rpc测试
func RpcTraceTest(ctx *laya.GrpcContext) (*RpcData, error) {
	conn := gcal.GetRpcConn(serviceName2)
	if conn == nil {
		return nil, errors.New("连接不存在")
	}

	c := pb.NewGreeterClient(conn)

	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: "q1mi"})
	if err != nil {
		return nil, err
	}
	return &RpcData{
		Message: res.Message,
	}, err
}
`