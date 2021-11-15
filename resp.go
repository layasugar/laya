package laya

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/layasugar/glogs"
	"github.com/layasugar/laya/gconf"
	"github.com/layasugar/laya/genv"
	"net/http"
)

const requestIDName = glogs.RequestIDName

type Resp struct{}

type response struct {
	StatusCode uint32      `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	RequestID  string      `json:"request_id"`
}

type Pagination struct {
	Total       int64 `json:"total"`
	Count       int64 `json:"count"`
	PerPage     int64 `json:"per_page"`
	CurrentPage int64 `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}

// RspError 错误处理
type RspError struct {
	Code uint32
	Msg  string
}

func (re *RspError) Error() string {
	return re.Msg
}

func Err(code uint32) (err error) {
	err = &RspError{
		Code: code,
	}
	return err
}

// Render 渲染
func (re *RspError) render() (code uint32, msg string) {
	key := fmt.Sprintf("err_code.%d", re.Code)
	s := gconf.C.GetString(key)
	if s == "" {
		s = "sorry, system err"
	}
	re.Msg = s
	return re.Code, re.Msg
}

func (res *Resp) Suc(c *gin.Context, data interface{}, msg ...string) {
	rr := new(response)
	rr.StatusCode = http.StatusOK
	if len(msg) == 0 {
		rr.Message = "success"
	} else {
		for _, v := range msg {
			rr.Message += "," + v
		}
	}
	rr.Data = data
	rr.RequestID = c.GetHeader(requestIDName)
	if genv.ParamLog() {
		log, _ := json.Marshal(&rr)
		glogs.InfoF(c, "出参", string(log))
	}

	c.JSON(http.StatusOK, &rr)
}

func (res *Resp) Fail(c *gin.Context, err error) {
	rr := new(response)
	switch err.(type) {
	case *RspError:
		rr.StatusCode, rr.Message = err.(*RspError).render()
	default:
		rr.StatusCode = 400
		rr.Message = err.Error()
	}
	rr.RequestID = c.GetHeader(requestIDName)
	if genv.ParamLog() {
		log, _ := json.Marshal(&rr)
		glogs.InfoF(c, "出参", string(log))
	}

	c.JSON(http.StatusOK, &rr)
}

// RawJSONString json 数据返回
func (res *Resp) RawJSONString(c *gin.Context, data string) {
	w := c.Writer
	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write([]byte(data))
}

// RawString raw 数据返回
func (res *Resp) RawString(c *gin.Context, data string) {
	w := c.Writer
	w.WriteHeader(200)
	_, _ = w.Write([]byte(data))
}
