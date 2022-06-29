package main

import (
	"suitbim.com/go-media-admin/core"
	"suitbim.com/go-media-admin/router"
)

// @title 这里写标题
// @version 1.0
// @description 这里写描述信息
// @termsOfService http://swagger.io/terms/

// @contact.name 这里写联系人信息
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host http://localhost:8005
// @BasePath /
func main() {
	// 启动流媒体服务,用来处理流请求
	dispatcher := core.NewDispatcher()
	go dispatcher.Start()
	// 启动后台接口服务,提供给业务系统使用
	router.InitRouterGroup(dispatcher)
}
