package routers

import (
	"SuperArch/core"
	"SuperArch/middleware"
	"github.com/kataras/iris/v12"
)

func Register(app *iris.Application) {
	root := app.Party("/", middleware.AuthToken).AllowMethods(iris.MethodGet)
	{
		javaUtil := root.Party("/task")
		{
			javaUtil.Post("/add", core.Task{}.AddTask)
			javaUtil.Post("/status", core.Task{}.TaskStatus)
		}
	}
}