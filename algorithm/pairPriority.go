package algorithm

import (
	"sort"
)

// 
type StudyPair struct {
	StudentName   string
	TeacherName   string
	LessonName    string
	Priority      float64
	SpareTime     map[int]map[int]struct{}
	CandidateDays []CandidateDay
}

type CandidateDay struct {
	WeekDay  int
	Duration int
	WorkNum  int
}

// 按照总优先级生成并排序所有学习对
func SortStudyPairs(studyPairs *[]StudyPair, students map[string]Student, teachers map[string]Teacher) {
	for tName, tDict := range teachers {
		for _, lesson := range tDict.Lesson {
			sDcit := students[lesson.StudentName]
			tDict.SetDaysDistribution(sDcit)

			days, spareTime := matchDays(sDcit.SpareTime, tDict.SpareTime, tDict.Holidays)
			priority := sDcit.Priority * tDict.Priority / float64(days)
			*studyPairs = append(*studyPairs, StudyPair{
				StudentName: sDcit.StudentName,
				TeacherName: tName,
				LessonName:  lesson.LessonName,
				Priority:    priority,
				SpareTime:   spareTime,
				CandidateDays: make([]CandidateDay, 0),
			})
		}
	}
	teachersDaysOfStudent := studentForTeacher(students, teachers)

	for _, studyPair := range *studyPairs {
		studyPair.CandidateDays = candidateDays(teachersDaysOfStudent[studyPair.StudentName], studyPair.SpareTime)
	}

	sort.Slice(*studyPairs, func(i, j int) bool {
		return (*studyPairs)[i].Priority > (*studyPairs)[j].Priority
	})
}

func PriorityTable(studyPairs []StudyPair) [7][3]float64 {
	var pTable [7][3]float64
	for _, pair := range studyPairs {
		for day, dayDuration := range pair.SpareTime {
			for duration := range dayDuration {
				pTable[day][duration] += pair.Priority
			}
		}
	}

	return pTable
}

// 计算学生和老师共同时间段数量
func matchDays(studentTime, teacherTime map[int]map[int]struct{}, teacherHolidy map[int]struct{}) (int, map[int]map[int]struct{}) {
	days, spareTime := 0, make(map[int]map[int]struct{})
	for tDay, tTime := range teacherTime {
		if _, ok := teacherHolidy[tDay]; ok {
			continue
		}
		if sTime, ok := studentTime[tDay]; ok {
			for sDuration := range sTime {
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

	return days, spareTime
}

// 计算一个学生的老师们的总可能工作时间
func studentForTeacher(students map[string]Student, teachers map[string]Teacher) map[string][7][3]int {
	teachersDaysOfStudent := make(map[string][7][3]int)
	for sName, sDict := range students {
		var teachersDays [7][3]int
		for _, lesson := range sDict.Lesson {
			t := teachers[lesson.Teacher].DaysDistribution
			addMatrix(&teachersDays, &t)
		}
		teachersDaysOfStudent[sName] = teachersDays
	}

	return teachersDaysOfStudent
}

// 按照可能工作量对学习对的可选时间进行排序，期望选择该学生的老师们可能工作时间小的时间段
func candidateDays(teachersDays [7][3]int, spareTime map[int]map[int]struct{}) []CandidateDay {
	candidate := make([]CandidateDay, 0)
	for day, time := range spareTime {
		for duration := range time {
			candidate = append(candidate, CandidateDay{
				WeekDay: day,
				Duration: duration,
				WorkNum: teachersDays[day][duration],
			})
		}
	}

	sort.Slice(candidate, func(i, j int) bool {
		return candidate[i].WorkNum < candidate[j].WorkNum
	})

	return candidate
}

func addMatrix(target, added *[7][3]int) {
	for i := 0; i < 7; i++ {
		for j := 0; j < 3; j++ {
			(*target)[i][j] += (*added)[i][j]
		}
	}
}
