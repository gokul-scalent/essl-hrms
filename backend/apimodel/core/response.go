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

type EmailListResponse struct {
	ID               int       `json:"ID"`
	UserID           int       `json:"userID"`
	Name             string    `json:"name"`
	TotalRecords     int       `json:"totalRecords"`
	ProcessedRecords int       `json:"processedRecords"`
	SafeCount        int       `json:"safeCount"`
	RiskyCount       int       `json:"riskyCount"`
	InvalIDCount     int       `json:"invalIDCount"`
	UnknownCount     int       `json:"unknownCount"`
	PendingCount     int       `json:"pendingCount"`
	Priority         string    `json:"priority"`
	CreatedAt        time.Time `json:"createdAt"`
}

type EmailListListResponse struct {
	TotalRecords       int                 `json:"totalRecords"`
	NoOfRecordsPerPage int                 `json:"noOfRecordsPerPage"`
	EmailList          []EmailListResponse `json:"emailList"`
}

type LeadResponse struct {
	ID                  int        `json:"ID"`
	EmailListID         int        `json:"emailListID"`
	Priority            string     `json:"priority"`
	Email               string     `json:"email"`
	EmailProvIDer       string     `json:"emailProvIDer"`
	FirstName           string     `json:"firstName"`
	LastName            string     `json:"lastName"`
	JobTitle            string     `json:"jobTitle"`
	Company             string     `json:"company"`
	City                string     `json:"city"`
	Country             string     `json:"country"`
	Industry            string     `json:"industry"`
	IsSafe              string     `json:"isSafe"`
	FinalStatus         string     `json:"finalStatus"`
	IsReachable         string     `json:"isReachable"`
	IsDisposable        bool       `json:"isDisposable"`
	IsRoleAccount       bool       `json:"isRoleAccount"`
	VerificationStatus  string     `json:"verificationStatus"`
	VerifiedOn          *time.Time `json:"verifiedOn"`
	VerificationAttempt int        `json:"verificationAttempt"`
}

type LeadListResponse struct {
	TotalRecords       int `json:"totalRecords"`
	NoOfRecordsPerPage int `json:"noOfRecordsPerPage"`

	SafeCount    int            `json:"safeCount"`
	RiskyCount   int            `json:"riskyCount"`
	InvalidCount int            `json:"invalidCount"`
	UnknownCount int            `json:"unknownCount"`
	PendingCount int            `json:"pendingCount"`
	TimeoutCount int            `json:"timeoutCount"`
	Lead         []LeadResponse `json:"lead"`
}

type RoleResponse struct {
	ID     int    `json:"ID"`
	Name   string `json:"name"`
	Code   string `json:"code"`
	Status string `json:"status"`
}
type UserSettingResponse struct {
	ID                   int `json:"ID"`
	UserID               int `json:"userID"`
	VerificationInterval int `json:"verificationInterval"`
}

type UserSettingListResponse struct {
	TotalRecords       int `json:"totalRecords"`
	NoOfRecordsPerPage int `json:"noOfRecordsPerPage"`
	UserSetting        []UserSettingResponse
}

//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
