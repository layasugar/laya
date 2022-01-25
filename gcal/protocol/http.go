package protocol

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync"
	"time"

	"gitlab.xthktech.cn/bs/gxe/cal/context"
	"gitlab.xthktech.cn/bs/gxe/cal/converter"
	"gitlab.xthktech.cn/bs/gxe/cal/service"
	"gitlab.xthktech.cn/bs/gxe/utils/produce"
	"gitlab.xthktech.cn/bs/gxe/version"
)

// UA just a flag
const UA = "CAL/" + version.VERSION + " (gxe cal http client)"
const HttpClientAlive time.Duration = 5 * time.Minute

// HTTPRequest http requst 对象，cal.Cal 函数必须传递这个类型的变量
type HTTPRequest struct {
	CustomHost string
	CustomPort int

	// TODO http.Header
	Header      map[string][]string
	Method      string
	Body        interface{}
	Path        string
	QueryParams url.Values
	LogID       string

	Converter converter.ConverterType

	Ctx context.RequestContext
}

// HTTPHead HTTPResponse，兼容历史
type HTTPHead struct {
	Status     string
	StatusCode int
	Proto      string
	// TODO http.Header
	Header        map[string][]string
	ContentLength int64
}

// HTTPProtocol http 协议
type HTTPProtocol struct {
	protocol string
	serv     service.Service
	logID    string

	originReq *HTTPRequest
	RawReq    *http.Request
	// OriginRsp *http.Response
}

// Protocol 返回类型
func (hp *HTTPProtocol) Protocol() string {
	return hp.protocol
}

// initLogID 生成logID
func (hp *HTTPProtocol) initLogID(ctx *context.Context) {
	logID := hp.originReq.LogID

	if logID == "" {
		if ctx.ReqContext != nil {
			logID = ctx.ReqContext.GetLogID()
		}
	}

	if logID == "" {
		//直接随机生成
		logID = produce.NewLogID()
	}

	hp.logID, ctx.LogID = logID, logID
}

// NewHTTPProtocol 创建一个 Http Protocol
func NewHTTPProtocol(ctx *context.Context, serv service.Service, req *HTTPRequest, isHTTPS bool) (hp *HTTPProtocol, err error) {
	hp = &HTTPProtocol{
		serv:      serv,
		originReq: req,
		protocol:  "http",
	}
	if isHTTPS {
		hp.protocol = "https"
	}

	ctx.ReqContext = req.Ctx
	hp.initLogID(ctx)
	ctx.Method = strings.ToLower(req.Method)

	hp.RawReq = &http.Request{
		Method:     strings.ToUpper(req.Method),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     req.Header,
		Body:       http.NoBody,
		GetBody:    func() (io.ReadCloser, error) { return http.NoBody, nil },
		// URL:        u,
		// Host:       u.Host,
	}
	if hp.RawReq.Header == nil {
		hp.RawReq.Header = make(http.Header)
	}

	bb := []byte{}
	if req.Body != nil {
		conver, _ := converter.GetConverter(req.Converter)
		if conver == nil {
			err = fmt.Errorf("get convert %s failed", req.Converter)
			return
		}

		ctx.PackStatis.StartPoint = time.Now()
		bb, err = conver.Pack(req.Body)
		ctx.PackStatis.StopPoint = time.Now()
		if err != nil {
			return
		}
	}

	if len(bb) > 0 {
		body := strings.NewReader(string(bb))
		hp.RawReq.ContentLength = int64(body.Len())
		hp.RawReq.Body = ioutil.NopCloser(body)
		snapshot := *body
		hp.RawReq.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return ioutil.NopCloser(&r), nil
		}
	}

	ctx.ReqLen = hp.RawReq.ContentLength

	commonHeaders, err := serv.HeaderInfo()
	if err != nil {
		return nil, err
	}

	hp.RawReq.Header.Set(commonHeaders[service.HeaderLogIDKey], hp.logID)
	delete(commonHeaders, service.HeaderLogIDKey)

	// 优先使用用户配置 Host
	if hosts := req.Header["Host"]; len(hosts) > 0 {
		hp.RawReq.Host = hosts[0]
	} else if host, ok := commonHeaders["Host"]; ok {
		// 使用BNS的 Host
		hp.RawReq.Host = host
		// TODO context log
		delete(commonHeaders, "Host")
	}

	// If the user doesn't set User-Agent, set the default User-Agent
	if hp.RawReq.Header.Get("User-Agent") == "" {
		hp.RawReq.Header.Set("User-Agent", UA)
	}

	return
}

// Do 发送请求
func (hp *HTTPProtocol) Do(ctx *context.Context, addr *service.Addr) (rsp *Response, err error) {
	var host string
	if hp.originReq.CustomHost != "" {
		host = fmt.Sprintf("%s:%d", hp.originReq.CustomHost, hp.originReq.CustomPort)
		ctx.CurRecord().IDC = "custom"
	} else {
		host = addr.String()
		ctx.CurRecord().IDC = addr.IDC
	}
	ctx.CurRecord().IPPort = host

	path := hp.originReq.Path
	// 天路
	if path == "" {
		path = addr.Path
	}

	if len(hp.originReq.QueryParams) > 0 {
		// 不接受又是path有参数又设置get参数的行为
		path += "?"
		path += hp.originReq.QueryParams.Encode()
	}

	fullPath := fmt.Sprintf("%s://%s/%s", hp.Protocol(), host, path)
	u, err := url.Parse(fullPath)
	if err != nil {
		return nil, err
	}

	ctx.CurRecord().Path = u.Path

	hp.RawReq.URL = u
	if hp.RawReq.Host == "" {
		hp.RawReq.Host = u.Host // BNS或者用户没有设置
	}

	ctx.CurRecord().Host = hp.RawReq.Host

	//GetConn为从域名解析开始算起，connectStart为对目标地址的ip发起连接请求，
	//write时间计算方法为连接建立完毕到wirte request结束，这个不太准确
	//其中部分行为可能会执行多次，这里只记录最后一次的

	// +---------+-----------+----------+---------------+--------------+----------+----------------+---------------+-----------------------+---------------------+--------------+
	// |         |           |          |               |              |          |                |               |                       |                     |              |
	// | GetConn |  DNSStart |  DNSDone |  ConnectStart |  ConnectDone |  GotConn |  WroteHeaders  |  WroteRequest |  GotFirstResponseByte |  GotAllResponseByte | CloseConnect |
	// |         |           |          |               |              |          |                |               |                       |                     |              |
	// +---------+-----------+----------+---------------+--------------+----------+----------------+---------------+-----------------------+---------------------+--------------+
	// |<-----------------------------connect------------------------------------>|
	// 																		      |<-------------write------------>|
	// 																											   |<---------------------read------------------>|
	// |<--------------------------------------------------------------------------talk------------------------------------------------------------------------->|
	// |<--------------------------------------------------------------------------------cost---------------------------------------------------------------------------------->|
	trace := &httptrace.ClientTrace{
		GetConn: func(hostport string) {
			ctx.TimeStatisStart("connect")
			ctx.TimeStatisStart("talk")
			ctx.CurRecord().RecordTimePoint("req_start_time")
		},
		GotConn: func(info httptrace.GotConnInfo) {
			ctx.TimeStatisStop("connect")
			ctx.TimeStatisStart("write")
		},
		ConnectStart: func(network, addr string) {
			//ctx.TimeStatisStart("talk")
		},
		ConnectDone: func(network, addr string, err error) {
			//ctx.TimeStatisStart("talk")
		},
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			ctx.TimeStatisStart("dnslookup")
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			ctx.TimeStatisStop("dnslookup")
		},

		WroteRequest: func(writeRequestInfo httptrace.WroteRequestInfo) {
			ctx.TimeStatisStop("write")
		},
	}

	httpReq := hp.RawReq.WithContext(httptrace.WithClientTrace(hp.RawReq.Context(), trace))

	client, err := hp.getClient(ctx)
	if err != nil {
		return nil, err
	}
	if hp.serv.GetReuse() {
		defer hp.tryReuseClient(client)
	}

	originRsp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	defer func() {
		originRsp.Body.Close()
		ctx.TimeStatisStop("talk")
	}()

	ctx.CurRecord().RspCode = originRsp.StatusCode

	ctx.TimeStatisStart("read")
	raw, err := ioutil.ReadAll(originRsp.Body)
	ctx.TimeStatisStop("read")
	if err != nil {
		return nil, err
	}
	rsp = &Response{
		Request: originRsp.Request,
		Head: HTTPHead{
			Status:        originRsp.Status,
			StatusCode:    originRsp.StatusCode,
			Proto:         originRsp.Proto,
			Header:        originRsp.Header,
			ContentLength: originRsp.ContentLength,
		},
		Body:      raw,
		OriginRsp: originRsp,
	}

	ctx.RspLen = int64(len(raw))

	return
}

func (hp *HTTPProtocol) tryReuseClient(cli *http.Client) {
	service2httpClientMap.Store(hp.serv.GetName(), cli)
}

var service2httpClientMap sync.Map
var lock sync.Mutex

func (hp *HTTPProtocol) getClient(ctx *context.Context) (client *http.Client, err error) {
	if !hp.serv.GetReuse() {
		return DefaultHTTPClientFactory(hp.serv)
	}
	cli, ok := service2httpClientMap.Load(hp.serv.GetName())
	if !ok {
		lock.Lock()
		cli, ok = service2httpClientMap.Load(hp.serv.GetName())
		if !ok {
			client, err = DefaultHTTPClientFactory(hp.serv)
			service2httpClientMap.Store(hp.serv.GetName(), client)
			lock.Unlock()
			go func(name string) {
				select {
				case <-time.After(HttpClientAlive):
					service2httpClientMap.Delete(name)
				}
			}(hp.serv.GetName())
			return
		}
		lock.Unlock()
	}
	return cli.(*http.Client), nil
}

// DefaultHTTPClientFactory 默认的 http client factory
var DefaultHTTPClientFactory = func(serv service.Service) (cli *http.Client, err error) {
	var proxyURL *url.URL
	var proxy string

	if serv.GetHTTPProxy() != "" {
		proxy = serv.GetHTTPProxy()
	} else if serv.GetHTTPSProxy() != "" {
		proxy = serv.GetHTTPSProxy()
	}
	if proxy != "" {
		proxyURL, err = url.Parse(proxy)
		if err != nil {
			return nil, err
		}
	}

	perHost := -1
	if serv.GetReuse() {
		// 如果不加这个，每发一个请求会产生一个gorouter，85秒后回收
		// 见 $GOROOT/src/net/http/transport.goa line :765
		perHost = 2
	}

	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL), //代理设置
			DialContext: (&net.Dialer{
				Timeout:   serv.GetConnTimeout(), //连接超时时间
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConnsPerHost:   perHost,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		//总的超时时间
		Timeout: serv.GetTotalTimeout(),
	}, nil
}
