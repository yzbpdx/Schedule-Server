package algorithm

import (
	"encoding/json"
	"fmt"
	"schedule/common"
	"schedule/logs"
)

type LessonJSON struct {
	Day      int               `json:"day"`
	Duration int               `json:"duration"`
	Lesson   common.LessonDict `json:"lesson"`
}

func Result2JSON(result map[string][7][common.DayDuration]LessonResult) {
	var lessonResult []LessonJSON
	for _, lesson := range result {
		for day, dayResult := range lesson {
			for duration, durationResult := range dayResult {
				if durationResult.IfSchedule {
					lessonResult = append(lessonResult, LessonJSON{
						Day:      day,
						Duration: duration,
						Lesson:   durationResult.Lesson.Dict,
					})
				}
			}
		}
	}

	lessonResultJSON, err := json.Marshal(lessonResult)
	if err != nil {
		logs.GetInstance().Logger.Errorf("marshal lesson result error %s", err)
	}

	fmt.Println(string(lessonResultJSON))
}
