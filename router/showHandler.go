package router

import (
	"net/http"
	"schedule/logs"
	"schedule/mysql"

	"github.com/gin-gonic/gin"
)

type ShowType struct {
	ShowType string `json:"showType"`
}

func ShowHandler(ctx *gin.Context) {
	var showType ShowType
	if err := ctx.ShouldBindJSON(&showType); err != nil {
		logs.GetInstance().Logger.Errorf("bind show type json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	switch showType.ShowType {
	case "student":
		showStudentsHandler(ctx)
	case "class":
		showClassesHandler(ctx)
	case "teacher":
		showTeachersHandler(ctx)
	case "lesson":
		showLessonsHandler(ctx)
	}
}

func showStudentsHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.StudentSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}

func showClassesHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.ClassSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query class error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}

func showTeachersHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.TeacherSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query teacher error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}

func showLessonsHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.LessonSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query class error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}
