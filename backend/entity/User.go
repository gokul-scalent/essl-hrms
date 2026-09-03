package entity

import "time"

type User struct {
	ID            int
	Email         string
	Password      string
	IsPasswordSet string
	Role          Role
	Status        string
	LastLoginAt   time.Time
	SessionToken  string
	EmpID         string
	EmpName       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}
