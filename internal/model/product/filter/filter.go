package model

import (
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ProductFilter struct {
	filter.BaseFilter

	ProductName        string
	Manufacturer       string
	Barcode            string
	Status             *bool
	SupplierID         *int64
	Version            *int
	MinCostPrice       *float64
	MaxCostPrice       *float64
	MinSalePrice       *float64
	MaxSalePrice       *float64
	MinStockQuantity   *int
	MaxStockQuantity   *int
	AllowDiscount      *bool
	MinDiscountPercent *float64
	MaxDiscountPercent *float64
	CreatedFrom        *time.Time
	CreatedTo          *time.Time
	UpdatedFrom        *time.Time
	UpdatedTo          *time.Time
}

func (f *ProductFilter) Validate() error {
	if err := f.BaseFilter.Validate(); err != nil {
		return err
	}

	// Valida intervalos de preço de custo
	if f.MinCostPrice != nil && f.MaxCostPrice != nil && *f.MinCostPrice > *f.MaxCostPrice {
		return &validators.ValidationError{
			Field:   "MinCostPrice/MaxCostPrice",
			Message: "intervalo de preço de custo inválido",
		}
	}

	// Valida intervalos de preço de venda
	if f.MinSalePrice != nil && f.MaxSalePrice != nil && *f.MinSalePrice > *f.MaxSalePrice {
		return &validators.ValidationError{
			Field:   "MinSalePrice/MaxSalePrice",
			Message: "intervalo de preço de venda inválido",
		}
	}

	// Valida estoque
	if f.MinStockQuantity != nil && f.MaxStockQuantity != nil && *f.MinStockQuantity > *f.MaxStockQuantity {
		return &validators.ValidationError{
			Field:   "MinStockQuantity/MaxStockQuantity",
			Message: "intervalo de estoque inválido",
		}
	}

	// Valida intervalo de desconto
	if f.MinDiscountPercent != nil && f.MaxDiscountPercent != nil && *f.MinDiscountPercent > *f.MaxDiscountPercent {
		return &validators.ValidationError{
			Field:   "MinDiscountPercent/MaxDiscountPercent",
			Message: "intervalo de desconto inválido",
		}
	}

	// Valida intervalo de criação
	if f.CreatedFrom != nil && f.CreatedTo != nil && f.CreatedFrom.After(*f.CreatedTo) {
		return &validators.ValidationError{
			Field:   "CreatedFrom/CreatedTo",
			Message: "intervalo de criação inválido",
		}
	}

	// Valida intervalo de atualização
	if f.UpdatedFrom != nil && f.UpdatedTo != nil && f.UpdatedFrom.After(*f.UpdatedTo) {
		return &validators.ValidationError{
			Field:   "UpdatedFrom/UpdatedTo",
			Message: "intervalo de atualização inválido",
		}
	}

	return nil
}
