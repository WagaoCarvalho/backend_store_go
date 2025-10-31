package model

import (
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SaleItem struct {
	ID          int64
	SaleID      int64
	ProductID   int64
	Quantity    int
	UnitPrice   float64
	Discount    float64
	Tax         float64
	Subtotal    float64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// --- Validação estrutural ---
func (s *SaleItem) ValidateStructural() error {
	var errs validators.ValidationErrors

	if s.SaleID <= 0 {
		errs = append(errs, validators.ValidationError{Field: "sale_id", Message: validators.MsgRequiredField})
	}
	if s.ProductID <= 0 {
		errs = append(errs, validators.ValidationError{Field: "product_id", Message: validators.MsgRequiredField})
	}
	if s.Quantity <= 0 {
		errs = append(errs, validators.ValidationError{Field: "quantity", Message: "must be greater than 0"})
	}
	if s.UnitPrice < 0 {
		errs = append(errs, validators.ValidationError{Field: "unit_price", Message: "must be >= 0"})
	}
	if s.Discount < 0 {
		errs = append(errs, validators.ValidationError{Field: "discount", Message: "must be >= 0"})
	}
	if s.Tax < 0 {
		errs = append(errs, validators.ValidationError{Field: "tax", Message: "must be >= 0"})
	}
	if s.Subtotal < 0 {
		errs = append(errs, validators.ValidationError{Field: "subtotal", Message: "must be >= 0"})
	}
	if len(s.Description) > 500 {
		errs = append(errs, validators.ValidationError{Field: "description", Message: "max 500 characters"})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// --- Regras de negócio ---
func (s *SaleItem) ValidateBusinessRules() error {
	var errs validators.ValidationErrors

	if s.Subtotal != float64(s.Quantity)*s.UnitPrice-(s.Discount)+s.Tax {
		errs = append(errs, validators.ValidationError{
			Field:   "subtotal",
			Message: "subtotal must equal (quantity * unit_price - discount + tax)",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
