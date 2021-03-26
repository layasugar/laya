package glogs

import (
	"context"
	"errors"
	"github.com/openzipkin/zipkin-go"
	"net/http"

	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

var Tracer *zipkin.Tracer //tracer 引擎

var spanContextKey string //ctx key，约定ctx的key名称

type Trace struct {
	ServiceName     string //服务名
	ServiceEndpoint string //当前服务节点
	ZipkinAddr      string //zipkin地址
	Mod             uint64 //采样率,0==不进行链路追踪，1==全量。值越大，采样率月底，对性能影响越小
}

//获取默认，用于demo
func GetDefaultTrace(zipkinAddr string) *Trace {
	return GetNewTrace("laya_go_template_trace", "localhost:80", zipkinAddr, "zipkin_span", 1)
}

//获取配置
func GetNewTrace(serviceName, serviceEndpoint, zipkinAddr, sContextKey string, mod uint64) *Trace {
	spanContextKey = sContextKey
	return &Trace{
		ServiceName:     serviceName,
		ServiceEndpoint: serviceEndpoint,
		ZipkinAddr:      zipkinAddr,
		Mod:             mod,
	}
}

//初始化tracer
func (t *Trace) InitTracer() error {
	var err error
	Tracer, err = t.GetTrace()
	return err
}

//获取tracer
func (t *Trace) GetTrace() (*zipkin.Tracer, error) {
	if t == (&Trace{}) {
		return nil, errors.New("trace is not init")
	}
	// create a reporter to be used by the tracer
	reporter := httpreporter.NewReporter(t.ZipkinAddr)
	// set-up the local endpoint for our service
	endpoint, err := zipkin.NewEndpoint(t.ServiceName, t.ServiceEndpoint)
	if err != nil {
		return nil, err
	}
	// set-up our sampling strategy
	sampler := zipkin.NewModuloSampler(t.Mod)
	if t.Mod == 0 {
		sampler = zipkin.NeverSample
	}
	// initialize the tracer
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	return tracer, err
}

//根据上下文创建span
func StartSpan(ctx context.Context, trace *zipkin.Tracer, name string) zipkin.Span {
	spanChild := trace.StartSpan(name)
	if ctx == nil {
		return spanChild
	}
	spanI := ctx.Value(spanContextKey)
	if spanContext, ok := spanI.(model.SpanContext); ok {
		spanChild = trace.StartSpan(name, zipkin.Parent(spanContext))
	}
	return spanChild
}

//根据请求头创建span
func StartSpanFromReq(r *http.Request, trace *zipkin.Tracer, name string) zipkin.Span {
	if r != (&http.Request{}) {
		return trace.StartSpan(name, zipkin.Parent(trace.Extract(b3.ExtractHTTP(copyRequest(r)))))
	}
	return trace.StartSpan(name)
}

//注入span信息到请求头
func InjectToReq(ctx context.Context, r *http.Request) error {
	injector := b3.InjectHTTP(r)
	if ctx == nil {
		return errors.New("ctx is nil")
	}
	spanI := ctx.Value(spanContextKey)
	if spanContext, ok := spanI.(model.SpanContext); ok {
		err := injector(spanContext)
		return err
	}
	return nil
}

func copyRequest(r *http.Request) *http.Request {
	req := &http.Request{}
	req = r
	return req
}
func GetSpanContextKey() string { return spanContextKey }
