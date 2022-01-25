package protocol

import (
	"fmt"
	"net"
	"strconv"
	"sync"

	"gitlab.xthktech.cn/bs/gxe/cal/context"
	"gitlab.xthktech.cn/bs/gxe/cal/service"
	"gitlab.xthktech.cn/bs/gxe/rpc/pbrpc"
	"gitlab.xthktech.cn/bs/gxe/utils/produce"
	
	"github.com/two/tcpool"
)

var pbTc = &tcpool.Pool{
}
var pbMu sync.Mutex

// PbRPCRequest PbRpc 请求
type PbRPCRequest struct {
	CustomHost string
	CustomPort int
	Data       *pbrpc.Package
	LogID      string

	Ctx context.RequestContext
}

// PbRPCHead PbRPC头
type PbRPCHead struct {
	Header     pbrpc.Header
	Meta       pbrpc.RpcMeta
	Attachment []byte
}

// PbRPCProtocol pbrpc 协议
type PbRPCProtocol struct {
	serv       service.Service
	originReq  *PbRPCRequest
	curConnKey tcpool.Key
	logID      int64
}

// NewPbRPCProtocol 创建 PbRPC协议
func NewPbRPCProtocol(ctx *context.Context, serv service.Service, req *PbRPCRequest) (hp *PbRPCProtocol, err error) {
	hp = &PbRPCProtocol{
		serv:      serv,
		originReq: req,
		logID:     req.Data.GetLogId(),
	}
	ctx.ReqContext = req.Ctx

	hp.initLogID(ctx)
	return
}

func (hp *PbRPCProtocol) initLogID(ctx *context.Context) {
	logID := hp.logID

	if logID == 0 {
		if ctx.ReqContext != nil {
			logID, _ = strconv.ParseInt(ctx.ReqContext.GetLogID(), 10, 64)
		}
	}

	if logID == 0 {
		logID = produce.NewLogIDInt()
	}

	hp.logID, ctx.LogID = logID, logID
	hp.originReq.Data.SetLogId(logID)
}

// Do 执行
func (hp *PbRPCProtocol) Do(ctx *context.Context, addr *service.Addr) (rsp *Response, err error) {
	conn, err := hp.getClient(ctx, addr)
	if err != nil {
		return nil, err
	}
	if hp.serv.GetReuse() {
		c := tcpool.Func{
			Factory: func() (interface{}, error) {
				d := net.Dialer{Timeout: hp.serv.GetConnTimeout()}
				return d.Dial("tcp", addr.String())
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
	originRsp := pbrpc.NewPackage()
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

func (hp *PbRPCProtocol) getClient(ctx *context.Context, addr *service.Addr) (conn net.Conn, err error) {
	url := ""
	if hp.originReq.CustomHost != "" {
		url = fmt.Sprintf("%s:%d", hp.originReq.CustomHost, hp.originReq.CustomPort)
		ctx.CurRecord().IDC = "custom"
		ctx.CurRecord().Host = url
	} else {
		url = addr.String()
		ctx.CurRecord().IDC = addr.IDC
		ctx.CurRecord().Host = addr.GetHostName()
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
