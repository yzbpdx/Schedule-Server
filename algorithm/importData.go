package algorithm

import (
	"encoding/json"
	"schedule/common"
	"schedule/logs"
	"schedule/mysql"

	"gorm.io/gorm"
)

type Student struct {
	Dict   common.StudentDict
	Extend StudentExtend
	Lessons map[string]Lesson
}

type StudentExtend struct {
	Priority  float64
	LessonNum int
	SpareNum  int
}

type Teacher struct {
	Dict   common.TeacherDict
	Extend TeacherExtend
	Lessons map[string]Lesson
}

type TeacherExtend struct {
	Priority  float64
	LessonNum int
	SpareNum  int
	DaysDistribution [7][3]int
	WorkDays         [7]WorkDay
	Holidays         map[int]struct{}
}

type WorkDay struct {
	Day int
	WorkNum int
}

type Class struct {
	Dict common.ClassDict
	Extend ClassExtend
	Lessons map[string]Lesson
}

type ClassExtend struct {
	Priority  float64
	LessonNum int
	SpareNum  int
	SpareTime map[int]map[int]struct{}
}

type Lesson struct {
	Dict common.LessonDict
	Extend LessonExtend
}

type LessonExtend struct {
	Priority float64
	SpareTime map[int]map[int]struct{}
	CandidateDays []CandidateDay
}

type CandidateDay struct {
	Day int
	Duration int
	Priority float64
}

// 从sql中加载学生信息
func ImportStudents(db *gorm.DB) map[string]Student {
	students := make(map[string]Student)
	rows := make([]mysql.StudentSql, 0)
	result := db.Find(&rows, "status = ?", true)
	if result.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", result.Error)
		return students
	}

	for _, row := range rows {
		var spareTime map[int]map[int]struct{}
		if err := json.Unmarshal(row.SpareTime, &spareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}

		students[row.StudentName] = Student{
			Dict: common.StudentDict{
				StudentId: row.StudentId,
				StudentName: row.StudentName,
				SpareTime: spareTime,
			},
			Extend: StudentExtend{},
			Lessons: make(map[string]Lesson),
		}
	}

	return students
}

// 从sql中加载老师信息
func ImportTeachers(db *gorm.DB) map[string]Teacher {
	teachers := make((map[string]Teacher))
	rows := make([]mysql.TeacherSql, 0)
	result := db.Find(&rows, "status = ?", true)
	if result.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", result.Error)
		return teachers
	}

	for _, row := range rows {
		var spareTime map[int]map[int]struct{}
		if err := json.Unmarshal(row.SpareTime, &spareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}

		teachers[row.TeacherName] = Teacher{
			Dict: common.TeacherDict{
				TeacherId: row.TeacherId,
				TeacherName: row.TeacherName,
				SpareTime: spareTime,
				HolidayNum: row.HolidayNum,
			},
			Extend: TeacherExtend{
				Holidays: make(map[int]struct{}),
				WorkDays: [7]WorkDay{},
				DaysDistribution: [7][3]int{},
			},
			Lessons: make(map[string]Lesson),
		}
	}

	return teachers
}

// 从sql中加载班级信息
func ImportClasses(db *gorm.DB) map[string]Class {
	classes := make((map[string]Class))
	rows := make([]mysql.ClassSql, 0)
	result := db.Find(&rows, "status = ?", true)
	if result.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", result.Error)
		return classes
	}

	for _, row := range rows {
		var classMates map[string]struct{}
		if err := json.Unmarshal(row.ClassMates, &classMates); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}

		classes[row.ClassName] = Class{
			Dict: common.ClassDict{
				ClassId: row.ClassId,
				ClassName: row.ClassName,
				ClassMates: classMates,
			},
			Extend: ClassExtend{
				SpareTime: make(map[int]map[int]struct{}),
			},
			Lessons: make(map[string]Lesson),
		}
	}

	return classes
}

// 从sql中加载课程信息
func ImportLessons(db *gorm.DB) []Lesson {
	rows := make([]mysql.LessonSql, 0)
	result := db.Find(&rows)
	if result.Error != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", result.Error)
		return []Lesson{}
	}

	lessons := make([]Lesson, len(rows))
	for i, row := range rows {
		lessons[i] = Lesson{
			Dict: common.LessonDict{
				LessonId: row.LessonId,
				LessonName: row.LessonName,
				TeacherName: row.TeacherName,
				StudyName: row.StudyName,
				StudentNum: row.StudentNum,
			},
			Extend: LessonExtend{
				SpareTime: make(map[int]map[int]struct{}),
				CandidateDays: make([]CandidateDay, 0),
			},
		}
	}

	return lessons
}
