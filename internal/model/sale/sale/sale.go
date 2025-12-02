package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Sale struct {
	ID                 int64
	ClientID           *int64
	UserID             *int64
	SaleDate           time.Time
	TotalItemsAmount   float64
	TotalItemsDiscount float64
	TotalSaleDiscount  float64
	TotalAmount        float64
	PaymentType        string
	Status             string
	Notes              string
	Version            int
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (s *Sale) ValidateStructural() error {
	var errs validators.ValidationErrors

	if s.TotalItemsAmount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_items_amount", Message: "must be >= 0"})
	}
	if s.TotalItemsDiscount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_items_discount", Message: "must be >= 0"})
	}
	if s.TotalSaleDiscount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_sale_discount", Message: "must be >= 0"})
	}
	if s.TotalAmount < 0 {
		errs = append(errs, validators.ValidationError{Field: "total_amount", Message: "must be >= 0"})
	}

	if validators.IsBlank(s.PaymentType) {
		errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgRequiredField})
	} else {
		allowed := map[string]bool{"cash": true, "card": true, "credit": true, "pix": true}
		if !allowed[s.PaymentType] {
			errs = append(errs, validators.ValidationError{Field: "payment_type", Message: "invalid payment type"})
		}
		if len(s.PaymentType) > 50 {
			errs = append(errs, validators.ValidationError{Field: "payment_type", Message: validators.MsgMax50})
		}
	}

	if validators.IsBlank(s.Status) {
		errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgRequiredField})
	} else {
		allowed := map[string]bool{"active": true, "canceled": true, "returned": true, "completed": true}
		if !allowed[s.Status] {
			errs = append(errs, validators.ValidationError{Field: "status", Message: "invalid status"})
		}
		if len(s.Status) > 50 {
			errs = append(errs, validators.ValidationError{Field: "status", Message: validators.MsgMax50})
		}
	}

	if len(s.Notes) > 500 {
		errs = append(errs, validators.ValidationError{Field: "notes", Message: "max 500 characters"})
	}
	if s.Version < 1 {
		errs = append(errs, validators.ValidationError{
			Field:   "version",
			Message: "must be >= 1",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (s *Sale) ValidateBusinessRules() error {
	var errs validators.ValidationErrors

	totalDiscount := s.TotalItemsDiscount + s.TotalSaleDiscount

	if totalDiscount > s.TotalAmount {
		errs = append(errs, validators.ValidationError{
			Field:   "total_amount",
			Message: "sum of discounts cannot exceed total amount",
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
