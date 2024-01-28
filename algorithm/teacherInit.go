package algorithm

import (
	"schedule/common"
	"schedule/logs"
)

type Teachers struct {
	Teachers map[string]common.TeacherDict
}

func AddTeacherLesson(teacherName, studentName, lessonName string, lessonId int) {
	if techerDict, ok := teachers.Teachers[teacherName]; ok {
		techerDict.Lesson = append(techerDict.Lesson, common.LessonForTeacher{
			LessonId:    lessonId,
			LessonName:  lessonName,
			StudentName: studentName,
		})
	} else {
		logs.GetInstance().Logger.Warnf("find no register teacher %s", teacherName)
	}
}

var teachers Teachers
