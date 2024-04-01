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
		ctx.Redirect(http.StatusFound, "/home")
	})
	ginRouter.GET("/home", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", gin.H{})
	})
	ginRouter.GET("/student", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "student.html", gin.H{})
	})
	ginRouter.GET("/teacher", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "teacher.html", gin.H{})
	})
	ginRouter.GET("/lesson", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "lesson.html", gin.H{})
	})
	ginRouter.GET("/class", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "class.html", gin.H{})
	})
	ginRouter.GET("/schedule", ScheduleHandler)

	ginRouter.POST("/show", ShowHandler)
	ginRouter.POST("/update", UpdateHandler)
	ginRouter.POST("/add", AddHandler)
	ginRouter.POST("/delete", DeleteHandler)

	return ginRouter
}
