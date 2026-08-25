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

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
