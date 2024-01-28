package algorithm

import "schedule/common"

type Students struct {
	Students map[string]common.StudentDict
}

func (s *Students) CaculateStudentPriority() {
	for sName, sDict := range s.Students {
		lessonNum, spareNum := len(sDict.Lesson), 0
		for _, spareTime := range sDict.SpareTime {
			spareNum += len(spareTime.Duration)
		}
		sDict.Priority = float64(lessonNum) / float64(spareNum)

		for _, lesson := range sDict.Lesson {
			AddTeacherLesson(lesson.Teacher, sName, lesson.LessonName, lesson.LessonId)
		}
	}
}

var students Students
