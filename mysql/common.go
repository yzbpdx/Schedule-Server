package mysql

import "schedule/common"

type StudentSql struct {
	StudentId   int    `json:"studentId" gorm:"column:studentId"`
	StudentName string `json:"studentName" gorm:"column:studentName"`
	SpareTime   []byte `json:"spareTime" gorm:"column:spareTime"`
	Status      bool   `json:"status" gorm:"column:status"`
}

func (s StudentSql) TableName() string {
	return "students"
}

type TeacherSql struct {
	TeacherId   int    `json:"teacherId" gorm:"column:teacherId"`
	TeacherName string `json:"teacherName" gorm:"column:teacherName"`
	SpareTime   []byte `json:"spareTime" gorm:"column:spareTime"`
	HolidayNum  int    `json:"holidayNum" gorm:"column:holidayNum"`
	Status      bool   `json:"status" gorm:"column:status"`
}

func (t TeacherSql) TableName() string {
	return "teachers"
}

type ClassSql struct {
	ClassId    int    `json:"classId" gorm:"column:classId"`
	ClassName  string `json:"className" gorm:"column:className"`
	ClassMates []byte `json:"classMates" gorm:"column:classMates"`
	Status     bool   `json:"status" gorm:"column:status"`
}

func (c ClassSql) TableName() string {
	return "classes"
}

type LessonSql struct {
	common.LessonDict
}

func (l LessonSql) TableName() string {
	return "lessons"
}
