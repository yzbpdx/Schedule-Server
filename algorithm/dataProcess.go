package algorithm

import (
	"schedule/common"
	"schedule/logs"
	"sort"
)

// 分发课程到每个实体
func DispatchLessons(lessons []*Lesson, students map[string]*Student, teachers map[string]*Teacher, classes map[string]*Class) {
	for _, lesson := range lessons {
		lessonDict := lesson.Dict
		sName, tName := lessonDict.StudyName, lessonDict.TeacherName
		if lessonDict.StudentNum > 1 {
			class := classes[sName]
			for studentName := range class.Dict.ClassMates {
				student := students[studentName]
				student.dispatchToStudent()
			}
			class.dispatchToClass()
			class.Lessons[lessonDict.LessonName] = lesson
		} else if lessonDict.StudentNum == 1 {
			student := students[sName]
			student.dispatchToStudent()
			student.Lessons[lessonDict.LessonName] = lesson
		}
		teacher := teachers[tName]
		teacher.dispatchToTeacher()
		teacher.Lessons[lessonDict.LessonName] = lesson
	}
	logs.GetInstance().Logger.Infof("dispatch lessons success")
}

// 处理学生extend数据
func ProcessStudents(students map[string]*Student, originDuration bool) {
	for _, student := range students {
		for _, time := range student.Dict.SpareTime {
			for duration := range time {
				if originDuration && checkDuration(duration) {
					student.Extend.SpareNum++
				} else if !originDuration && !checkDuration(duration) {
					student.Extend.SpareNum++
				}
			}
		}

		student.Extend.Priority = float64(student.Extend.LessonNum) / float64(student.Extend.SpareNum)
	}
}

// 处理老师extend数据
func ProcessTeachers(teachers map[string]*Teacher, classes map[string]*Class, students map[string]*Student, originDuration bool) {
	for _, teacher := range teachers {
		for _, time := range teacher.Dict.SpareTime {
			for duration := range time {
				if originDuration && checkDuration(duration) {
					teacher.Extend.SpareNum++
				} else if !originDuration && !checkDuration(duration) {
					teacher.Extend.SpareNum++
				}
			}
		}

		teacher.Extend.Priority = float64(teacher.Extend.LessonNum) / float64(teacher.Extend.SpareNum)

		for _, lesson := range teacher.Lessons {
			sName := lesson.Dict.StudyName
			if lesson.Dict.StudentNum > 1 {
				class := classes[sName]
				teacher.setDaysDistribution(class.Extend.SpareTime, originDuration)
			} else if lesson.Dict.StudentNum == 1 {
				student := students[sName]
				teacher.setDaysDistribution(student.Dict.SpareTime, originDuration)
			}
		}

		for i, time := range teacher.Extend.DaysDistribution {
			for _, duration := range time {
				teacher.Extend.WorkDays[i].Day = i
				teacher.Extend.WorkDays[i].WorkNum += duration
			}
		}

		if originDuration {
			slice := teacher.Extend.WorkDays[:]
			sort.Slice(slice, func(i, j int) bool {
				return slice[i].WorkNum < slice[j].WorkNum
			})

			for i := 0; i < teacher.Dict.HolidayNum; i++ {
				teacher.Extend.Holidays[teacher.Extend.WorkDays[i].Day] = struct{}{}
			}
			// teacher.updateDaysDistribution()
		}
	}
}

// 处理班级extend数据
func ProcessClasses(classes map[string]*Class, students map[string]*Student, originDuration bool) {
	for _, class := range classes {
		classSpareTime := make(map[int]map[int]int)
		for sName := range class.Dict.ClassMates {
			student := students[sName]
			for day, time := range student.Dict.SpareTime {
				if _, ok := classSpareTime[day]; !ok {
					classSpareTime[day] = make(map[int]int)
				}
				spare := classSpareTime[day]
				for duration := range time {
					if originDuration && checkDuration(duration) {
						spare[duration]++
					} else if !originDuration && !checkDuration(duration) {
						spare[duration]++
					}
				}
			}
		}

		for day, time := range classSpareTime {
			for duration, num := range time {
				if num == len(class.Dict.ClassMates) {
					if _, ok := class.Extend.SpareTime[day]; !ok {
						class.Extend.SpareTime[day] = make(map[int]struct{})
					}
					spare := class.Extend.SpareTime[day]
					spare[duration] = struct{}{}
					class.Extend.SpareNum++
				}
			}
		}

		class.Extend.Priority = float64(class.Extend.LessonNum) / float64(class.Extend.SpareNum)
		// logs.GetInstance().Logger.Infof("%v spare time %v with %v", class.Dict.ClassName, class.Extend.SpareTime, class.Extend.Priority)
	}
}

// 处理课程extend数据
func ProcessLessons(lessons []*Lesson, students map[string]*Student, teachers map[string]*Teacher, classes map[string]*Class, originDuration bool) {
	for _, lesson := range lessons {
		teacher := teachers[lesson.Dict.TeacherName]
		teacherTime := teacher.Dict.SpareTime
		var studyTime map[int]map[int]struct{}
		var studyPrioirty float64
		var studyLessons map[string]*Lesson
		if lesson.Dict.StudentNum > 1 {
			class := classes[lesson.Dict.StudyName]
			studyTime = class.Extend.SpareTime
			studyPrioirty = class.Extend.Priority
			studyLessons = class.Lessons
		} else if lesson.Dict.StudentNum == 1 {
			student := students[lesson.Dict.StudyName]
			studyTime = student.Dict.SpareTime
			studyPrioirty = student.Extend.Priority
			studyLessons = student.Lessons
		}
		possibleDayNum, lessonTime := getLessonPossibleDays(studyTime, teacherTime, teacher.Extend.Holidays, originDuration)
		lessonPriority := studyPrioirty * teacher.Extend.Priority / float64(possibleDayNum)
		// logs.GetInstance().Logger.Infof("lesson %v with %v and %v", lesson.Dict.LessonName, lessonPriority, studyPrioirty)
		lesson.Extend.SpareTime = lessonTime
		lesson.Extend.Priority = lessonPriority

		var candidateDays []CandidateDay
		var teachersDays [7][3]float64
		for _, studyLesson := range studyLessons {
			lessonTeacher := teachers[studyLesson.Dict.TeacherName]
			for i, time := range lessonTeacher.Extend.DaysDistribution {
				if _, ok := lessonTeacher.Extend.Holidays[i]; ok {
					continue
				}
				for j, workNum := range time {
					if originDuration {
						teachersDays[i][j] = lessonTeacher.Extend.Priority * float64(workNum)
					} else {
						teachersDays[i][common.OtherDuration[j]] = lessonTeacher.Extend.Priority * float64(workNum)
					}
				}
			}
		}
		for day, time := range lessonTime {
			for duration := range time {
				var priority float64
				if originDuration {
					priority = teachersDays[day][duration]
				} else {
					priority = teachersDays[day][common.OtherDuration[duration]]
				}
				candidateDays = append(candidateDays, CandidateDay{
					Day:      day,
					Duration: duration,
					Priority: priority,
				})
			}
		}
		sort.Slice(candidateDays, func(i, j int) bool {
			return candidateDays[i].Priority > candidateDays[j].Priority
		})
		lesson.Extend.CandidateDays = candidateDays
	}
	sort.Slice(lessons, func(i, j int) bool {
		return lessons[i].Extend.Priority > lessons[j].Extend.Priority
	})
}

// 课程分发至学生
func (s *Student) dispatchToStudent() {
	s.Extend.LessonNum++
}

// 课程分发至老师
func (t *Teacher) dispatchToTeacher() {
	t.Extend.LessonNum++
}

// 课程分发至班级
func (c *Class) dispatchToClass() {
	c.Extend.LessonNum++
}

// 计算老师可能工作时间分布
func (t *Teacher) setDaysDistribution(spareTimeMap map[int]map[int]struct{}, originDuration bool) {
	var spareTime [7][3]int
	for day, time := range spareTimeMap {
		for duration := range time {
			if originDuration && checkDuration(duration) {
				spareTime[day][duration]++
			} else if !originDuration && !checkDuration(duration) {
				spareTime[day][common.OtherDuration[duration]]++
			}
		}
	}

	t.addDays(spareTime)
}

// 添加老师可能工作时间分布
func (t *Teacher) addDays(workNum [7][3]int) {
	for i, work := range workNum {
		for j, num := range work {
			t.Extend.DaysDistribution[i][j] += num
		}
	}
}

// 更新老师可能工作时间分布
func (t *Teacher) updateDaysDistribution() {
	for holiday := range t.Extend.Holidays {
		for i := 0; i < 3; i++ {
			t.Extend.DaysDistribution[holiday][i] = 0
		}
	}
}

// 检查时间段是否在初始排课范围
func checkDuration(duration int) bool {
	if duration >= 0 && duration < 3 {
		return true
	}
	return false
}

// 得到课程可能安排的时间
func getLessonPossibleDays(studyTime, teacherTime map[int]map[int]struct{}, teacherHolidy map[int]struct{}, originDuration bool) (int, map[int]map[int]struct{}) {
	days, spareTime := 0, make(map[int]map[int]struct{})
	for tDay, tTime := range teacherTime {
		if _, ok := teacherHolidy[tDay]; ok {
			continue
		}
		if sTime, ok := studyTime[tDay]; ok {
			for sDuration := range sTime {
				if (originDuration && checkDuration(sDuration)) || (!originDuration && !checkDuration(sDuration)) {
					if _, ok := tTime[sDuration]; ok {
						days++
						if _, ok := spareTime[tDay]; !ok {
							spareTime[tDay] = make(map[int]struct{})
						}
						spareTime[tDay][sDuration] = struct{}{}
					}
				}
			}
		}
	}

	return days, spareTime
}
