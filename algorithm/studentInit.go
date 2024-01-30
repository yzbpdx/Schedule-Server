package algorithm

import (
	"database/sql"
	"encoding/json"
	"schedule/common"
	"schedule/logs"
)

type Student struct {
	common.StudentDict
	Priority  float64
	LessonNum int
	SpareNum  int
}

func (s *Student) studentInClass(classLenssonNum int) {
	s.LessonNum += classLenssonNum
	s.SpareNum -= classLenssonNum
	s.Priority = float64(s.LessonNum) / float64(s.SpareNum)
}

func CaculateStudentPriority(students map[string]Student, teachers map[string]Teacher) int {
	pairsNum := 0
	for _, sDict := range students {
		lessonNum, spareNum := len(sDict.Lesson), 0
		for _, spareTime := range sDict.SpareTime {
			spareNum += len(spareTime)
		}
		if sDict.Class != "VIP" && len(sDict.Classmates) != 0 {
			for classmateName := range sDict.Classmates {
				classmateDict := students[classmateName]
				classmateDict.studentInClass(len(sDict.Lesson))
			}
		}
		sDict.Priority = float64(lessonNum) / float64(spareNum)
		sDict.LessonNum += lessonNum
		sDict.SpareNum += spareNum

		for _, lesson := range sDict.Lesson {
			if teacher, ok := teachers[lesson.Teacher]; ok {
				teacher.AddTeacherLesson(lesson.LessonName, lesson.LessonId, &sDict)
			}
			pairsNum++
		}
	}

	return pairsNum
}

func ImportStudents(clien *sql.DB, students map[string]Student) {
	rows, err := clien.Query("select * from student")
	if err != nil {
		logs.GetInstance().Logger.Errorf("sql query student error %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var student Student
		var spareTimeJSON, lessonJSON, classmatesJSON []byte
		if err := rows.Scan(&student.StudentId, &student.StudentName, &student.Class, &classmatesJSON, &spareTimeJSON, &lessonJSON); err != nil {
			logs.GetInstance().Logger.Errorf("scan student dict error %s", err)
			continue
		}
		if err := json.Unmarshal(spareTimeJSON, &student.SpareTime); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal spare time err %s", err)
			continue
		}
		if err := json.Unmarshal(lessonJSON, &student.Lesson); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal lesson err %s", err)
			continue
		}
		if err := json.Unmarshal(classmatesJSON, &student.Classmates); err != nil {
			logs.GetInstance().Logger.Warnf("unmarshal classmates err %s", err)
			continue
		}

		students[student.StudentName] = student
	}
}
