package router

import (
	"encoding/json"
	"net/http"
	"schedule/logs"
	"schedule/mysql"
	"strings"

	"github.com/gin-gonic/gin"
)

type Add struct {
	AddType     string        `json:"addType"`
	StudentData StudentResult `json:"studentData"`
	TeacherData TeacherResult `json:"teacherData"`
	ClassData   ClassResult   `json:"classData"`
	LessonData  LessonResult  `json:"lessonData"`
}

// 处理新建信息
func AddHandler(ctx *gin.Context) {
	var add Add
	if err := ctx.ShouldBindJSON(&add); err != nil {
		logs.GetInstance().Logger.Errorf("bind create json error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	switch add.AddType {
	case "student":
		addStudentHandler(&add, ctx)
	case "class":
		addClassHandler(&add, ctx)
	case "teacher":
		addTeacherHandler(&add, ctx)
	case "lesson":
		addLessonHandler(&add, ctx)
	}
}

// 处理新建学生信息
func addStudentHandler(add *Add, ctx *gin.Context) {
	var student mysql.StudentSql
	student.StudentName = add.StudentData.StudentName
	student.Status = add.StudentData.Status
	spareTime := make(map[int]map[int]struct{})
	ConvertSpareTime(spareTime, add.StudentData.SpareTime)
	spareTimeJson, err := json.Marshal(spareTime)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal spare time error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	student.SpareTime = spareTimeJson

	db := mysql.GetClient()
	e := db.Create(&student)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create student sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}

// 处理新建课程信息
func addClassHandler(add *Add, ctx *gin.Context) {
	var class mysql.ClassSql
	class.ClassName = add.ClassData.ClassName
	class.Status = add.ClassData.Status
	classMates := strings.Split(add.ClassData.ClassMates, " ")
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
	e := db.Create(&class)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create class sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}

// 处理新建老师信息
func addTeacherHandler(add *Add, ctx *gin.Context) {
	var teacher mysql.TeacherSql
	teacher.TeacherName = add.TeacherData.TeacherName
	teacher.Status = add.TeacherData.Status
	teacher.HolidayNum = add.TeacherData.HolidayNum
	spareTime := make(map[int]map[int]struct{})
	ConvertSpareTime(spareTime, add.TeacherData.SpareTime)
	spareTimeJson, err := json.Marshal(spareTime)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal spare time error %s", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	teacher.SpareTime = spareTimeJson

	db := mysql.GetClient()
	e := db.Create(&teacher)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create teacher sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}

// 处理新建课程信息
func addLessonHandler(add *Add, ctx *gin.Context) {
	var lesson mysql.LessonSql
	lesson.LessonName = add.LessonData.LessonName
	lesson.StudyName = add.LessonData.StudyName
	lesson.TeacherName = add.LessonData.TeacherName
	lesson.StudentNum = add.LessonData.StudentNum

	db := mysql.GetClient()
	e := db.Create(&lesson)
	if e.Error != nil {
		logs.GetInstance().Logger.Errorf("create lesson sql error %s", e.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
}
