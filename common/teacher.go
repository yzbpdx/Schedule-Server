package common

type TeacherDict struct {
	TeacherId   int                `json:"teacherId"`
	TeacherName string             `json:"teacherName"`
	SpareTime   []SpareTime        `json:"spareTime"`
	Lesson      []LessonForTeacher `json:"-"`
	Priority    float64            `json:"-"`
}

type LessonForTeacher struct {
	LessonId    int
	LessonName  string
	StudentName string
}
