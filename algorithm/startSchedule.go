package algorithm

import (
	"bytes"
	"encoding/gob"
	_ "fmt"
	"schedule/common"
	"schedule/logs"
	"schedule/mysql"
)

type LessonResult struct {
	IfSchedule bool
	Lesson *Lesson
}

type ScheduleResult struct {
	StudentResult map[string][7][3]LessonResult
	TeacherResult map[string][7][3]LessonResult
}

type Schedule struct {
	Students map[string]*Student
	Classes  map[string]*Class
	Teachers map[string]*Teacher
	Lessons  []*Lesson

	FinalResult       ScheduleResult
	TmpResult         ScheduleResult
	UnScheduleLessons []*Lesson
	UnPerfectResult   ScheduleResult

	LessonNum  int
	FindResult bool
	FinishNum  int
	StartFrom  int
}

// 回溯初排课程，得到初排结果
func StartSchedule() {
	students := ImportStudents(mysql.GetClient())
	classes := ImportClasses(mysql.GetClient())
	teachers := ImportTeachers(mysql.GetClient())
	lessons := ImportLessons(mysql.GetClient())
	DispatchLessons(lessons, students, teachers, classes)
	ProcessStudents(students, true)
	ProcessClasses(classes, students, true)
	ProcessTeachers(teachers, classes, students, true)
	ProcessLessons(lessons, students, teachers, classes, true)

	var schedule Schedule
	schedule.Students = students
	schedule.Classes = classes
	schedule.Teachers = teachers
	schedule.Lessons = lessons
	schedule.FinalResult = initScheduleResult()
	schedule.TmpResult = initScheduleResult()
	schedule.UnScheduleLessons = make([]*Lesson, 0)
	schedule.UnPerfectResult = initScheduleResult()
	schedule.LessonNum = len(schedule.Lessons)

	schedule.backTrackingSchedule(0, 0, true)
	logs.GetInstance().Logger.Infof("finish schedule with %v unscheduled lesson", len(schedule.UnScheduleLessons))

	if len(schedule.UnScheduleLessons) > 0 {
		OpenOtherDuration(&schedule)
	}
}

// 回溯课程的可能安排时间得到最终结果或者得到未完全成功和冲突课程的结果
func (s *Schedule) backTrackingSchedule(startIndex, finishNum int, originDuration bool) {
	if s.FindResult {
		return
	}
	if startIndex == s.LessonNum {
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
			if _, ok := s.TmpResult.StudentResult[sName]; !ok {
				s.TmpResult.StudentResult[sName] = [7][3]LessonResult{}
			}
		}

		for _, candidateDay := range lesson.Extend.CandidateDays {
			day, duration := candidateDay.Day, candidateDay.Duration
			if !originDuration {
				duration = common.OtherDuration[duration]
			}
			if s.checkLessonTimeConfict(day, duration, sName, tName, isClass, originDuration) {
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
			logs.GetInstance().Logger.Infof("set %v %v at day %v duration %v with %v", lesson.Dict.StudyName, lesson.Dict.TeacherName, day, duration, lesson.Extend.Priority)
			s.TmpResult.TeacherResult[tName] = teacherResult
			s.backTrackingSchedule(i + 1, finishNum, originDuration)
			if s.FindResult {
				return
			}
			finishNum--
			if isClass {
				classMates := s.Classes[sName].Dict.ClassMates
				for stu := range classMates {
					studentResult := s.TmpResult.StudentResult[stu]
					studentResult[day][duration] = LessonResult{
						IfSchedule: false,
						Lesson: &Lesson{},
					}
					s.TmpResult.StudentResult[stu] = studentResult
				}
			} else {
				studentResult := s.TmpResult.StudentResult[sName]
				studentResult[day][duration] = LessonResult{
					IfSchedule: false,
					Lesson: &Lesson{},
				}
				s.TmpResult.StudentResult[sName] = studentResult
			}
			teacherResult = s.TmpResult.TeacherResult[tName]
			teacherResult[day][duration] = LessonResult{
				IfSchedule: false,
				Lesson: &Lesson{},
			}
			s.TmpResult.TeacherResult[tName] = teacherResult
		}

		s.TmpResult = *deepCopy(s.UnPerfectResult)
		s.UnScheduleLessons = append(s.UnScheduleLessons, s.Lessons[s.StartFrom])
		logs.GetInstance().Logger.Infof("find conflict at %v", s.StartFrom)
		logs.GetInstance().Logger.Infof("studyName %v teacherName %v", s.Lessons[s.StartFrom].Dict.StudyName, s.Lessons[s.StartFrom].Dict.TeacherName)
		// s.backTrackingSchedule(s.StartFrom + 1, s.FinishNum + 1)
	}

	// if !s.FindResult {
	// 	s.TmpResult = *deepCopy(s.UnPerfectResult)
	// 	s.UnScheduleLessons = append(s.UnScheduleLessons, s.Lessons[s.StartFrom])
	// 	logs.GetInstance().Logger.Infof("find conflict at %v", s.StartFrom)
	// 	s.backTrackingSchedule(s.StartFrom + 1, s.FinishNum + 1)
	// }
}

// 检查当前时间是否与已经安排课程有冲突
func (s *Schedule) checkLessonTimeConfict(day, duration int, sName, tName string, isClass bool, originDuration bool) bool {
	if !originDuration {
		duration = common.OtherDuration[duration]
	}
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

// 深拷贝临时结果
func deepCopy(r ScheduleResult) *ScheduleResult {
	var buf bytes.Buffer
	gob.Register(ScheduleResult{})
	gob.Register(LessonResult{})
	gob.Register(Lesson{})
	gob.Register(common.LessonDict{})
	gob.Register(LessonExtend{})
	gob.Register(CandidateDay{})
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

func initScheduleResult() ScheduleResult {
	return ScheduleResult{
		StudentResult: make(map[string][7][3]LessonResult),
		TeacherResult: make(map[string][7][3]LessonResult),
	}
}
