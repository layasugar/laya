// Package protocol 提供了 HTTP、HTTPS、NSHead、ProtoBuffer 协议支持
package protocol

import (
	"fmt"

	"github.com/layasugar/laya/gcal/context"
	"github.com/layasugar/laya/gcal/service"
)

// Protocoler 协议的接口
// 协议本身只完成数据请求
type Protocoler interface {
	Do(ctx *context.Context, addr string) (*Response, error)
	Protocol() string
}

var (
	_ Protocoler = &HTTPProtocol{}
	_ Protocoler = &PbRPCProtocol{}
)

// NewProtocol 创建协议
func NewProtocol(ctx *context.Context, serv service.Service, req interface{}) (p Protocoler, err error) {
	protocolName := serv.GetProtocol()

	if protocolName == "http" || protocolName == "https" {
		tmp, ok := req.(HTTPRequest)
		if !ok {
			return nil, fmt.Errorf("%s: bad request type: %T", protocolName, req)
		}
		return NewHTTPProtocol(ctx, serv, &tmp, protocolName == "https")
	}

	if protocolName == "pbrpc" {
		tmp, ok := req.(PbRPCRequest)
		if !ok {
			return nil, fmt.Errorf("%s: bad request type: %T", protocolName, req)
		}
		return NewPbRPCProtocol(ctx, serv, &tmp)
	}

	return nil, fmt.Errorf("unknow protocol: %s ", protocolName)
}

// Response 通用的返回
type Response struct {
	// Raw []byte
	Body      interface{}
	Head      interface{}
	Request   interface{}
	OriginRsp interface{}
}
