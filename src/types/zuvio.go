package types

import (
	"github.com/NEO-TAT/tat_auto_roll_call_service/src/types/z_roll_call_record_type"
	"time"
)

// ZUserInfo is a struct that contains information about a user.
// The Zuvio IRS student data model, includes some fields we may use in some cases.
// Note that the real user data model is not just these fields,
// there are many more that exist in the API, but we do not use.
type ZUserInfo struct {
	// The `user_id` of this user.
	ID string `json:"id,omitempty"`
	// The accessToken of this user.
	AccessToken string `json:"access_token,omitempty"`
	// The name of this user.
	Name string `json:"name,omitempty"`
	// Indicates if this user has joined any Zuvio course.
	HasCourse bool `json:"has_course,omitempty"`
}

// ZCourse is The course in Zuvio IRS.
// Note that the real user data model is not just these fields,
// there are many more that exist in the API, but we do not use.
type ZCourse struct {
	// ID is the unique identifier of the course.
	ID string `json:"id,omitempty"`
	// SemesterID is the unique identifier of the semester.
	SemesterID string `json:"semester_id,omitempty"`
	// SemesterName is the name of the semester.
	SemesterName string `json:"semester_name,omitempty"`
	// TeacherName is the name of the teacher.
	TeacherName string `json:"teacher_name,omitempty"`
	// Name is the name of the course.
	Name string `json:"name,omitempty"`
	// IsSpecialCourse shows whether the course is a special course.
	IsSpecialCourse string `json:"is_special_course,omitempty"`
}

// ZRollCall is The roll call in Zuvio IRS.
type ZRollCall struct {
	// ID is The id of this rol-call.
	ID string `json:"id,omitempty"`
	// CourseID is The id of the course corresponding to this roll-call.
	CourseID string `json:"course_id,omitempty"`
	// CreateAt is The creation time of this roll-call.
	CreateAt time.Time `json:"create_at"`
	// Record is The record of this roll-call, see ZRollCallRecord for more details.
	Record ZRollCallRecord `json:"record"`
	// IsGps indicates whether this is a roll call that requires location information.
	IsGps bool `json:"is_gps,omitempty"`
}

// ZRollCallRecord is The `record` field of a ZRollCall.
type ZRollCallRecord struct {
	// Type is The type of this record, default is z_roll_call_record_type.NonArrival.
	Type z_roll_call_record_type.ZRollCallRecordType `json:"type,omitempty"`
	// Answered Indicates whether the roll call record is an already registered record.
	Answered bool `json:"answered,omitempty"`
	// Timestamp is The register time of this record.
	// It will be nil if Answered is not true.
	Timestamp *time.Time `json:"timestamp,omitempty"`
}
