package model

import (
	"time"
)

type PromoCode struct {
	ID                     int       `json:"id" db:"id"`
	Email                  string    `json:"email" db:"email"`
	Key                    string    `json:"key" db:"key"`
	AvailableRegistrations int       `json:"available_registrations" db:"available_registrations"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}
