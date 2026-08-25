package entity

import "time"

type Employee struct {
	ID        int
	UID       int
	EmpID     string
	EmpName   string
	Privilege int
	Password  string
	GroupID   string
	Card      string
	CreatedAt time.Time
	DeletedAt time.Time
}
