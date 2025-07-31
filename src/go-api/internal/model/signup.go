package model

import (
	"github.com/lib/pq"
	"time"
)

type Signup struct {
	ID             int         `json:"id" db:"id"`
	DayID          int         `json:"day_id" db:"day_id"`
	RegistrationID int         `json:"registration_id" db:"registration_id"`
	State          string      `json:"state" db:"state"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	DeletedAt      pq.NullTime `json:"-" db:"deleted_at"`
}
