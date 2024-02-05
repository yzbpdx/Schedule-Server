package router

import (
	"net/http"
	"schedule/logs"
	"schedule/mysql"

	"github.com/gin-gonic/gin"
)

type Update struct {
	UpdateType string `json:"updateType"`
}

func UpdateHandler(ctx *gin.Context) {
	var update Update
	if err := ctx.ShouldBindJSON(&update); err != nil {
		logs.GetInstance().Logger.Errorf("bind update type json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	switch update.UpdateType {
	case "student":
		updateStudentHandler(ctx)
	case "class":
		updateClassHandler(ctx)
	case "teacher":
		updateTeacherHandler(ctx)
	case "lesson":
		updateLessonHandler(ctx)
	}
}

func updateStudentHandler(ctx *gin.Context) {
	var student mysql.StudentSql
	if err := ctx.ShouldBindJSON(&student); err != nil {
		logs.GetInstance().Logger.Errorf("bind student json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Save(&student)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create student sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showStudentsHandler(ctx)
}

func updateClassHandler(ctx *gin.Context) {
	var class mysql.ClassSql
	if err := ctx.ShouldBindJSON(&class); err != nil {
		logs.GetInstance().Logger.Errorf("bind class json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Save(&class)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create class sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showClassesHandler(ctx)
}

func updateTeacherHandler(ctx *gin.Context) {
	var teacher mysql.TeacherSql
	if err := ctx.ShouldBindJSON(&teacher); err != nil {
		logs.GetInstance().Logger.Errorf("bind teacher json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Save(&teacher)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create teacher sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showTeachersHandler(ctx)
}

func updateLessonHandler(ctx *gin.Context) {
	var lesson mysql.LessonSql
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		logs.GetInstance().Logger.Errorf("bind lesson json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Save(&lesson)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create lesson sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showLessonsHandler(ctx)
}
