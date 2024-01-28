package common

type TeacherDict struct {
	TeacherId   int                      `json:"teacherId" gorm:"column:teacherId"`
	TeacherName string                   `json:"teacherName" gorm:"column:teacherName"`
	SpareTime   map[int]map[int]struct{} `json:"spareTime" gorm:"column:spareTime"`
	Lesson      []LessonForTeacher       `json:"-" gorm:"-"`
	Priority    float64                  `json:"-" gorm:"-"`
}

type LessonForTeacher struct {
	LessonId    int
	LessonName  string
	StudentName string
}
