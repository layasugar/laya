package http_tpl

const ModelsDaoCalHttpTestTraceTpl = `// 请求测试文件

package http_test

import (
	"errors"
	"github.com/layasugar/laya"
	"github.com/layasugar/laya/gcal"
	"net/http"
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

// HttpToHttpTraceTest Http测试, body是interface可以发送任何类型的数据
func HttpToHttpTraceTest(ctx *laya.WebContext) (*Data, error) {
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
`
