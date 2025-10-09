package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ClientContactRelations struct {
	ClientID  int64     `json:"client_id"`
	ContactID int64     `json:"contact_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (ccr *ClientContactRelations) Validate() error {
	if ccr.ClientID <= 0 {
		return &validators.ValidationError{
			Field:   "client_id",
			Message: "campo obrigatório e deve ser maior que zero",
		}
	}

	if ccr.ContactID <= 0 {
		return &validators.ValidationError{
			Field:   "contact_id",
			Message: "campo obrigatório e deve ser maior que zero",
		}
	}

	return nil
}
