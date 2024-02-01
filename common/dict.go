package common

type StudentDict struct {
	StudentId   int
	StudentName string
	SpareTime   map[int]map[int]struct{}
}

type TeacherDict struct {
	TeacherId   int
	TeacherName string
	SpareTime   map[int]map[int]struct{}
	HolidayNum  int
}

type LessonDict struct {
	LessonId    int    `json:"lessonId" gorm:"column:lessonId"`
	LessonName  string `json:"lessonName" gorm:"column:lessonName"`
	TeacherName string `json:"teacherName" gorm:"column:teacherName"`
	StudyName   string `json:"studyName" gorm:"column:studyName"`
	StudentNum  int    `json:"studentNum" gorm:"column:studentNum"`
}

type ClassDict struct {
	ClassId    int
	ClassName  string
	ClassMates map[string]struct{}
}
