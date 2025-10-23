package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type UserContactRelation struct {
	UserID    int64
	ContactID int64
	CreatedAt time.Time
}

func (ucr *UserContactRelation) Validate() error {
	if ucr.UserID <= 0 {
		return &validators.ValidationError{Field: "UserID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	if ucr.ContactID <= 0 {
		return &validators.ValidationError{Field: "ContactID", Message: "campo obrigatório e deve ser maior que zero"}
	}

	return nil
}
