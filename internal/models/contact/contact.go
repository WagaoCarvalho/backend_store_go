package models

import "time"

type Contact struct {
	ID              *int64    `json:"id"`
	UserID          *int64    `json:"user_id,omitempty"`
	ClientID        *int64    `json:"client_id,omitempty"`
	SupplierID      *int64    `json:"supplier_id,omitempty"`
	ContactName     string    `json:"contact_name"`
	ContactPosition string    `json:"contact_position,omitempty"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	Cell            string    `json:"cell,omitempty"`
	ContactType     string    `json:"contact_type,omitempty"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
