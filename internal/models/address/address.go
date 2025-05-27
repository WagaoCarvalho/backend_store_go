package models

import "time"

type Address struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id,omitempty"`
	ClientID   *int64    `json:"client_id,omitempty"`
	SupplierID *int64    `json:"supplier_id,omitempty"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
