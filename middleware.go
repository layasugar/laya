package laya

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/layasugar/glogs"
	"github.com/layasugar/laya/genv"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"strings"
)

// SetHeader implements the gin.handlerFunc
func SetHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	requestID := c.GetHeader(glogs.RequestIDName)
	if requestID == "" {
		c.Request.Header.Set(glogs.RequestIDName, uuid.NewV4().String())
	}
	c.Next()
}

// LogInParams 记录框架出入参
func LogInParams(c *gin.Context) {
	if genv.ParamLog() {
		requestData, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestData))
		ct := c.GetHeader("Content-Type")
		sct := strings.Split(ct, ";")
		switch sct[0] {
		case "application/json":
			var in map[string]interface{}
			_ = json.NewDecoder(bytes.NewBuffer(requestData)).Decode(&in)
			inJson, _ := json.Marshal(&in)
			glogs.InfoF(c, "入参", string(inJson), glogs.String("header", c.Request.Header))
		case "application/x-www-form-urlencoded", "multipart/form-data":
			glogs.InfoF(c, "入参", string(requestData), glogs.String("header", c.Request.Header))
		default:
			glogs.InfoF(c, "入参", string(requestData), glogs.String("header", c.Request.Header))
		}
	}

	c.Next()
}
