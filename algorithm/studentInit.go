package algorithm

import (
	"database/sql"
	"encoding/json"
	"schedule/common"
	"schedule/logs"
)

type Students struct {
	Students map[string]common.StudentDict
}

func (s *Students) CaculateStudentPriority(teachers *Teachers) int {
	pairsNum := 0
	for sName, sDict := range s.Students {
		lessonNum, spareNum := len(sDict.Lesson), 0
		for _, spareTime := range sDict.SpareTime {
			spareNum += len(spareTime)
		}
		sDict.Priority = float64(lessonNum) / float64(spareNum)

		for _, lesson := range sDict.Lesson {
			teachers.AddTeacherLesson(lesson.Teacher, sName, lesson.LessonName, lesson.LessonId)
			pairsNum++
		}
	}

	return pairsNum
}

func (s *Students) ImportStudents(clien *sql.DB) {
	rows, err := clien.Query("select * from student")
	if err != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var studentDcit common.StudentDict
		var spareTimeJSON, lessonJSON []byte
		if err := rows.Scan(&studentDcit.StudentId, &studentDcit.StudentName, &studentDcit.Class, &spareTimeJSON, &lessonJSON); err != nil {
			logs.GetInstance().Logger.Errorf("scan student dict error %s", err)
			continue
		}
		if err := json.Unmarshal(spareTimeJSON, &studentDcit.SpareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}
		if err := json.Unmarshal(lessonJSON, &studentDcit.Lesson); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal lesson err %s", err)
			continue
		}

		s.Students[studentDcit.StudentName] = studentDcit
	}
}
