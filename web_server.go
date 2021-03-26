package laya

import (
	"github.com/gin-gonic/gin"
)

// WebServer 基于http协议的服务
type WebServer struct {
	*gin.Engine
}

// NewWebServer 创建WebServer
func NewWebServer(mode string) *WebServer {
	gin.SetMode(mode)
	server := &WebServer{
		Engine: gin.New(),
	}
	return server
}

// RouterRegister 路由注册者
type RouterRegister func(*WebServer)

// RegisterRouter 注册路由
func (webServer *WebServer) RegisterRouter(rr RouterRegister) {
	rr(webServer)
}
