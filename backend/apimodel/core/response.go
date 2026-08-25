package apimodel

import "time"

type LoginResponse struct {
	Email         string `json:"email,omitempty"`
	Token         string `json:"token"`
	Role          string `json:"role"`
	IsPasswordSet string `json:"isPasswordSet"`
}

type UserResponse struct {
	ID            int        `json:"ID"`
	Email         string     `json:"email"`
	Status        string     `json:"status"`
	IsPasswordSet string     `json:"isPasswordSet,omitempty"`
	LastLoginAt   *time.Time `json:"lastLoginAt"`
}

type UserListResponse struct {
	TotalRecords       int            `json:"totalRecords"`
	NoOfRecordsPerPage int            `json:"noOfRecordsPerPage"`
	User               []UserResponse `json:"user"`
	Roles              RoleResponse   `json:"roles"`
}
type RoleResponse struct {
	ID     int    `json:"ID"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
