package router

import (
	v1 "gin_class/api/v1"
	"github.com/gin-gonic/gin"
)

func InitCoapRouter(Router *gin.RouterGroup) (R gin.IRoutes) {

	BaseRouter := Router.Group("coap")
	{
		BaseRouter.POST("send", v1.Send)
	}
	return BaseRouter
}
