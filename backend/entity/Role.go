package entity

import "time"

type Role struct {
	ID        int
	Name      string
	Code      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
