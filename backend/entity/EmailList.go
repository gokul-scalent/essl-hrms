package entity

import "time"

type EmailList struct {
	ID               int
	User             User
	Name             string
	TotalRecords     int
	ProcessedRecords int
	SafeCount        int
	RiskyCount       int
	InvalIDCount     int
	UnknownCount     int
	PendingCount     int
	Priority         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        time.Time
}
