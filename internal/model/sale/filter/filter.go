package model

import (
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SaleFilter struct {
	filter.BaseFilter

	ClientID              *int64
	UserID                *int64
	PaymentType           string
	Status                string
	MinTotalItemsAmount   *float64
	MaxTotalItemsAmount   *float64
	MinTotalItemsDiscount *float64
	MaxTotalItemsDiscount *float64
	MinTotalSaleDiscount  *float64
	MaxTotalSaleDiscount  *float64
	MinTotalAmount        *float64
	MaxTotalAmount        *float64
	Notes                 string
	SaleDateFrom          *time.Time
	SaleDateTo            *time.Time
	CreatedFrom           *time.Time
	CreatedTo             *time.Time
	UpdatedFrom           *time.Time
	UpdatedTo             *time.Time
}

func (f *SaleFilter) Validate() error {
	if err := f.BaseFilter.Validate(); err != nil {
		return err
	}

	// Valida o tipo de pagamento (se fornecido)
	if f.PaymentType != "" {
		allowedPaymentTypes := map[string]bool{"cash": true, "card": true, "credit": true, "pix": true}
		if !allowedPaymentTypes[f.PaymentType] {
			return &validators.ValidationError{
				Field:   "PaymentType",
				Message: "tipo de pagamento inválido. Valores permitidos: cash, card, credit, pix",
			}
		}
	}

	// Valida o status (se fornecido)
	if f.Status != "" {
		allowedStatuses := map[string]bool{"active": true, "canceled": true, "returned": true, "completed": true}
		if !allowedStatuses[f.Status] {
			return &validators.ValidationError{
				Field:   "Status",
				Message: "status inválido. Valores permitidos: active, canceled, returned, completed",
			}
		}
	}

	// Valida intervalos de TotalItemsAmount
	if f.MinTotalItemsAmount != nil && f.MaxTotalItemsAmount != nil && *f.MinTotalItemsAmount > *f.MaxTotalItemsAmount {
		return &validators.ValidationError{
			Field:   "MinTotalItemsAmount/MaxTotalItemsAmount",
			Message: "intervalo de valor total dos itens inválido",
		}
	}

	// Valida intervalos de TotalItemsDiscount
	if f.MinTotalItemsDiscount != nil && f.MaxTotalItemsDiscount != nil && *f.MinTotalItemsDiscount > *f.MaxTotalItemsDiscount {
		return &validators.ValidationError{
			Field:   "MinTotalItemsDiscount/MaxTotalItemsDiscount",
			Message: "intervalo de desconto dos itens inválido",
		}
	}

	// Valida intervalos de TotalSaleDiscount
	if f.MinTotalSaleDiscount != nil && f.MaxTotalSaleDiscount != nil && *f.MinTotalSaleDiscount > *f.MaxTotalSaleDiscount {
		return &validators.ValidationError{
			Field:   "MinTotalSaleDiscount/MaxTotalSaleDiscount",
			Message: "intervalo de desconto da venda inválido",
		}
	}

	// Valida intervalos de TotalAmount
	if f.MinTotalAmount != nil && f.MaxTotalAmount != nil && *f.MinTotalAmount > *f.MaxTotalAmount {
		return &validators.ValidationError{
			Field:   "MinTotalAmount/MaxTotalAmount",
			Message: "intervalo de valor total da venda inválido",
		}
	}

	// Valida que valores monetários não podem ser negativos
	if f.MinTotalItemsAmount != nil && *f.MinTotalItemsAmount < 0 {
		return &validators.ValidationError{
			Field:   "MinTotalItemsAmount",
			Message: "não pode ser negativo",
		}
	}
	if f.MaxTotalItemsAmount != nil && *f.MaxTotalItemsAmount < 0 {
		return &validators.ValidationError{
			Field:   "MaxTotalItemsAmount",
			Message: "não pode ser negativo",
		}
	}
	if f.MinTotalItemsDiscount != nil && *f.MinTotalItemsDiscount < 0 {
		return &validators.ValidationError{
			Field:   "MinTotalItemsDiscount",
			Message: "não pode ser negativo",
		}
	}
	if f.MaxTotalItemsDiscount != nil && *f.MaxTotalItemsDiscount < 0 {
		return &validators.ValidationError{
			Field:   "MaxTotalItemsDiscount",
			Message: "não pode ser negativo",
		}
	}
	if f.MinTotalSaleDiscount != nil && *f.MinTotalSaleDiscount < 0 {
		return &validators.ValidationError{
			Field:   "MinTotalSaleDiscount",
			Message: "não pode ser negativo",
		}
	}
	if f.MaxTotalSaleDiscount != nil && *f.MaxTotalSaleDiscount < 0 {
		return &validators.ValidationError{
			Field:   "MaxTotalSaleDiscount",
			Message: "não pode ser negativo",
		}
	}
	if f.MinTotalAmount != nil && *f.MinTotalAmount < 0 {
		return &validators.ValidationError{
			Field:   "MinTotalAmount",
			Message: "não pode ser negativo",
		}
	}
	if f.MaxTotalAmount != nil && *f.MaxTotalAmount < 0 {
		return &validators.ValidationError{
			Field:   "MaxTotalAmount",
			Message: "não pode ser negativo",
		}
	}

	// Valida intervalo de data da venda
	if f.SaleDateFrom != nil && f.SaleDateTo != nil && f.SaleDateFrom.After(*f.SaleDateTo) {
		return &validators.ValidationError{
			Field:   "SaleDateFrom/SaleDateTo",
			Message: "intervalo de data da venda inválido",
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

	// Valida se SaleDate não está no futuro (se apenas SaleDateFrom for fornecido)
	now := time.Now()
	if f.SaleDateFrom != nil && f.SaleDateFrom.After(now) {
		return &validators.ValidationError{
			Field:   "SaleDateFrom",
			Message: "data da venda não pode estar no futuro",
		}
	}

	// Valida se SaleDate não está no futuro (se apenas SaleDateTo for fornecido)
	if f.SaleDateTo != nil && f.SaleDateTo.After(now) {
		return &validators.ValidationError{
			Field:   "SaleDateTo",
			Message: "data da venda não pode estar no futuro",
		}
	}

	// Valida que descontos não podem ser maiores que o valor total
	if f.MinTotalItemsDiscount != nil && f.MaxTotalAmount != nil && *f.MinTotalItemsDiscount > *f.MaxTotalAmount {
		return &validators.ValidationError{
			Field:   "MinTotalItemsDiscount",
			Message: "desconto dos itens não pode ser maior que o valor total",
		}
	}
	if f.MinTotalSaleDiscount != nil && f.MaxTotalAmount != nil && *f.MinTotalSaleDiscount > *f.MaxTotalAmount {
		return &validators.ValidationError{
			Field:   "MinTotalSaleDiscount",
			Message: "desconto da venda não pode ser maior que o valor total",
		}
	}

	return nil
}
