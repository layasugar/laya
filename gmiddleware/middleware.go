package gmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/glogs"
	uuid "github.com/satori/go.uuid"
)

// implements the gin.handlerFunc
func SetHeader(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	requestID := c.GetHeader(glogs.RequestIDName)
	if requestID == "" {
		c.Request.Header.Set(glogs.RequestIDName, uuid.NewV4().String())
	}
	c.Next()
}
