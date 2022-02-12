package http_tpl

const ModelsDaoCalHttpTestServiceTestTpl = `// 链路追踪的测试文件

package cal_test

import (
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"io"
	"log"
	"net/http"
	"testing"
)

var RespSuc = []byte({{.goModName}}{
"data": {"code": "trace-http-test"},
"message": "操作成功",
"status_code": 200
}{{.goModName}})

var agentHost = "127.0.0.1:6831"

func TestStartHttp(t *testing.T) {
	tracer, closer, _ := NewJaeger("http_server")
	defer closer.Close()

	var serverMux = http.NewServeMux()
	serverMux.HandleFunc("/server-b/fast", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			log.Println(err.Error())
		}
		firstSpan := tracer.StartSpan(r.URL.Path, ext.RPCServerOption(spanCtx))
		ext.HTTPUrl.Set(firstSpan, r.URL.Path)
		ext.HTTPMethod.Set(firstSpan, r.Method)
		firstSpan.SetTag("is_debug", "1")
		firstSpan.Finish()
		w.Write(RespSuc)
	})

	http.ListenAndServe("0.0.0.0:10081", serverMux)
}

func NewJaeger(serverName string) (opentracing.Tracer, io.Closer, error) {
	var cfg = jaegerCfg.Configuration{
		ServiceName: serverName,

		Sampler: &jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentHost,
		},
	}
	jLogger := jaegerLog.StdLogger
	return cfg.NewTracer(
		jaegerCfg.Logger(jLogger),
	)
}
`
