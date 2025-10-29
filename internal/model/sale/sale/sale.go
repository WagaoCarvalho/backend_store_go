package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Sale struct {
	ID            int64
	ClientID      *int64
	UserID        *int64
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

// --- Validação estrutural ---
func (s *Sale) ValidateStructural() error {
	var errs validators.ValidationErrors

	if s.TotalAmount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_amount", Message: "must be >= 0"})
	}
	if s.TotalDiscount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_discount", Message: "must be >= 0"})
	}
	if validators.IsBlank(s.PaymentType) {
		errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgRequiredField})
	} else if len(s.PaymentType) > 50 {
		errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgMax50})
	}
	if validators.IsBlank(s.Status) {
		errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgRequiredField})
	} else if len(s.Status) > 50 {
		errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgMax50})
	}
	if len(s.Notes) > 500 {
		errs = append(errs, validators.ValidationError{Field: "notes", Message: "max 500 characters"})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// --- Regras de negócio ---
func (s *Sale) ValidateBusinessRules() error {
	var errs validators.ValidationErrors

	if s.TotalDiscount > s.TotalAmount {
		errs = append(errs, validators.ValidationError{
			Field:   "total_discount",
			Message: "discount cannot exceed total amount",
		})
	}
	if s.SaleDate.IsZero() {
		errs = append(errs, validators.ValidationError{
			Field:   "sale_date",
			Message: "sale_date cannot be empty",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
