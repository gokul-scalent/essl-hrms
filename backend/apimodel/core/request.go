package apimodel

import (
	"mime/multipart"
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

type EmailList struct {
	UserID           int    `json:"userID" binding:"required"`
	Name             string `json:"name" binding:"required"`
	TotalRecords     int    `json:"totalRecords" binding:"required"`
	ProcessedRecords int    `json:"processedRecords" binding:"required"`
	SafeCount        int    `json:"safeCount" binding:"required"`
	RiskyCount       int    `json:"riskyCount" binding:"required"`
	InvalIDCount     int    `json:"invalIDCount" binding:"required"`
	UnknownCount     int    `json:"unknownCount" binding:"required"`
	PendingCount     int    `json:"pendingCount" binding:"required"`
	Priority         string `json:"priority" binding:"required"`
}

// add multipart form data
type CreateEmailListRequest struct {
	Name     string `form:"name" binding:"required"`
	Priority string `form:"priority" binding:"required"`
}

type UpdateEmailListRequest struct {
	UserID           int    `json:"userID" binding:"omitempty"`
	Name             string `json:"name" binding:"omitempty"`
	TotalRecords     int    `json:"totalRecords" binding:"omitempty"`
	ProcessedRecords int    `json:"processedRecords" binding:"omitempty"`
	SafeCount        int    `json:"safeCount" binding:"omitempty"`
	RiskyCount       int    `json:"riskyCount" binding:"omitempty"`
	InvalIDCount     int    `json:"invalIDCount" binding:"omitempty"`
	UnknownCount     int    `json:"unknownCount" binding:"omitempty"`
	PendingCount     int    `json:"pendingCount" binding:"omitempty"`
	Priority         string `json:"priority" binding:"omitempty"`
}

type Lead struct {
	EmailListID        int        `json:"emailListID" binding:"required"`
	Priority           string     `json:"priority" binding:"required"`
	Email              string     `json:"email" binding:"required"`
	EmailProvIDer      string     `json:"emailProvIDer" binding:"omitempty"`
	FirstName          string     `json:"firstName" binding:"omitempty"`
	LastName           string     `json:"lastName" binding:"omitempty"`
	JobTitle           string     `json:"jobTitle" binding:"omitempty"`
	Company            string     `json:"company" binding:"omitempty"`
	City               string     `json:"city" binding:"omitempty"`
	Country            string     `json:"country" binding:"omitempty"`
	Industry           string     `json:"industry" binding:"omitempty"`
	IsSafe             string     `json:"isSafe" binding:"required"`
	FinalStatus        string     `json:"finalStatus" binding:"required"`
	IsReachable        string     `json:"isReachable"`
	IsDisposable       bool       `json:"isDisposable" `
	IsRoleAccount      bool       `json:"isRoleAccount" binding:"required"`
	VerificationStatus string     `json:"verificationStatus" binding:"required"`
	VerifiedOn         *time.Time `json:"verifiedOn" binding:"omitempty"`
}

type CreateLeadRequest struct {
	Type        string                `form:"type" binding:"required,oneof=single csv"`
	EmailListID int                   `form:"emailListID" binding:"required"`
	Email       string                `form:"email"`
	File        *multipart.FileHeader `form:"file"`
}

type UpdateLeadRequest struct {
	EmailListID        int        `json:"emailListID" binding:"omitempty"`
	Priority           string     `json:"priority" binding:"omitempty"`
	Email              string     `json:"email" binding:"omitempty"`
	EmailProvIDer      string     `json:"emailProvIDer" binding:"omitempty"`
	FirstName          string     `json:"firstName" binding:"omitempty"`
	LastName           string     `json:"lastName" binding:"omitempty"`
	JobTitle           string     `json:"jobTitle" binding:"omitempty"`
	Company            string     `json:"company" binding:"omitempty"`
	City               string     `json:"city" binding:"omitempty"`
	Country            string     `json:"country" binding:"omitempty"`
	Industry           string     `json:"industry" binding:"omitempty"`
	IsSafe             string     `json:"isSafe" binding:"omitempty"`
	FinalStatus        string     `json:"finalStatus" binding:"omitempty"`
	IsReachable        string     `json:"isReachable" binding:"omitempty"`
	IsDisposable       bool       `json:"isDisposable"`
	IsRoleAccount      bool       `json:"isRoleAccount"`
	VerificationStatus string     `json:"verificationStatus" binding:"omitempty"`
	VerifiedOn         *time.Time `json:"verifiedOn" binding:"omitempty"`
}

type UserSetting struct {
	VerificationInterval int `json:"verificationInterval"`
}

type UpdateUserSettingRequest struct {
	VerificationInterval int `json:"verificationInterval" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword     string `json:"oldPassword" binding:"required,min=8"`
	NewPassword     string `json:"newPassword" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required,min=8"`
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
