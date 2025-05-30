package models

import (
	"errors"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/utils"
)

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

func (a *Address) Validate() error {
	if a.Street == "cause_generic_error" {
		return errors.New("erro genérico na validação")
	}
	if a.UserID == nil && a.ClientID == nil && a.SupplierID == nil {
		return &utils.ValidationError{Field: "UserID/ClientID/SupplierID", Message: "pelo menos um deve ser informado"}
	}
	if a.Street == "" {
		return &utils.ValidationError{Field: "Street", Message: "campo obrigatório"}
	}
	if a.City == "" {
		return &utils.ValidationError{Field: "City", Message: "campo obrigatório"}
	}
	if a.State == "" {
		return &utils.ValidationError{Field: "State", Message: "campo obrigatório"}
	}
	if a.Country == "" {
		return &utils.ValidationError{Field: "Country", Message: "campo obrigatório"}
	}
	if a.PostalCode == "" {
		return &utils.ValidationError{Field: "PostalCode", Message: "campo obrigatório"}
	}
	return nil
}
