package algorithm

import (
	"database/sql"
	"encoding/json"
	"schedule/common"
	"schedule/logs"
)

type Teachers struct {
	Teachers map[string]common.TeacherDict
}

func (t *Teachers) ImportTeachers(clien *sql.DB) {
	rows, err := clien.Query("select * from teacher")
	if err != nil {
		logs.GetInstance().Logger.Errorf("sql query teacher error %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var teacherDcit common.TeacherDict
		var spareTimeJSON []byte
		if err := rows.Scan(&teacherDcit.TeacherId, &teacherDcit.TeacherName, &spareTimeJSON); err != nil {
			logs.GetInstance().Logger.Errorf("scan student dict error %s", err)
			continue
		}
		if err := json.Unmarshal(spareTimeJSON, &teacherDcit.SpareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}

		t.Teachers[teacherDcit.TeacherName] = teacherDcit
	}
}

func (t *Teachers) CaculateTeacherPriority() {
	for _, tDict := range t.Teachers {
		lessonNum, spareNum := len(tDict.Lesson), 0
		for _, spareTime := range tDict.SpareTime {
			spareNum += len(spareTime)
		}
		tDict.Priority = float64(lessonNum) / float64(spareNum)
	}
}

func (t *Teachers) AddTeacherLesson(teacherName, studentName, lessonName string, lessonId int) {
	if techerDict, ok := t.Teachers[teacherName]; ok {
		techerDict.Lesson = append(techerDict.Lesson, common.LessonForTeacher{
			LessonId:    lessonId,
			LessonName:  lessonName,
			StudentName: studentName,
		})
	} else {
		logs.GetInstance().Logger.Warnf("find no register teacher %s", teacherName)
	}
}
