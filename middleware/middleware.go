package middleware

import "github.com/gin-gonic/gin"

// Default validation middleware and return middleware and signature middleware
type Middleware struct {
	*gin.HandlerFunc
}
