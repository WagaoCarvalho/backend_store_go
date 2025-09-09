package models

import (
	"regexp"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Product struct {
	ID                 int64
	SupplierID         *int64
	ProductName        string
	Manufacturer       string
	Description        string
	CostPrice          float64
	SalePrice          float64
	StockQuantity      int
	MinStock           int
	MaxStock           *int
	Barcode            *string
	Status             bool
	Version            int
	AllowDiscount      bool
	MinDiscountPercent float64
	MaxDiscountPercent float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

var barcodeRegex = regexp.MustCompile(`^[0-9]{8,14}$`)

func (p *Product) Validate() error {
	var errs validators.ValidationErrors

	// --- Nome ---
	if validators.IsBlank(p.ProductName) {
		errs = append(errs, validators.ValidationError{
			Field:   "product_name",
			Message: validators.MsgRequiredField,
		})
	}

	// --- Fabricante ---
	if validators.IsBlank(p.Manufacturer) {
		errs = append(errs, validators.ValidationError{
			Field:   "manufacturer",
			Message: validators.MsgRequiredField,
		})
	}

	// --- Preços ---
	if p.CostPrice < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "cost_price",
			Message: validators.MsgCostNonNegative,
		})
	}
	if p.SalePrice < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "sale_price",
			Message: validators.MsgSaleNonNegative,
		})
	}
	if p.SalePrice < p.CostPrice {
		errs = append(errs, validators.ValidationError{
			Field:   "sale_price",
			Message: "preço de venda não pode ser menor que o de custo",
		})
	}

	// --- Estoque ---
	if p.StockQuantity < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "stock_quantity",
			Message: validators.MsgStockNegative,
		})
	}
	if p.MinStock < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "min_stock",
			Message: "estoque mínimo não pode ser negativo",
		})
	}
	if p.MaxStock != nil && *p.MaxStock < p.MinStock {
		errs = append(errs, validators.ValidationError{
			Field:   "max_stock",
			Message: "estoque máximo não pode ser menor que o mínimo",
		})
	}

	// --- Código de barras ---
	// --- Código de barras ---
	if p.Barcode != nil && !validators.IsBlank(*p.Barcode) && !barcodeRegex.MatchString(*p.Barcode) {
		errs = append(errs, validators.ValidationError{
			Field:   "barcode",
			Message: "código de barras inválido (esperado entre 8 e 14 dígitos numéricos)",
		})
	}

	// --- Fornecedor ---
	if p.SupplierID == nil {
		errs = append(errs, validators.ValidationError{
			Field:   "supplier_id",
			Message: "fornecedor é obrigatório",
		})
	}

	// --- Status ---
	if !p.Status {
		errs = append(errs, validators.ValidationError{
			Field:   "status",
			Message: "produto inativo não permitido",
		})
	}

	// --- Descontos ---
	if p.AllowDiscount {
		if p.MinDiscountPercent < 0 || p.MaxDiscountPercent < 0 {
			errs = append(errs, validators.ValidationError{
				Field:   "discount",
				Message: "desconto não pode ser negativo",
			})
		}
		if p.MinDiscountPercent > p.MaxDiscountPercent {
			errs = append(errs, validators.ValidationError{
				Field:   "discount_range",
				Message: "desconto mínimo não pode ser maior que o máximo",
			})
		}
		if p.MaxDiscountPercent > 100 {
			errs = append(errs, validators.ValidationError{
				Field:   "max_discount_percent",
				Message: "desconto máximo não pode exceder 100%",
			})
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
