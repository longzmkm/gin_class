package core

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"time"
)
// 可以优雅的 零停机时间重新启动
func initServer(address string, router *gin.Engine) server {
	s := endless.NewServer(address, router)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20
	return s
}
