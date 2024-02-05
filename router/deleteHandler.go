package router

import (
	"net/http"
	"schedule/logs"
	"schedule/mysql"

	"github.com/gin-gonic/gin"
)

type Delete struct {
	DeleteType string `json:"deleteType"`
	DeleteId   int    `json:"deleteId"`
}

func DeleteHandler(ctx *gin.Context) {
	var delete Delete
	if err := ctx.ShouldBindJSON(&delete); err != nil {
		logs.GetInstance().Logger.Errorf("bind delete json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	switch delete.DeleteType {
	case "student":
		deleteStudentHandler(ctx, delete.DeleteId)
	case "class":
		deleteClassHandler(ctx, delete.DeleteId)
	case "teacher":
		deleteTeacherHandler(ctx, delete.DeleteId)
	case "lesson":
		deleteLessonHandler(ctx, delete.DeleteId)
	}
}

func deleteStudentHandler(ctx *gin.Context, id int) {
	db := mysql.GetClient()
	err := db.Where("studentId = ?", id).Delete(&mysql.StudentSql{})
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql delete student error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	showStudentsHandler(ctx)
}

func deleteClassHandler(ctx *gin.Context, id int) {
	db := mysql.GetClient()
	err := db.Where("classId = ?", id).Delete(&mysql.ClassSql{})
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql delete class error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	showClassesHandler(ctx)
}

func deleteTeacherHandler(ctx *gin.Context, id int) {
	db := mysql.GetClient()
	err := db.Where("teacherId = ?", id).Delete(&mysql.TeacherSql{})
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql delete teacher error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	showTeachersHandler(ctx)
}

func deleteLessonHandler(ctx *gin.Context, id int) {
	db := mysql.GetClient()
	err := db.Where("lessonId = ?", id).Delete(&mysql.LessonSql{})
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql delete student error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	showLessonsHandler(ctx)
}
