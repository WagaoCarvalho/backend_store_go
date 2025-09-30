package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Sale struct {
	ID            int64
	ClientID      *int64
	UserID        int64
	SaleDate      time.Time
	TotalAmount   float64
	TotalDiscount float64
	PaymentType   string
	Status        string
	Notes         string
	Version       int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (s *Sale) Validate() error {
	var errs validators.ValidationErrors

	// --- Associação obrigatória ---
	if s.UserID == 0 {
		errs = append(errs, validators.ValidationError{Field: "user_id", Message: validators.MsgRequiredField})
	}

	if s.SaleDate.IsZero() {
		s.SaleDate = time.Now()
	}

	// --- Totais ---
	if s.TotalAmount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_amount", Message: "total_amount must be >= 0"})
	}
	if s.TotalDiscount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_discount", Message: "total_discount must be >= 0"})
	}

	// --- PaymentType ---
	if validators.IsBlank(s.PaymentType) {
		errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgRequiredField})
	} else if len(s.PaymentType) > 50 {
		errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgMax50})
	}

	// --- Status ---
	if validators.IsBlank(s.Status) {
		errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgRequiredField})
	} else if len(s.Status) > 50 {
		errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgMax50})
	}

	// --- Notes ---
	if len(s.Notes) > 500 {
		errs = append(errs, validators.ValidationError{Field: "notes", Message: "notes max 500 characters"})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
