package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 注册路由
func RouterInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob("html/*")

	ginRouter.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusFound, "/login")
	})

	return ginRouter
}
