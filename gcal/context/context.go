// Package context 提供每次 RAL 请求的上下文对象，主要用来输出日志。
package context

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.xthktech.cn/bs/gxe/env"
	"gitlab.xthktech.cn/bs/gxe/cal/log"
	"gitlab.xthktech.cn/bs/gxe/utils/produce"
)

// RequestContext Web请求的上下文
type RequestContext interface {
	GetLogID() string
	GetClientIP() string
}

// Context 用作日志记录
type Context struct {
	ReqContext RequestContext

	Caller      string
	ServiceName string
	ReqLen      int64
	RspLen      int64
	Method      string
	LogID       interface{}
	Protocol    string
	BalanceName string

	PackStatis *StatisItem

	MaxTry int

	curTryIndex   int
	invokeRecords []*InvokeRecord
	lock          *sync.RWMutex
}

// NewContext 创建一个context
func NewContext() (ctx *Context) {
	return &Context{
		PackStatis: &StatisItem{},
		LogID:      produce.NewLogIDInt(),
		lock:       new(sync.RWMutex),
	}
}

// CurRecord 当前的访问记录
func (ctx *Context) CurRecord() *InvokeRecord {
	for len(ctx.invokeRecords) < ctx.curTryIndex+1 {
		ctx.invokeRecords = append(ctx.invokeRecords, &InvokeRecord{
			timeStatis: map[string]*StatisItem{},
			index:      ctx.curTryIndex,
			timePoints: map[string]time.Time{},
			lock:       new(sync.RWMutex),
		})
	}

	return ctx.invokeRecords[ctx.curTryIndex]
}

// NextRecord 将访问记录往后移一位
func (ctx *Context) NextRecord() {
	ctx.curTryIndex++
}

// StatisItem 时间统计项
type StatisItem struct {
	StartPoint time.Time
	StopPoint  time.Time
}

// GetSpan 得到耗时
func (si *StatisItem) GetSpan() string {
	if si == nil || si.StartPoint.IsZero() || si.StopPoint.IsZero() {
		return "0"
	}

	span := si.StopPoint.Sub(si.StartPoint)
	return fmt.Sprintf("%.3f", float64(span/time.Nanosecond)/1000000)
}

// TimeStatisStart 开始一个统计项
func (ctx *Context) TimeStatisStart(topic string) {
	ctx.lock.RLock()
	if ctx.CurRecord().timeStatis[topic] != nil { // 被设置过了
		ctx.lock.RUnlock()
		return
	}
	ctx.lock.RUnlock()
	ctx.lock.Lock()
	defer ctx.lock.Unlock()
	if _, ok := ctx.CurRecord().timeStatis[topic]; !ok {
		ctx.CurRecord().timeStatis[topic] = &StatisItem{
			StartPoint: time.Now(),
		}
	}

}

// TimeStatisStop 停止一个统计项
func (ctx *Context) TimeStatisStop(topic string) {
	ctx.lock.RLock()
	defer ctx.lock.RUnlock()
	tmp := ctx.CurRecord().timeStatis[topic]
	if tmp == nil {
		return
	}
	tmp.StopPoint = time.Now()
}

// FlushLog 将日志写入磁盘
func (ctx *Context) FlushLog() {
	records := ctx.invokeRecords
	if len(records) == 0 {
		records = defaultRecords
	}
	for _, r := range records {
		kvs := [][2]string{}
		for _, f := range statisItems {
			kvs = append(kvs, [2]string{f, handlers[f](ctx, r)})
		}
		log.GetCalWorkerLogger().Notice(kvs)
		if r.Error != nil {
			log.GetCalWorkerLogger().Warn(kvs)
		}
	}

}

// Err2ErrorHandler 错误转换为错误码
// protocol 请求协议 当前有 http, nshead, pbrpc, mysql, redis
type Err2ErrorHandler func(protocol string, errMsg string) (errno string, dealSucc bool)

const (
	ErrnoHTTPEmptyBody               = "700"
	ErrnoHTTPAwaitingHeadersExceeded = "701"
	ErrnoHTTPIOTimeout               = "702"
	ErrnoUnKnown                     = "999"
)

// Err2ErrorHandlers 错误转换处理者
var Err2ErrorHandlers = []Err2ErrorHandler{
	// http: ContentLength=562 with Body length 0
	func(protocol string, errMsg string) (errno string, dealSucc bool) {
		if protocol != "http" {
			return
		}

		if strings.HasSuffix(errMsg, "with Body length 0") {
			return ErrnoHTTPEmptyBody, true
		}

		return
	},

	// net/http: request canceled (Client.Timeout exceeded while awaiting headers)
	func(protocol string, errMsg string) (errno string, dealSucc bool) {
		if protocol != "http" {
			return
		}

		if strings.HasSuffix(errMsg, "net/http: request canceled (Client.Timeout exceeded while awaiting headers)") {
			return ErrnoHTTPAwaitingHeadersExceeded, true
		}

		return
	},

	// dial tcp 10.26.7.174:8000: i/o timeout
	func(protocol string, errMsg string) (errno string, dealSucc bool) {
		if protocol != "http" {
			return
		}

		if strings.HasSuffix(errMsg, "i/o timeout") {
			return ErrnoHTTPEmptyBody, true
		}

		return
	},
}

var defaultRecords = []*InvokeRecord{&InvokeRecord{
	lock: new(sync.RWMutex),
}}

// InvokeRecord 访问日志，因为重试可能有多条
type InvokeRecord struct {
	// RspCode 请求的响应码
	// http 代表 http status code，200 为正常，700+是自定义的错误码，表示发送请求时发生了error
	// nshead 等有自己的规则，不统一描述
	RspCode int

	// Path 请求的路径
	// http 相对path， 形如： /foo/bar
	Path string

	// IDC 访问的IDC
	IDC string

	// IPPort ip和端口号
	IPPort string

	// Host 域名，可能和IPPort 一致
	Host string

	// 一次请求最多一条错误日志
	Error error

	timeStatis map[string]*StatisItem
	timePoints map[string]time.Time
	index      int
	lock       *sync.RWMutex
}

// GetTimeStatis 获取一个统计项
func (invokeRecord *InvokeRecord) GetTimeStatis(topic string) string {
	invokeRecord.lock.RLock()
	defer invokeRecord.lock.RUnlock()
	tmp := invokeRecord.timeStatis[topic]
	if tmp == nil {
		return "0"
	}
	return tmp.GetSpan()
}

// RecordTimePoint 打下一个时间点
func (invokeRecord *InvokeRecord) RecordTimePoint(topic string) {
	if _, ok := invokeRecord.timePoints[topic]; ok {
		return
	}
	invokeRecord.timePoints[topic] = time.Now()
}

// GetTimePoint 得到一个时间点 毫秒
func (invokeRecord *InvokeRecord) GetTimePoint(topic string) string {
	t := invokeRecord.timePoints[topic]
	if t.IsZero() {
		return "0"
	}

	return strconv.FormatInt(t.UnixNano()/1000000, 10)
}

var handlers = map[string]func(ctx *Context, invokeRecord *InvokeRecord) string{
	"appname": func(ctx *Context, invokeRecord *InvokeRecord) string { //自身模块名
		return env.AppName()
	},
	"uri": func(ctx *Context, invokeRecord *InvokeRecord) string { // 此cal请求的调用uri
		return invokeRecord.Path
	},
	"service": func(ctx *Context, invokeRecord *InvokeRecord) string { // 请求下游服务名字（服务名）
		return ctx.ServiceName
	},
	"req_len": func(ctx *Context, invokeRecord *InvokeRecord) string { // 打包后数据长度，即网络上发送的协议body size(kb)
		return strconv.FormatInt(ctx.ReqLen, 10)
	},
	"res_len": func(ctx *Context, invokeRecord *InvokeRecord) string { // 解包前的数据长度，即网络上接受到数据的body size(kb)
		return strconv.FormatInt(ctx.RspLen, 10)
	},
	"errno": func(ctx *Context, invokeRecord *InvokeRecord) string {
		if invokeRecord.Error == nil {
			return strconv.Itoa(invokeRecord.RspCode)
		}
		errMsg := invokeRecord.Error.Error()
		for _, handler := range Err2ErrorHandlers {
			if no, ok := handler(ctx.Protocol, errMsg); ok {
				return no
			}
		}

		return ErrnoUnKnown
	},
	"retry": func(ctx *Context, invokeRecord *InvokeRecord) string { //retry[0/2] 第1次交互(0表示未开始重试)，最多2次重试
		return fmt.Sprintf("%d/%d", invokeRecord.index, ctx.MaxTry)
	},
	"cost": func(ctx *Context, invokeRecord *InvokeRecord) string { // 总耗时,ms
		return invokeRecord.GetTimeStatis("cost")
	},
	"api": func(ctx *Context, invokeRecord *InvokeRecord) string { // 此cal请求的调用api
		return strings.Replace(strings.TrimLeft(strings.SplitN(invokeRecord.Path, "?", 1)[0], "/"), "/", "_", -1)
	},
	"logid": func(ctx *Context, invokeRecord *InvokeRecord) string { // 日志logid，需要去掉前面的0
		return strings.TrimLeft(fmt.Sprintf("%v", ctx.LogID), "0")
	},
	"caller": func(ctx *Context, invokeRecord *InvokeRecord) string { // 打印该条日志的对象，标识该条日志属于RAL或其他网络交互库
		return ctx.Caller
	},
	"method": func(ctx *Context, invokeRecord *InvokeRecord) string { // 协议的请求类型，对http包括get、post、put和delete
		return ctx.Method
	},
	"protocol": func(ctx *Context, invokeRecord *InvokeRecord) string { // 协议类型，如prot=http
		return ctx.Protocol
	},
	"balance": func(ctx *Context, invokeRecord *InvokeRecord) string { // 协议类型，如prot=http
		return ctx.BalanceName
	},
	"user_ip": func(ctx *Context, invokeRecord *InvokeRecord) string { // 	本次请求实际用户的ip
		if ctx.ReqContext != nil {
			return ctx.ReqContext.GetClientIP()
		}
		return ""
	},
	"idc": func(ctx *Context, invokeRecord *InvokeRecord) string { // 本机IDC
		return env.IDC()
	},
	"local_ip": func(ctx *Context, invokeRecord *InvokeRecord) string { // 本机ip
		return env.LocalIP()
	},
	"remote_ip": func(ctx *Context, invokeRecord *InvokeRecord) string { // 本次请求后端服务的ip:port
		return invokeRecord.IPPort
	},
	"remote_idc": func(ctx *Context, invokeRecord *InvokeRecord) string { // 本机IDC
		return invokeRecord.IDC
	},
	"remote_host": func(ctx *Context, invokeRecord *InvokeRecord) string { // 请求后端服务的域名
		return invokeRecord.Host
	},
	"uniqid": func(ctx *Context, invokeRecord *InvokeRecord) string { // 每次cal调用的唯一uniqid(失败重试的时候会改变)
		return fmt.Sprintf("%v%d", ctx.LogID, invokeRecord.index)
	},
	"talk": func(ctx *Context, invokeRecord *InvokeRecord) string { // 本次交互耗时，ms
		return invokeRecord.GetTimeStatis("talk")
	},
	"connect": func(ctx *Context, invokeRecord *InvokeRecord) string { // 链接耗时，ms
		return invokeRecord.GetTimeStatis("connect")
	},
	"write": func(ctx *Context, invokeRecord *InvokeRecord) string { // 写耗时ms
		return invokeRecord.GetTimeStatis("write")
	},
	"read": func(ctx *Context, invokeRecord *InvokeRecord) string { // 读耗时
		return invokeRecord.GetTimeStatis("read")
	},
	"pack": func(ctx *Context, invokeRecord *InvokeRecord) string { // 打包的耗时 ms
		return ctx.PackStatis.GetSpan()
	},
	"unpack": func(ctx *Context, invokeRecord *InvokeRecord) string { // 解包的耗时 ms
		return invokeRecord.GetTimeStatis("unpack")
	},
	"req_start_time": func(ctx *Context, invokeRecord *InvokeRecord) string { // 请求开始执行时间(打包之前) ms
		return invokeRecord.GetTimePoint("req_start_time")
	},
	"talk_start_time": func(ctx *Context, invokeRecord *InvokeRecord) string { // 请求准备连接时间(打包和负载均衡之后) ms
		return invokeRecord.GetTimePoint("talk_start_time")
	},
	"errmsg": func(ctx *Context, invokeRecord *InvokeRecord) string { // cal的错误码对应的错误信息
		if invokeRecord.Error == nil {
			return ""
		}
		return invokeRecord.Error.Error()
	},
}

var statisItems = []string{"appname", "uri", "service", "req_len", "res_len", "errno",
	"retry", "cost", "api", "logid", "caller", "method", "protocol", "balance",
	"user_ip", "local_ip", "idc", "remote_ip", "remote_idc", "remote_host",
	"uniqid", "talk", "connect", "write", "read", "pack", "unpack",
	"req_start_time", "talk_start_time", "errmsg",
}
