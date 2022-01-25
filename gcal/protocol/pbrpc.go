package protocol

import (
	"fmt"
	"github.com/layasugar/laya/gcal/context"
	"github.com/layasugar/laya/gcal/service"
	"github.com/layasugar/laya/gcal/tcpool"
	"github.com/layasugar/laya/grpcx"
	pbrpc "github.com/layasugar/laya/grpcx/pbrpc"
	"github.com/layasugar/laya/gutils"
	"net"
	"sync"
)

var pbTc = &tcpool.Pool{}
var pbMu sync.Mutex

// PbRPCRequest PbRpc 请求
type PbRPCRequest struct {
	CustomAddr string
	Data       *grpcx.Package
	TraceId    string

	Ctx context.RequestContext
}

// PbRPCHead PbRPC头
type PbRPCHead struct {
	Header     grpcx.Header
	Meta       pbrpc.RpcMeta
	Attachment []byte
}

// PbRPCProtocol pbrpc 协议
type PbRPCProtocol struct {
	serv       service.Service
	originReq  *PbRPCRequest
	curConnKey tcpool.Key
	traceId    string
}

// NewPbRPCProtocol 创建 PbRPC协议
func NewPbRPCProtocol(ctx *context.Context, serv service.Service, req *PbRPCRequest) (hp *PbRPCProtocol, err error) {
	hp = &PbRPCProtocol{
		serv:      serv,
		originReq: req,
		traceId:   req.Data.GetTraceId(),
	}
	ctx.ReqContext = req.Ctx

	hp.initTraceId(ctx)
	return
}

func (hp *PbRPCProtocol) initTraceId(ctx *context.Context) {
	traceId := hp.traceId

	if traceId == "" {
		if ctx.ReqContext != nil {
			traceId = ctx.ReqContext.GetTraceId()
		}
	}

	if traceId == "" {
		traceId = gutils.GenerateTraceId()
	}

	hp.traceId, ctx.TraceID = traceId, traceId
	hp.originReq.Data.SetTraceId(traceId)
}

// Do 执行
func (hp *PbRPCProtocol) Do(ctx *context.Context, addr string) (rsp *Response, err error) {
	conn, err := hp.getClient(ctx, addr)
	if err != nil {
		return nil, err
	}
	if hp.serv.GetReuse() {
		c := tcpool.Func{
			Factory: func() (interface{}, error) {
				d := net.Dialer{Timeout: hp.serv.GetConnTimeout()}
				return d.Dial("tcp", addr)
			},
			Close: func(v interface{}) error { return v.(net.Conn).Close() },
		}
		pbTc.SetFunc(hp.curConnKey, c)
		defer pbTc.Put(hp.curConnKey, conn)
	} else {
		defer conn.Close()
	}

	_, err = hp.originReq.Data.WriteIO(conn)
	if err != nil {
		return nil, err
	}
	originRsp := grpcx.NewPackage()
	if err := originRsp.ReadIO(conn); err != nil {
		return nil, err
	}

	rsp = &Response{
		Head: PbRPCHead{
			Header:     originRsp.Header,
			Meta:       originRsp.Meta,
			Attachment: originRsp.Attachment,
		},
		OriginRsp: originRsp,
		Body:      originRsp.Data,
		Request:   hp.originReq,
	}
	return
}

func (hp *PbRPCProtocol) getClient(ctx *context.Context, addr string) (conn net.Conn, err error) {
	url := ""
	if hp.originReq.CustomAddr != "" {
		url = fmt.Sprintf("%s", hp.originReq.CustomAddr)
		ctx.CurRecord().Host = url
	} else {
		url = addr
		ctx.CurRecord().Host = addr
	}

	ctx.CurRecord().IPPort = url

	if !hp.serv.GetReuse() {
		d := net.Dialer{Timeout: hp.serv.GetConnTimeout()}
		return d.Dial("tcp", url)
	}

	hp.curConnKey = tcpool.Key{
		Schema: "tcp",
		Addr:   url,
	}
	tcConn, err := pbTc.Get(hp.curConnKey)
	if tcConn == nil {
		d := net.Dialer{Timeout: hp.serv.GetConnTimeout()}
		return d.Dial("tcp", url)
	}
	conn = tcConn.(net.Conn)

	return
}

// Protocol 返回类型
func (hp *PbRPCProtocol) Protocol() string {
	return "pbrpc"
}
