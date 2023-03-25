package middleware

import (
	"github.com/kataras/iris/v12"
)


func AuthToken(ctx iris.Context){
	//token := ctx.URLParamDefault("token", "notoken")
	//logrus.Info("Token: "+token)
	ctx.Next()
}