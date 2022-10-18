package z_roll_call_record_type

type ZRollCallRecordType int

const (
	// NonArrival is the Default status, means not have roll-call yet.
	NonArrival ZRollCallRecordType = iota
	// Late means The Roll Call result is late but still roll-called successfully.
	Late
	// Punctual means Roll call on time.
	Punctual
	// Leave means the student has been Marked as leave on the Roll Call.
	Leave
	// Unknown is The unknown status.
	Unknown
)
