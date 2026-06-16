package models

import "time"

type Participant struct {
	ID           int        `json:"id" db:"id"`
	EventID      int        `json:"event_id" db:"event_id"`
	LastName     string     `json:"last_name" db:"last_name"`
	FirstName    string     `json:"first_name" db:"first_name"`
	MiddleName   *string    `json:"middle_name,omitempty" db:"middle_name"`
	Email        string     `json:"email" db:"email"`
	Phone        *string    `json:"phone,omitempty" db:"phone"`
	CarNumber    *string    `json:"car_number,omitempty" db:"car_number"`
	QRToken      string     `json:"qr_token" db:"qr_token"`
	VisitStatus  string     `json:"visit_status" db:"visit_status"`
	CheckedInAt  *time.Time `json:"checked_in_at,omitempty" db:"checked_in_at"`
	RegisteredAt time.Time  `json:"registered_at" db:"registered_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
