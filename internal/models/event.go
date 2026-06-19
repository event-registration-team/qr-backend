package models

import "time"

type Event struct {
	ID                 int        `json:"id" db:"id"`
	Title              string     `json:"title" db:"title"`
	Description        *string    `json:"description,omitempty" db:"description"`
	Location           string     `json:"location" db:"location"`
	StartAt            time.Time  `json:"start_at" db:"start_at"`
	EndAt              time.Time  `json:"end_at" db:"end_at"`
	RegistrationStatus string     `json:"registration_status" db:"registration_status"`
	RegistrationLink   string     `json:"registration_link" db:"registration_link"`
	MaxParticipants    *int       `json:"max_participants,omitempty" db:"max_participants"`
	MaterialsLink      *string    `json:"materials_link,omitempty" db:"materials_link"`
	RequirePhone       bool       `json:"require_phone" db:"require_phone"`
	RequireCarNumber   bool       `json:"require_car_number" db:"require_car_number"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}