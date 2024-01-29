package algorithm

import (
	"database/sql"
	"encoding/json"
	"schedule/common"
	"schedule/logs"
)

type Teacher struct {
	common.TeacherDict
	WorkDays [7][3]int
	Holidays []int
}

func ImportTeachers(clien *sql.DB, teachers map[string]Teacher) {
	rows, err := clien.Query("select * from teacher")
	if err != nil {
		logs.GetInstance().Logger.Errorf("sql query teacher error %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var teacher Teacher
		var spareTimeJSON []byte
		if err := rows.Scan(&teacher.TeacherId, &teacher.TeacherName, &spareTimeJSON, &teacher.Holiday); err != nil {
			logs.GetInstance().Logger.Errorf("scan student dict error %s", err)
			continue
		}
		if err := json.Unmarshal(spareTimeJSON, &teacher.SpareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}

		teachers[teacher.TeacherName] = teacher
	}
}

func CaculateTeacherPriority(teachers map[string]Teacher) {
	for _, tDict := range teachers {
		lessonNum, spareNum := len(tDict.Lesson), 0
		for _, spareTime := range tDict.SpareTime {
			spareNum += len(spareTime)
		}
		tDict.Priority = float64(lessonNum) / float64(spareNum)
	}
}

func (t *Teacher) AddTeacherLesson(lessonName string, lessonId int, sDict *Student) {
	t.Lesson = append(t.Lesson, common.LessonForTeacher{
		LessonId:    lessonId,
		LessonName:  lessonName,
		StudentName: sDict.StudentName,
	})
}

func (t *Teacher) getWorkDays() {
	
}
