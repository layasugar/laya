package laya

import (
	log "github.com/LaYa-op/laya/logger"
	"github.com/LaYa-op/laya/store/db"
	"github.com/LaYa-op/laya/store/redis"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/web"
)

func Init() {
	db.Init()
	redis.Init()
	log.Init()
}

// 注册路由,注册全局中间件
func NewServer() *gin.Engine {
	// create new web service
	service := web.NewService(
		web.Name("tb.controllers.hall"),
		web.Version("1.0"),
		web.Address(":8080"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.ZapLog.Error(err)
	}

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	//r.Use(middleware.Sign(), middleware.Response())
	//r.Use(middleware.Base.Auth)
	//r.Use(middleware.Base.Response)
	//r.Use(middleware.Base.Sign)
	//r.Use(middleware.Base.Test)
	service.Handle("/", r)

	// initialise route
	routers.Init(r)

	return r
	// run service
	//if err := service.Run(); err != nil {
	//	log.ZapLog.Error(err)
	//}
}

func Run()  {
	if err := service.Run(); err != nil {
		log.ZapLog.Error(err)
	}

}