package algorithm

import "schedule/logs"

func OpenOtherDuration(schedule *Schedule) {
	students := make(map[string]*Student)
	for sName, student := range schedule.Students {
		students[sName] = &Student{
			Dict:    student.Dict,
			Extend:  StudentExtend{},
			Lessons: make(map[string]*Lesson),
		}
	}
	teachers := make(map[string]*Teacher)
	for tName, teacher := range schedule.Teachers {
		teachers[tName] = &Teacher{
			Dict: teacher.Dict,
			Extend: TeacherExtend{
				Holidays: teacher.Extend.Holidays,
				WorkDays: teacher.Extend.WorkDays,
			},
			Lessons: make(map[string]*Lesson),
		}
	}
	classes := make(map[string]*Class)
	for cName, class := range schedule.Classes {
		classes[cName] = &Class{
			Dict:    class.Dict,
			Extend:  ClassExtend{},
			Lessons: make(map[string]*Lesson),
		}
	}
	lessonsMap := schedule.UnScheduleLessons
	var lessons []*Lesson
	for _, lesson := range lessonsMap {
		lessons = append(lessons, lesson)
	}
	DispatchLessons(lessons, students, teachers, classes)
	ProcessStudents(students, false)
	ProcessClasses(classes, students, false)
	ProcessTeachers(teachers, classes, students, false)
	ProcessLessons(lessons, students, teachers, classes, false)

	var scheduleOpenDuration Schedule
	scheduleOpenDuration.Students = students
	scheduleOpenDuration.Classes = classes
	scheduleOpenDuration.Teachers = teachers
	scheduleOpenDuration.Lessons = lessons
	scheduleOpenDuration.FinalResult = initScheduleResult()
	scheduleOpenDuration.TmpResult = initScheduleResult()
	scheduleOpenDuration.UnScheduleLessons = make(map[int]*Lesson)
	scheduleOpenDuration.UnPerfectResult = initScheduleResult()
	scheduleOpenDuration.LessonNum = len(scheduleOpenDuration.Lessons)

	scheduleOpenDuration.backTrackingSchedule(0, 0, false)
	logs.GetInstance().Logger.Infof("finish open other duration with %v unscheduled lesson", len(scheduleOpenDuration.UnScheduleLessons))
}
