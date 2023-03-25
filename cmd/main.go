package main

import (
	"SuperArch/conf"
	"SuperArch/middleware/process"
	"SuperArch/middleware/register"
	"SuperArch/routers"
	"github.com/kataras/iris/v12"
	"os"
)

func NewApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug") //设置日志级别
	routers.Register(app) // 注册路由

	return app
}

func main(){
	conf.InitConfig()
	register.SchedulerRegister.Init()

	process.CronProcess()

	app := NewApp()
	var addr string
	if len(os.Getenv("DEBUG")) > 0{
		addr = ":7777"
	}else{
		addr = ":80"
	}
	app.Run(iris.Addr(addr), iris.WithoutServerError(iris.ErrServerClosed))
}


