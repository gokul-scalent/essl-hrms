package entity

import "time"

type Lead struct {
	ID                 int
	EmailList          EmailList
	Priority           string
	Email              string
	EmailProvIDer      string
	FirstName          string
	LastName           string
	JobTitle           string
	Company            string
	City               string
	Country            string
	Industry           string
	IsSafe             string
	FinalStatus        string
	IsReachable        string
	IsDisposable       bool
	IsRoleAccount      bool
	VerificationStatus string
	VerifiedOn         *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
	RetryCount         int
	NextRetryAt        *time.Time
}

// summary of lead status count
type LeadStatusCount struct {
	SafeCount    int `db:"safe_count"`
	RiskyCount   int `db:"risky_count"`
	InvalidCount int `db:"invalid_count"`
	UnknownCount int `db:"unknown_count"`
	PendingCount int `db:"pending_count"`
	TimeoutCount int `db:"timeout_count"`
}
