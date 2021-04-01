package initialize

import (
	"fmt"
	_ "gin_class/docs"
	"gin_class/middleware"
	"gin_class/router"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	//"net/http"
)

func Routers() *gin.Engine {
	var Router = gin.Default()
	// 初始化 跨域服务
	Router.Use(middleware.Cors())
	// 初始化 swagger
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // 使用 Swagger

	fmt.Println("register swagger handler")

	PublicGroup := Router.Group("")
	{
		router.InitBaseRouter(PublicGroup)
		router.InitCoapRouter(PublicGroup)
		router.InitWebsocketRouter(PublicGroup)
	}
	PrivateGroup := Router.Group("")

	PrivateGroup.Use(middleware.JWTAuth())
	{

	}
	//global.GVA_LOG.Info("router register success")
	return Router
}
