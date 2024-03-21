package router

import (
	"encoding/json"
	"net/http"
	"schedule/logs"
	"schedule/mysql"
	"strings"

	"github.com/gin-gonic/gin"
)

type Update struct {
	UpdateType  string        `json:"updateType"`
	StudentData StudentResult `json:"studentData"`
	TeacherData TeacherResult `json:"teacherData"`
	ClassData   ClassResult   `json:"classData"`
	LessonData  LessonResult  `json:"lessonData"`
}

// 处理更新
func UpdateHandler(ctx *gin.Context) {
	var update Update
	if err := ctx.ShouldBindJSON(&update); err != nil {
		logs.GetInstance().Logger.Errorf("bind update type json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	switch update.UpdateType {
	case "student":
		updateStudentHandler(&update, ctx)
	case "class":
		updateClassHandler(&update, ctx)
	case "teacher":
		updateTeacherHandler(&update, ctx)
	case "lesson":
		updateLessonHandler(&update, ctx)
	}
}

// 处理更新学生
func updateStudentHandler(update *Update, ctx *gin.Context) {
	var student mysql.StudentSql
	student.StudentId = update.StudentData.StudentId
	student.StudentName = update.StudentData.StudentName
	student.Status = update.StudentData.Status
	spareTime := make(map[int]map[int]struct{})
	ConvertSpareTime(spareTime, update.StudentData.SpareTime)
	spareTimeJson, err := json.Marshal(spareTime)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal spare time error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	student.SpareTime = spareTimeJson

	db := mysql.GetClient()
	e := db.Where("studentId = ?", student.StudentId).Save(&student)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create student sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// showStudentsHandler(ctx)
}

// 处理更新课程
func updateClassHandler(update *Update, ctx *gin.Context) {
	var class mysql.ClassSql
	class.ClassId = update.ClassData.ClassId
	class.ClassName = update.ClassData.ClassName
	class.Status = update.ClassData.Status
	classMates := strings.Split(update.ClassData.ClassMates, " ")
	classMatesMap := make(map[string]struct{})
	for _, mate := range classMates {
		classMatesMap[mate] = struct{}{}
	}
	classMatesJson, err := json.Marshal(classMatesMap)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal class mates error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	class.ClassMates = classMatesJson

	db := mysql.GetClient()
	e := db.Where("classId = ?", class.ClassId).Save(&class)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create class sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// showClassesHandler(ctx)
}

// 处理更新老师
func updateTeacherHandler(update *Update, ctx *gin.Context) {
	var teacher mysql.TeacherSql
	teacher.TeacherId = update.TeacherData.TeacherId
	teacher.TeacherName = update.TeacherData.TeacherName
	teacher.Status = update.TeacherData.Status
	teacher.HolidayNum = update.TeacherData.HolidayNum
	spareTime := make(map[int]map[int]struct{})
	ConvertSpareTime(spareTime, update.TeacherData.SpareTime)
	spareTimeJson, err := json.Marshal(spareTime)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal spare time error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	teacher.SpareTime = spareTimeJson

	db := mysql.GetClient()
	e := db.Where("teacherId = ?", teacher.TeacherId).Save(&teacher)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create teacher sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// showTeachersHandler(ctx)
}

// 处理更新课程
func updateLessonHandler(update *Update, ctx *gin.Context) {
	var lesson mysql.LessonSql
	lesson.LessonId = update.LessonData.LessonId
	lesson.LessonName = update.LessonData.LessonName
	lesson.StudyName = update.LessonData.StudyName
	lesson.TeacherName = update.LessonData.TeacherName
	lesson.StudentNum = update.LessonData.StudentNum

	db := mysql.GetClient()
	err := db.Where("lessonId = ?", lesson.LessonId).Save(&lesson)
	if err.Error != nil {
		logs.GetInstance().Logger.Errorf("create lesson sql error %s", err.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	// showLessonsHandler(ctx)
}
