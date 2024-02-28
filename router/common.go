package router

import "schedule/common"

type StudentResult struct {
	StudentId   int                                     `json:"studentId"`
	StudentName string                                  `json:"studentName"`
	SpareTime   [common.DayNum][common.DayDuration]bool `json:"spareTime"`
	Status      bool                                    `json:"status"`
}

type TeacherResult struct {
	TeacherId   int                                     `json:"teacherId"`
	TeacherName string                                  `json:"teacherName"`
	SpareTime   [common.DayNum][common.DayDuration]bool `json:"spareTime"`
	HolidayNum  int                                     `json:"holidayNum"`
	Status      bool                                    `json:"status"`
}

type ClassResult struct {
	ClassId    int    `json:"classId"`
	ClassName  string `json:"className"`
	ClassMates string  `json:"classMates"`
	Status     bool   `json:"status"`
}

type LessonResult struct {
	LessonId    int    `json:"lessonId"`
	LessonName  string `json:"lessonName"`
	StudyName   string `json:"studyName"`
	TeacherName string `json:"teacherName"`
	StudentNum  int    `json:"studentNum"`
}

func ConvertSpareTime(newSpareTime map[int]map[int]struct{}, spareTime [common.DayNum][common.DayDuration]bool) {
	for day, durations := range spareTime {
		for duration, ifSpare := range durations {
			if !ifSpare {
				continue
			}
			if _, ok := newSpareTime[day]; !ok {
				newSpareTime[day] = make(map[int]struct{})
			}
			spareTime[day][duration - 1] = true
		}
	}
}
