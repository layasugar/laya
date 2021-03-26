package gmiddleware

import (
	"github.com/gin-gonic/gin"
	"github.com/layatips/laya/gresp"
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
	var r gresp.Response
	c.JSON(http.StatusOK, r.GetResponse(params, al))
}
