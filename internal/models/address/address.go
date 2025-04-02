package models

import "time"

type Address struct {
	ID         int       `json:"id"`
	UserID     *int      `json:"user_id,omitempty"`
	ClientID   *int      `json:"client_id,omitempty"`
	SupplierID *int      `json:"supplier_id,omitempty"`
	Street     string    `json:"street"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Country    string    `json:"country"`
	PostalCode string    `json:"postal_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
