package router

import (
	"net/http"
	"schedule/algorithm"

	"github.com/gin-gonic/gin"
)

func ScheduleHandler(ctx *gin.Context) {
	algorithm.StartSchedule()
	ctx.JSON(http.StatusOK, gin.H{})
}
