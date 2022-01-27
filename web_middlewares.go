package laya

import (
	"bytes"
	"encoding/json"
	"github.com/layasugar/laya/genv"
	"github.com/layasugar/laya/glogs"
	"github.com/layasugar/laya/gutils"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"strings"
)

// 不需要打印入参和出参的路由
// 不需要打印入参和出参的前缀
// 不需要打印入参和出参的后缀
type logParams struct {
	NoLogParams       map[string]string
	NoLogParamsPrefix []string
	NoLogParamsSuffix []string
}

// 不想打印的路由分组
var noLogParamsRules logParams

// CheckNoLogParams 判断是否需要打印入参出参日志, 不需要打印返回true
func CheckNoLogParams(origin string) bool {
	if len(noLogParamsRules.NoLogParams) > 0 {
		if _, ok := noLogParamsRules.NoLogParams[origin]; ok {
			return true
		}
	}

	if len(noLogParamsRules.NoLogParamsPrefix) > 0 {
		for _, v := range noLogParamsRules.NoLogParamsPrefix {
			if strings.HasPrefix(origin, v) {
				return true
			}
		}
	}

	if len(noLogParamsRules.NoLogParamsSuffix) > 0 {
		for _, v := range noLogParamsRules.NoLogParamsSuffix {
			if strings.HasSuffix(origin, v) {
				return true
			}
		}
	}

	return false
}

// LogParams 记录框架出入参
func LogParams(ctx *WebContext) {
	w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: ctx.Writer}
	ctx.Writer = w

	if genv.ParamLog() {
		if !CheckNoLogParams(ctx.Request.RequestURI) {
			requestData, _ := ctx.GetRawData()
			ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData))
			ct := ctx.GetHeader("Content-Type")
			sct := strings.Split(ct, ";")
			switch sct[0] {
			case "application/json":
				var in map[string]interface{}
				_ = json.NewDecoder(bytes.NewBuffer(requestData)).Decode(&in)
				inJson, _ := json.Marshal(&in)
				ctx.InfoF("%s", string(inJson), glogs.String("header", gutils.GetString(ctx.Request.Header)), ctx.Field("title", "入参"))
			case "application/x-www-form-urlencoded", "multipart/form-data":
				ctx.InfoF("%s", string(requestData), glogs.String("header", gutils.GetString(ctx.Request.Header)), ctx.Field("title", "入参"))
			default:
				ctx.InfoF("%s", string(requestData), glogs.String("header", gutils.GetString(ctx.Request.Header)), ctx.Field("title", "入参"))
			}
		}
	}

	ctx.Next()
	if genv.ParamLog() {
		if !CheckNoLogParams(ctx.Request.RequestURI) {
			ctx.InfoF("%s", w.body.String(), ctx.Field("title", "出参"))
		}
	}
}

// SetTrace 注入链路追踪 优先级 request_id > trace_id
func SetTrace(ctx *WebContext) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	requestID := ctx.GetHeader(glogs.RequestIDName)
	if requestID == "" {
		ctx.Request.Header.Set(glogs.RequestIDName, gutils.Md5(uuid.NewV4().String()))
	}
	ctx.Next()
	//span := glogs.StartSpanR(ctx.Request, ctx.Request.RequestURI)
	//if span != nil {
	//	span.Tag(glogs.RequestIDName, ctx.GetHeader(glogs.RequestIDName))
	//	_ = glogs.Inject(context.WithValue(context.Background(), glogs.GetSpanContextKey(), span.Context()), ctx.Request)
	//	ctx.Next()
	//	span.Finish()
	//}
}
