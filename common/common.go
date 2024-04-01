package common

const (
	DayNum            = 7
	DayDuration       = 6
	OriginDayDuration = 3
)

var OtherDuration = map[int]int{
	-1: 0,
	3:  1,
	4:  2,
}

var ExcelDayMap = map[int]string{
	0: "B",
	1: "C",
	2: "D",
	3: "E",
	4: "F",
	5: "G",
	6: "H",
}

var ExcelDurationMap = map[int]string{
	0: "3",
	1: "4",
	2: "5",
}
