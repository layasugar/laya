package laya

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/layasugar/laya/gcnf"
)

type HttpResp struct{}

type Response struct {
	StatusCode uint32      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	RequestID  string      `json:"request_id"`
}

type Pagination struct {
	Total  int64 `json:"total"`
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
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
	msg := gcnf.LoadErrMsg(re.Code)
	if msg == "" {
		msg = "sorry, system err"
	}
	re.Msg = msg
	return re.Code, re.Msg
}

func (res *HttpResp) Suc(ctx *gin.Context, data interface{}, msg ...string) {
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

func (res *HttpResp) Fail(ctx *gin.Context, err error) {
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
func (res *HttpResp) RawJSONString(ctx *gin.Context, data string) {
	w := ctx.Writer
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(data))
}

// RawString raw 数据返回
func (res *HttpResp) RawString(ctx *gin.Context, data string) {
	w := ctx.Writer
	w.WriteHeader(200)
	_, _ = w.Write([]byte(data))
}
