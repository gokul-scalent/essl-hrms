package entity

import (
	"database/sql"
	"time"
)

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
	Synthesized     bool
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

type DailyAttendanceLogHours struct {
	EmpID        string
	EmpName      string
	Date         time.Time
	CheckInTime  sql.NullTime
	CheckOutTime sql.NullTime
	WorkingHours string
	Status       string
}

type WorkingHoursSummary struct {
	EmpID            string
	EmpName          string
	Date             string
	CheckInTime      *time.Time
	CheckOutTime     *time.Time
	WorkingHours     string
	OutOfOfficeHours string
	Status           string
}
