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
