package algorithm

import (
	"schedule/gorm"
)

func StartSchedule() {
	students, teachers := new(Students), new(Teachers)
	pairsNum := 0
	students.ImportStudents(gorm.GetClient("student"))
	teachers.ImportTeachers(gorm.GetClient("teacher"))

	pairsNum = students.CaculateStudentPriority(teachers)
	teachers.CaculateTeacherPriority()

	studyPairs := make([]StudyPair, 0, pairsNum)
	SortStudyPairs(&studyPairs, students.Students, teachers.Teachers)
}
