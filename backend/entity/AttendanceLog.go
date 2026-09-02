package entity

import "time"

type AttendanceLog struct {
	ID              int
	UID             int
	EmpID           string
	EmpName         string
	Timestamp       time.Time
	Status          int
	Punch           int
	AttendanceState string
	DeviceName      string
	CreatedAt       time.Time
}

type AttendancePunch struct {
	CheckIn  *time.Time
	CheckOut *time.Time
}

type DailyAttendanceLog struct {
	EmpID        string
	EmpName      string
	Date         time.Time
	Punches      []AttendancePunch
	WorkingHours string
	Status       string
}
