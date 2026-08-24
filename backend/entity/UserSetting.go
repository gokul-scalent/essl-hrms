package entity

import "time"

type UserSetting struct {
	ID                   int
	User                 User
	VerificationInterval int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            time.Time
}
