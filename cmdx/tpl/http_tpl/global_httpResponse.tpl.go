package http_tpl

const GlobalHttpResponseTpl = `package global

import (
	"github.com/layasugar/laya"
	"github.com/layasugar/laya/gconf"
	"net/http"
)

type HttpResp struct{}

type Response struct {
	StatusCode uint32      {{.tagName}}json:"status_code"{{.tagName}}
	Message    string      {{.tagName}}json:"message"{{.tagName}}
	Data       interface{} {{.tagName}}json:"data"{{.tagName}}
	RequestID  string      {{.tagName}}json:"request_id"{{.tagName}}
}

type Pagination struct {
	Total       int64 {{.tagName}}json:"total"{{.tagName}}
	Count       int64 {{.tagName}}json:"count"{{.tagName}}
	PerPage     int64 {{.tagName}}json:"per_page"{{.tagName}}
	CurrentPage int64 {{.tagName}}json:"current_page"{{.tagName}}
	TotalPages  int64 {{.tagName}}json:"total_pages"{{.tagName}}
}

// rspError 错误处理
type rspError struct {
	Code uint32
	Msg  string
}

func (re *rspError) Error() string {
	return re.Msg
}

func Err(code uint32) (err error) {
	err = &rspError{
		Code: code,
	}
	return err
}

// Render 渲染
func (re *rspError) render() (uint32, string) {
	msg := gconf.LoadErrMsg(re.Code)
	if msg == "" {
		msg = "sorry, system err"
	}
	re.Msg = msg
	return re.Code, re.Msg
}

func (res *HttpResp) Suc(ctx *laya.WebContext, data interface{}, msg ...string) {
	rr := new(Response)
	rr.StatusCode = http.StatusOK
	if len(msg) == 0 {
		rr.Message = "success"
	} else {
		for _, v := range msg {
			rr.Message += "," + v
		}
	}
	rr.Data = data
	rr.RequestID = ctx.GetLogId()
	ctx.JSON(http.StatusOK, &rr)
}

func (res *HttpResp) Fail(ctx *laya.WebContext, err error) {
	rr := new(Response)
	switch err.(type) {
	case *rspError:
		rr.StatusCode, rr.Message = err.(*rspError).render()
	default:
		rr.StatusCode = 400
		rr.Message = err.Error()
	}
	rr.RequestID = ctx.GetLogId()
	ctx.JSON(http.StatusOK, &rr)
}

// RawJSONString json 数据返回
func (res *HttpResp) RawJSONString(ctx *laya.WebContext, data string) {
	w := ctx.Writer
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(data))
}

// RawString raw 数据返回
func (res *HttpResp) RawString(ctx *laya.WebContext, data string) {
	w := ctx.Writer
	w.WriteHeader(200)
	_, _ = w.Write([]byte(data))
}
`
