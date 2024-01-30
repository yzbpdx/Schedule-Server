package algorithm

import (
	"schedule/gorm"
)

var (
	students map[string]Student
	teachers map[string]Teacher
)

type StudentSchedule struct {

}

// 开始排课
func StartSchedule() {
	ImportStudents(gorm.GetClient("student"), students)
	ImportTeachers(gorm.GetClient("teacher"), teachers)

	pairsNum := CaculateStudentPriority(students, teachers)
	CaculateTeacherPriority(teachers)

	studyPairs := make([]StudyPair, 0, pairsNum)
	SortStudyPairs(&studyPairs, students, teachers)
}

// 粗粒度生成课程表
func CoarseSchedule(studyPairs []StudyPair) {
	
}
