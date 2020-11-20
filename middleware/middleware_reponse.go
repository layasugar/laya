package middleware

import (
	"github.com/LaYa-op/laya/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (*Middleware) Response(c *gin.Context) {
	c.Next()
	if c.Writer.Written() {
		return
	}

	params := c.Keys
	if len(params) == 0 {
		return
	}

	al := c.GetHeader("Accept-Language")
	var r response.Response
	c.JSON(http.StatusOK, r.GetResponse(params, al))
}
