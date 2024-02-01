package algorithm

import (
	"bytes"
	"encoding/gob"
	"schedule/logs"
	"schedule/mysql"
)

type LessonResult struct {
	IfSchedule bool
	Lesson Lesson
}

type ScheduleResult struct {
	StudentResult map[string][7][3]LessonResult
	TeacherResult map[string][7][3]LessonResult
}

type Schedule struct {
	Students map[string]Student
	Classes map[string]Class
	Teachers map[string]Teacher
	Lessons []Lesson

	FinalResult ScheduleResult
	TmpResult ScheduleResult
	UnScheduleLessons []Lesson
	UnPerfectResult ScheduleResult

	LessonNum int
	FindResult bool
	FinishNum int
	StartFrom int
}

func StartSchedule() {
	students := ImportStudents(mysql.GetClient())
	classes := ImportClasses(mysql.GetClient())
	teachers := ImportTeachers(mysql.GetClient())
	lessons := ImportLessons(mysql.GetClient())
	DispatchLessons(lessons, students, teachers, classes)
	ProcessStudents(students)
	ProcessClasses(classes, students)
	ProcessTeachers(teachers, classes, students)
	ProcessLessons(lessons, students, teachers, classes)

	var schedule Schedule
	schedule.Students = students
	schedule.Classes = classes
	schedule.Teachers = teachers
	schedule.Lessons = lessons
	schedule.FinalResult = ScheduleResult{
		StudentResult: make(map[string][7][3]LessonResult),
		TeacherResult: make(map[string][7][3]LessonResult),
	}
	schedule.TmpResult = ScheduleResult{
		StudentResult: make(map[string][7][3]LessonResult),
		TeacherResult: make(map[string][7][3]LessonResult),
	}
	schedule.UnScheduleLessons = make([]Lesson, 0)
	schedule.UnPerfectResult = ScheduleResult{
		StudentResult: make(map[string][7][3]LessonResult),
		TeacherResult: make(map[string][7][3]LessonResult),
	}
	schedule.LessonNum = len(schedule.Lessons)

	schedule.BacktrackingSchedule(0, 0)
}

func (s *Schedule) BacktrackingSchedule(startIndex, finishNum int) {
	if s.FindResult {
		return
	}
	if finishNum == s.LessonNum {
		s.FindResult = true
		s.FinalResult = *deepCopy(s.TmpResult)
		return
	}

	if finishNum > s.FinishNum {
		s.FinishNum = finishNum
		s.StartFrom = startIndex
		s.UnPerfectResult = *deepCopy(s.TmpResult)
	}

	for i := startIndex; i < s.LessonNum; i++ {
		lesson := s.Lessons[i]
		sName, tName := lesson.Dict.StudyName, lesson.Dict.TeacherName
		isClass := false
		if _, ok := s.TmpResult.TeacherResult[tName]; !ok {
			s.TmpResult.TeacherResult[tName] = [7][3]LessonResult{}
		}
		if lesson.Dict.StudentNum > 1 {
			class := s.Classes[lesson.Dict.StudyName]
			for stu := range class.Dict.ClassMates {
				if _, ok := s.TmpResult.StudentResult[stu]; !ok {
					s.TmpResult.StudentResult[stu] = [7][3]LessonResult{}
				}
			}
			isClass = true
		} else if lesson.Dict.StudentNum == 1 {
			s.TmpResult.StudentResult[sName] = [7][3]LessonResult{}
		}

		for _, candidateDay := range lesson.Extend.CandidateDays {
			day, duration := candidateDay.Day, candidateDay.Duration
			if s.checkLessonTimeConfict(day, duration, sName, tName, isClass) {
				continue
			}
			finishNum++
			if isClass {
				classMates := s.Classes[sName].Dict.ClassMates
				for stu := range classMates {
					studentResult := s.TmpResult.StudentResult[stu]
					studentResult[day][duration] = LessonResult{
						IfSchedule: true,
						Lesson: lesson,
					}
					s.TmpResult.StudentResult[stu] = studentResult
				}
			} else {
				studentResult := s.TmpResult.StudentResult[sName]
				studentResult[day][duration] = LessonResult{
					IfSchedule: true,
					Lesson: lesson,
				}
				s.TmpResult.StudentResult[sName] = studentResult
			}
			teacherResult := s.TmpResult.TeacherResult[tName]
			teacherResult[day][duration] = LessonResult{
				IfSchedule: true,
				Lesson: lesson,
			}
			s.TmpResult.TeacherResult[tName] = teacherResult
			s.BacktrackingSchedule(i + 1, finishNum)
			finishNum--
			if isClass {
				classMates := s.Classes[sName].Dict.ClassMates
				for stu := range classMates {
					studentResult := s.TmpResult.StudentResult[stu]
					studentResult[day][duration] = LessonResult{
						IfSchedule: false,
						Lesson: Lesson{},
					}
					s.TmpResult.StudentResult[stu] = studentResult
				}
			} else {
				studentResult := s.TmpResult.StudentResult[sName]
				studentResult[day][duration] = LessonResult{
					IfSchedule: false,
					Lesson: Lesson{},
				}
				s.TmpResult.StudentResult[sName] = studentResult
			}
			teacherResult = s.TmpResult.TeacherResult[tName]
			teacherResult[day][duration] = LessonResult{
				IfSchedule: false,
				Lesson: Lesson{},
			}
			s.TmpResult.TeacherResult[tName] = teacherResult
		}
	}

	if !s.FindResult {
		s.TmpResult = *deepCopy(s.UnPerfectResult)
		s.UnScheduleLessons = append(s.UnScheduleLessons, s.Lessons[s.StartFrom])
		s.BacktrackingSchedule(s.StartFrom + 1, s.FinishNum + 1)
	}
}

func (s *Schedule) checkLessonTimeConfict(day, duration int, sName, tName string, isClass bool) bool {
	if s.TmpResult.TeacherResult[tName][day][duration].IfSchedule {
		return true
	}
	if isClass {
		classMates := s.Classes[sName].Dict.ClassMates
		for stu := range classMates {
			if s.TmpResult.StudentResult[stu][day][duration].IfSchedule {
				return true
			}
		}
	} else if s.TmpResult.StudentResult[sName][day][duration].IfSchedule {
		return true
	}

	return false
}

func deepCopy(r ScheduleResult) *ScheduleResult {
	var buf bytes.Buffer
	encode := gob.NewEncoder(&buf)
	decode := gob.NewDecoder(&buf)

	err := encode.Encode(&r)
	if err != nil {
		logs.GetInstance().Logger.Errorf("deepCopy error %s", err)
		return nil
	}

	var copy ScheduleResult
	err = decode.Decode(&copy)
	if err != nil {
		logs.GetInstance().Logger.Errorf("deepCopy error %s", err)
		return nil
	}

	return &copy
}
