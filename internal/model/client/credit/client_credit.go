package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ClientCredit struct {
	ID            int64
	ClientID      int64
	AllowCredit   bool
	CreditLimit   float64
	CreditBalance float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (cc *ClientCredit) Validate() error {
	if cc.ClientID <= 0 {
		return &validators.ValidationError{Field: "ClientID", Message: "campo obrigatório"}
	}
	if cc.CreditLimit < 0 {
		return &validators.ValidationError{Field: "CreditLimit", Message: "não pode ser negativo"}
	}
	if cc.CreditBalance < 0 {
		return &validators.ValidationError{Field: "CreditBalance", Message: "não pode ser negativo"}
	}
	if cc.CreditBalance > cc.CreditLimit {
		return &validators.ValidationError{Field: "CreditBalance", Message: "não pode ser maior que o limite"}
	}

	return nil
}
