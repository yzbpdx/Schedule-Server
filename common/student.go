package common

type StudentDict struct {
	StudentId   int                `json:"studentId"`
	StudentName string             `json:"name"`
	Class       string             `json:"class"`
	SpareTime   []SpareTime        `json:"spareTime"`
	Lesson      []LessonForStudent `json:"lesson"`
	Priority    float64            `json:"-"`
}

type LessonForStudent struct {
	LessonId   int    `json:"lessonId"`
	LessonName string `json:"lessonName"`
	Teacher    string `json:"teacher"`
}
