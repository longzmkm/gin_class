package router

import (
	v1 "gin_class/api/v1"
	"github.com/gin-gonic/gin"
)

func InitWebsocketRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	BaseRouter := Router.Group("ws")
	//Router.GET("/ws",v1.Stack)
	{
		BaseRouter.GET("ping", v1.Stack)
	}
	return Router
}
