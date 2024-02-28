package utils

import (
	"schedule/algorithm"
	"schedule/common"
	"schedule/logs"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func FormatStudents(studentResult map[string][7][3]algorithm.LessonResult) {
	file := excelize.NewFile()

	for sName, sResult := range studentResult {
		file.NewSheet(sName)
		setCellValue(file, sName)

		for day, dayResult := range sResult {
			for duration, lessonResult := range dayResult {
				if lessonResult.IfSchedule {
					cell := common.ExcelDayMap[day] + common.ExcelDurationMap[duration]
					content := lessonResult.Lesson.Dict.LessonName + lessonResult.Lesson.Dict.TeacherName
					file.SetCellValue(sName, cell, content)
				}
			}
		}
	}

	file.SetActiveSheet(1)

	if err := file.SaveAs("students.xlsx"); err != nil {
		logs.GetInstance().Logger.Warnf("save students.xlsx error %s", err)
	}
}

func setCellValue(file *excelize.File, sName string) {
	file.SetCellValue(sName, "B1", "星期一")
	file.SetCellValue(sName, "C1", "星期二")
	file.SetCellValue(sName, "D1", "星期三")
	file.SetCellValue(sName, "E1", "星期四")
	file.SetCellValue(sName, "F1", "星期五")
	file.SetCellValue(sName, "G1", "星期六")
	file.SetCellValue(sName, "H1", "星期日")

	file.SetCellValue(sName, "A2", "8:00-10:00")
	file.SetCellValue(sName, "A3", "10:00-12:00")
	file.SetCellValue(sName, "A4", "13:00-15:00")
	file.SetCellValue(sName, "A5", "15:00-17:00")
	file.SetCellValue(sName, "A6", "18:00-20:00")
	file.SetCellValue(sName, "A7", "20:00-22:00")
}
