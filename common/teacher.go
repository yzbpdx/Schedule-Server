package common

type TeacherDict struct {
	TeacherId   int                      `json:"teacherId" gorm:"column:teacherId"`
	TeacherName string                   `json:"teacherName" gorm:"column:teacherName"`
	SpareTime   map[int]map[int]struct{} `json:"spareTime" gorm:"column:spareTime"`
	HolidayNum     int                      `json:"holidayNum" gorm:"column:holidayNum"`
}
