package apimodel

import "time"

type LoginResponse struct {
	Email         string `json:"email,omitempty"`
	Token         string `json:"token"`
	Role          string `json:"role"`
	IsPasswordSet string `json:"isPasswordSet"`
}

type UserResponse struct {
	ID            int          `json:"ID"`
	Email         string       `json:"email"`
	Status        string       `json:"status"`
	IsPasswordSet string       `json:"isPasswordSet,omitempty"`
	LastLoginAt   *time.Time   `json:"lastLoginAt"`
	Role          RoleResponse `json:"role,omitempty"`
}

type UserListResponse struct {
	TotalRecords       int            `json:"totalRecords"`
	NoOfRecordsPerPage int            `json:"noOfRecordsPerPage"`
	User               []UserResponse `json:"user"`
}
type RoleResponse struct {
	ID     int    `json:"ID"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}
type EmployeeResponse struct {
	ID        int    `json:"ID"`
	UID       int    `json:"uID"`
	EmpID     string `json:"empID"`
	EmpName   string `json:"empName"`
	Privilege int    `json:"privilege"`
	// Password  string `json:"password"`
	GroupID string       `json:"groupID"`
	Card    string       `json:"card"`
	Role    RoleResponse `json:"role"`
}

type EmployeeListResponse struct {
	TotalRecords       int `json:"totalRecords"`
	NoOfRecordsPerPage int `json:"noOfRecordsPerPage"`
	Employee           []EmployeeResponse
}

type AttendanceLogResponse struct {
	ID              int       `json:"ID"`
	UID             int       `json:"uID"`
	EmpID           string    `json:"empID"`
	Timestamp       time.Time `json:"timestamp"`
	Status          int       `json:"status"`
	Punch           int       `json:"punch"`
	AttendanceState string    `json:"attendanceState"`
	DeviceName      string    `json:"deviceName"`
}

type AttendanceLogListResponse struct {
	TotalRecords       int `json:"totalRecords"`
	NoOfRecordsPerPage int `json:"noOfRecordsPerPage"`
	AttendanceLog      []AttendanceLogResponse
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
