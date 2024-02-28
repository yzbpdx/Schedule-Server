package router

import (
	"encoding/json"
	"net/http"
	"schedule/logs"
	"schedule/mysql"

	"github.com/gin-gonic/gin"
)

type ShowType struct {
	ShowType string `json:"showType"`
}

// 处理展示
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

// 处理展示学生
func showStudentsHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.StudentSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}

	studentResults := make([]StudentResult, len(results))
	for i, result := range results {
		studentResults[i].StudentId = result.StudentId
		studentResults[i].StudentName = result.StudentName
		studentResults[i].Status = result.Status
		var spareTime map[int]map[int]struct{}
		if err := json.Unmarshal(result.SpareTime, &spareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}
		for day, durations := range spareTime {
			for duration := range durations {
				studentResults[i].SpareTime[day][duration + 1] = true
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"results": studentResults})
}

// 处理展示班级
func showClassesHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.ClassSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query class error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}
	classResults := make([]ClassResult, len(results))
	for i, result := range results {
		classResults[i].ClassId = result.ClassId
		classResults[i].ClassName = result.ClassName
		classResults[i].Status = result.Status
		classMates := make(map[string]struct{})
		if err := json.Unmarshal(result.ClassMates, &classMates); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}
		var classMatesString string
		for mate := range classMates {
			classMatesString += mate + " "
		}
		classResults[i].ClassMates = classMatesString
	}	

	ctx.JSON(http.StatusOK, gin.H{"results": classResults})
}

// 处理展示老师
func showTeachersHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.TeacherSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query teacher error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}
	teacherResults := make([]TeacherResult, len(results))
	for i, result := range results {
		teacherResults[i].TeacherId = result.TeacherId
		teacherResults[i].TeacherName = result.TeacherName
		teacherResults[i].HolidayNum = result.HolidayNum
		teacherResults[i].Status = result.Status
		var spareTime map[int]map[int]struct{}
		if err := json.Unmarshal(result.SpareTime, &spareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}
		for day, durations := range spareTime {
			for duration := range durations {
				teacherResults[i].SpareTime[day][duration + 1] = true
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"results": teacherResults})
}

// 处理展示课程
func showLessonsHandler(ctx *gin.Context) {
	db := mysql.GetClient()
	results := make([]mysql.LessonSql, 0)
	if err := db.Find(&results); err.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query class error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
	}
	// lessonResults := make([]LessonResult, len(results))
	// for i, result := range results {
	// 	lessonResults[i].LessonId = result.LessonId
	// 	lessonResults[i].LessonName = result.LessonName
	// 	lessonResults[i].TeacherName = result.TeacherName
	// 	lessonResults[i].StudyName = result.StudyName
	// 	lessonResults[i].StudentNum = result.StudentNum
	// }

	ctx.JSON(http.StatusOK, gin.H{"results": results})
}
