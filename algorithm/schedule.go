package algorithm

import (
	"schedule/gorm"
)

var (
	students map[string]Student
	teachers map[string]Teacher
)

func StartSchedule() {
	pairsNum := 0
	ImportStudents(gorm.GetClient("student"), students)
	ImportTeachers(gorm.GetClient("teacher"), teachers)

	pairsNum = CaculateStudentPriority(students, teachers)
	CaculateTeacherPriority(teachers)

	studyPairs := make([]StudyPair, 0, pairsNum)
	SortStudyPairs(&studyPairs, students, teachers)
}
