package apimodel

import (
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateUser struct {
	Email  string `json:"email" binding:"required,email"`
	Status string `json:"status" binding:"omitempty,oneof=ACTIVE INACTIVE"`
}
type User struct {
	Email        string     `json:"email" binding:"omitempty"`
	Password     string     `json:"password" binding:"omitempty"`
	Status       string     `json:"status" binding:"omitempty"`
	LastLoginAt  *time.Time `json:"lastLoginAt" binding:"omitempty"`
	SessionToken string     `json:"sessionToken" binding:"omitempty"`
}

type UpdateUserRequest struct {
	Email        string    `json:"email" binding:"omitempty"`
	Password     string    `json:"password" binding:"omitempty"`
	Status       string    `json:"status" binding:"omitempty"`
	LastLoginAt  time.Time `json:"lastLoginAt" binding:"omitempty"`
	SessionToken string    `json:"sessionToken" binding:"omitempty"`
}
type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" binding:"required,min=8"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

type Employee struct {
	UID       int    `json:"uID" binding:"required"`
	EmpID     string `json:"empID" binding:"required"`
	EmpName   string `json:"empName" binding:"required"`
	Privilege int    `json:"privilege" binding:"omitempty"`
	Password  string `json:"password" binding:"omitempty"`
	GroupID   string `json:"groupID" binding:"omitempty"`
	Card      string `json:"card" binding:"omitempty"`
}

type UpdateEmployeeRequest struct {
	UID       int    `json:"uID" binding:"omitempty"`
	EmpID     string `json:"empID" binding:"omitempty"`
	EmpName   string `json:"empName" binding:"omitempty"`
	Privilege int    `json:"privilege" binding:"omitempty"`
	Password  string `json:"password" binding:"omitempty"`
	GroupID   string `json:"groupID" binding:"omitempty"`
	Card      string `json:"card" binding:"omitempty"`
}

type AttendanceLog struct {
	UID             int       `json:"uID" binding:"required"`
	EmpID           string    `json:"empID" binding:"required"`
	Timestamp       time.Time `json:"timestamp" binding:"required"`
	Status          int       `json:"status" binding:"required"`
	Punch           int       `json:"punch"` // 0- check in and 1 -check out
	AttendanceState string    `json:"attendanceState" binding:"required"`
	DeviceName      string    `json:"deviceName" binding:"required"`
}

type UpdateAttendanceLogRequest struct {
	UID             int       `json:"uID" binding:"omitempty"`
	EmpID           string    `json:"empID" binding:"omitempty"`
	Timestamp       time.Time `json:"timestamp" binding:"omitempty"`
	Status          int       `json:"status" binding:"omitempty"`
	Punch           int       `json:"punch" binding:"omitempty"`
	AttendanceState string    `json:"attendanceState" binding:"omitempty"`
	DeviceName      string    `json:"deviceName" binding:"omitempty"`
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
