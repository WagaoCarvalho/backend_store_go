package models

import (
	"regexp"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Product struct {
	ID                 int64     `json:"id"`
	SupplierID         *int64    `json:"supplier_id,omitempty"`
	ProductName        string    `json:"product_name"`
	Manufacturer       string    `json:"manufacturer"`
	Description        string    `json:"product_description,omitempty"`
	CostPrice          float64   `json:"cost_price"`
	SalePrice          float64   `json:"sale_price"`
	StockQuantity      int       `json:"stock_quantity"`
	Barcode            string    `json:"barcode,omitempty"`
	Status             bool      `json:"status"`
	Version            int       `json:"version"`
	AllowDiscount      bool      `json:"allow_discount"`
	MinDiscountPercent float64   `json:"min_discount_percent"`
	MaxDiscountPercent float64   `json:"max_discount_percent"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
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

	// --- Código de barras ---
	if !validators.IsBlank(p.Barcode) && !barcodeRegex.MatchString(p.Barcode) {
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
