package common

type StudentDict struct {
	StudentId   int                      `json:"studentId" gorm:"column:studentId"`
	StudentName string                   `json:"studentName" gorm:"column:studentName"`
	Class       string                   `json:"class" gorm:"column:class"`
	SpareTime   map[int]map[int]struct{} `json:"spareTime" gorm:"-"`
	Lesson      []LessonForStudent       `json:"lesson" gorm:"-"`
	Priority    float64                  `json:"-" gorm:"-"`
}

type LessonForStudent struct {
	LessonId   int    `json:"lessonId"`
	LessonName string `json:"lessonName"`
	Teacher    string `json:"teacher"`
}
