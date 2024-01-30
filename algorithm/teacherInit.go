package algorithm

import (
	"database/sql"
	"encoding/json"
	"schedule/common"
	"schedule/logs"
	"sort"
)

type Teacher struct {
	common.TeacherDict
	Lesson           []LessonForTeacher
	Priority         float64
	DaysDistribution [7][3]int
	WorkDays         [7]WorkDay
	Holidays         map[int]struct{}
}

type LessonForTeacher struct {
	LessonId    int
	LessonName  string
	StudentName string
}

type WorkDay struct {
	WeekDay int
	WorkNum int
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
		if err := rows.Scan(&teacher.TeacherId, &teacher.TeacherName, &spareTimeJSON, &teacher.HolidayNum); err != nil {
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
	t.Lesson = append(t.Lesson, LessonForTeacher{
		LessonId:    lessonId,
		LessonName:  lessonName,
		StudentName: sDict.StudentName,
	})
}

func (t *Teacher) SetDaysDistribution(student Student) {
	for day, spareTime := range student.SpareTime {
		for duration := range spareTime {
			t.DaysDistribution[day][duration]++
		}
	}

	t.preConfirmHolidey()
}

func (t *Teacher) preConfirmHolidey() {
	for day, distribution := range t.DaysDistribution {
		t.WorkDays[day] = WorkDay{
			WeekDay: day,
			WorkNum: sum(distribution[:]...),
		}
	}

	sort.Slice(t.WorkDays, func(i, j int) bool {
		return t.WorkDays[i].WorkNum < t.WorkDays[j].WorkNum
	})

	for i := 0; i < t.HolidayNum; i++ {
		t.Holidays[t.WorkDays[i].WeekDay] = struct{}{}
	}

	t.updateDaysDistribution()
}

func (t *Teacher) updateDaysDistribution() {
	for holiday := range t.Holidays {
		for i := 0; i < 3; i++ {
			t.DaysDistribution[holiday][i] = 0
		}
	}
}

func sum(nums ...int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}
