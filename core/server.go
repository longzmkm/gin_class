package core

import (
	"fmt"
	"gin_class/global"
	"gin_class/initialize"
	"go.uber.org/zap"
	"time"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")
	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	time.Sleep(10 * time.Microsecond)
	fmt.Println("server run success on address", address)
	//s.ListenAndServe().Error()
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Printf(`
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
`, address)

	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
