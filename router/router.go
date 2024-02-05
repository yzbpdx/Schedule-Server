package router

import "github.com/gin-gonic/gin"

func RouterInit() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	ginRouter := gin.Default()
	ginRouter.LoadHTMLGlob("html/*")

	return ginRouter
}
