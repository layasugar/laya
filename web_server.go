package laya

import (
	"github.com/gin-gonic/gin"
)

// WebServer 基于http协议的服务
// 这里的实现是基于gin框架，封装了gin的所有的方法
// gin 的核心是高效路由，但是gin.Engine和gin.IRouter(s)的高耦合让我们无法复用，gin的作者认为它的路由就是引擎吧
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
