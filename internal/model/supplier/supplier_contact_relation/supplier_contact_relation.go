package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SupplierContactRelation struct {
	ContactID  int64     `json:"contact_id"`
	SupplierID int64     `json:"supplier_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}

func (csr *SupplierContactRelation) Validate() error {
	if csr.ContactID <= 0 {
		return &validators.ValidationError{
			Field:   "contact_id",
			Message: "campo obrigatório e deve ser maior que zero",
		}
	}

	if csr.SupplierID <= 0 {
		return &validators.ValidationError{
			Field:   "supplier_id",
			Message: "campo obrigatório e deve ser maior que zero",
		}
	}

	return nil
}
