package middleware

import (
	"github.com/LaYa-op/laya"
	"github.com/LaYa-op/laya/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Default validation middleware and return middleware and signature middleware
type Middleware struct {

}

func (*Middleware) Autha() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		uid, err := laya.Redis.Get("user:token:" + token).Result()
		if err != nil {
			c.Set("$.TokenErr.code", response.TokenErr)
			c.Abort()
			return
		}

		ID, _ := strconv.ParseInt(uid, 10, 64)
		c.Set("uid", ID)
		c.Next()
	}
}
