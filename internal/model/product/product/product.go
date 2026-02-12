package model

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

func (p *Product) Validate(isUpdate bool) error {
	var errs validators.ValidationErrors

	// --- ProductName ---
	if validators.IsBlank(p.ProductName) {
		errs = append(errs, validators.ValidationError{
			Field:   "product_name",
			Message: validators.MsgRequiredField,
		})
	} else if len(p.ProductName) > 255 {
		errs = append(errs, validators.ValidationError{
			Field:   "product_name",
			Message: "nome do produto máximo 255 caracteres",
		})
	}

	// --- Manufacturer ---
	if validators.IsBlank(p.Manufacturer) {
		errs = append(errs, validators.ValidationError{
			Field:   "manufacturer",
			Message: validators.MsgRequiredField,
		})
	} else if len(p.Manufacturer) > 255 {
		errs = append(errs, validators.ValidationError{
			Field:   "manufacturer",
			Message: "fabricante máximo 255 caracteres",
		})
	}

	// --- SupplierID (obrigatório apenas na criação) ---
	if !isUpdate && p.SupplierID == nil {
		errs = append(errs, validators.ValidationError{
			Field:   "supplier_id",
			Message: "fornecedor é obrigatório",
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

	// --- Barcode ---
	if p.Barcode != nil && !validators.IsBlank(*p.Barcode) && !barcodeRegex.MatchString(*p.Barcode) {
		errs = append(errs, validators.ValidationError{
			Field:   "barcode",
			Message: "código de barras inválido (8-14 dígitos numéricos)",
		})
	}

	// --- Descontos (sempre valida valores, mesmo se AllowDiscount = false) ---
	if p.MinDiscountPercent < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "min_discount_percent",
			Message: "desconto mínimo não pode ser negativo",
		})
	}
	if p.MaxDiscountPercent < 0 {
		errs = append(errs, validators.ValidationError{
			Field:   "max_discount_percent",
			Message: "desconto máximo não pode ser negativo",
		})
	}
	if p.MaxDiscountPercent > 100 {
		errs = append(errs, validators.ValidationError{
			Field:   "max_discount_percent",
			Message: "desconto máximo não pode exceder 100%",
		})
	}
	if p.MinDiscountPercent > p.MaxDiscountPercent {
		errs = append(errs, validators.ValidationError{
			Field:   "discount_range",
			Message: "desconto mínimo não pode ser maior que o máximo",
		})
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
