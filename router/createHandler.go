package router

import (
	"net/http"
	"schedule/logs"
	"schedule/mysql"

	"github.com/gin-gonic/gin"
)

type Create struct {
	CreateType string `json:"createType"`
}

// 处理新建信息
func CreateHandler(ctx *gin.Context) {
	var create Create
	if err := ctx.ShouldBindJSON(&create); err != nil {
		logs.GetInstance().Logger.Errorf("bind create json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	switch create.CreateType {
	case "student":
		createStudentHandler(ctx)
	case "class":
		createClassHandler(ctx)
	case "teacher":
		createTeacherHandler(ctx)
	case "lesson":
		createLessonHandler(ctx)
	}
}

// 处理新建学生信息
func createStudentHandler(ctx *gin.Context) {
	var student mysql.StudentSql
	if err := ctx.ShouldBindJSON(&student); err != nil {
		logs.GetInstance().Logger.Errorf("bind student json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Create(&student)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create student sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showStudentsHandler(ctx)
}

// 处理新建课程信息
func createClassHandler(ctx *gin.Context) {
	var class mysql.ClassSql
	if err := ctx.ShouldBindJSON(&class); err != nil {
		logs.GetInstance().Logger.Errorf("bind class json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Create(&class)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create class sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showClassesHandler(ctx)
}

// 处理新建老师信息
func createTeacherHandler(ctx *gin.Context) {
	var teacher mysql.TeacherSql
	if err := ctx.ShouldBindJSON(&teacher); err != nil {
		logs.GetInstance().Logger.Errorf("bind teacher json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Create(&teacher)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create teacher sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showTeachersHandler(ctx)
}

// 处理新建课程信息
func createLessonHandler(ctx *gin.Context) {
	var lesson mysql.LessonSql
	if err := ctx.ShouldBindJSON(&lesson); err != nil {
		logs.GetInstance().Logger.Errorf("bind lesson json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	db := mysql.GetClient()
	err := db.Create(&lesson)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create lesson sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	showLessonsHandler(ctx)
}
