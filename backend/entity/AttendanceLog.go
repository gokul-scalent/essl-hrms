package entity

import "time"

type AttendanceLog struct {
	ID              int
	UID             int
	EmpID           string
	Timestamp       time.Time
	Status          int
	Punch           int
	AttendanceState string
	DeviceName      string
	CreatedAt       time.Time
}

type DailyAttendanceLog struct {
	EmpID        string
	EmpName      string
	Date         time.Time
	CheckIn      *time.Time
	CheckOut     *time.Time
	WorkingHours string
	Status       string
}
