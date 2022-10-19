package types

import "time"

// AutoRollCallScheduleTimeRange is The time range info of a scheduled auto roll call.
type AutoRollCallScheduleTimeRange struct {
	// Period is the Time range of automatic roll call monitoring.
	Period TimeOfDayPeriod `json:"period"`
	// SelectedWeekDay is The weekday of the schedule.
	SelectedWeekDay time.Weekday `json:"selected_week_day,omitempty"`
}

// AutoRollCallSchedule is A schedule for auto roll call.
type AutoRollCallSchedule struct {
	// ID is The id of this schedule.
	ID string `json:"id,omitempty"`
	// Enabled Represents the enabled state of this schedule.
	Enabled bool `json:"enabled,omitempty"`
	// TimeRange is the Action time definition for this schedule.
	TimeRange AutoRollCallScheduleTimeRange `json:"time_range"`
	// The target Zuvio course for this schedule.
	TargetCourse ZCourse `json:"target_course"`
	// The target zuvio user for this schedule.
	TargetUser ZUserInfo `json:"target_user"`
}
