package common

type SpareTime struct {
	WeekDay  int        `json:"weekDay"`
	Duration []Duration `json:"duration"`
}

type Duration struct {
	Start string `json:"start"`
	End   string `json:"end"`
}
