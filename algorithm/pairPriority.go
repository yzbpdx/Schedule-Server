package algorithm

import (
	"sort"
)

type StudyPair struct {
	StudentName string
	TeacherName string
	LessonName  string
	Priority    float64
	SpareTime   map[int]map[int]struct{}
}

func SortStudyPairs(studyPairs *[]StudyPair, students map[string]Student, teachers map[string]Teacher) {
	for sName, sDict := range students {
		for _, lesson := range sDict.Lesson {
			teacher := lesson.Teacher
			days, spareTime := matchDays(sDict.SpareTime, teachers[teacher].SpareTime)
			priority := sDict.Priority * teachers[teacher].Priority / float64(days)
			*studyPairs = append(*studyPairs, StudyPair{
				StudentName: sName,
				TeacherName: teacher,
				LessonName:  lesson.LessonName,
				Priority:    priority,
				SpareTime:   spareTime,
			})
		}
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

func matchDays(studentTime, teacherTime map[int]map[int]struct{}) (int, map[int]map[int]struct{}) {
	days, spareTime := 0, make(map[int]map[int]struct{})
	for sDay, sTime := range studentTime {
		if tTime, ok := teacherTime[sDay]; ok {
			for sDuration := range sTime {
				if _, ok := tTime[sDuration]; ok {
					days++
					if _, ok := spareTime[sDay]; !ok {
						spareTime[sDay] = make(map[int]struct{})
					}
					spareTime[sDay][sDuration] = struct{}{}
				}
			}
		}
	}

	return days, spareTime
}
