package types

type TimeOfDay struct {
	Hour   int `json:"hour,omitempty"`
	Minute int `json:"minute,omitempty"`
}

type TimeOfDayPeriod struct {
	StartTime TimeOfDay `json:"start_time"`
	EndTime   TimeOfDay `json:"end_time"`
}
