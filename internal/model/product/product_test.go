package models

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestProduct_Validate(t *testing.T) {
	validSupplierID := int64(1)

	tests := []struct {
		name     string
		input    Product
		wantErr  bool
		errField string
	}{
		{
			name:     "nome em branco",
			input:    Product{ProductName: "", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "product_name",
		},
		{
			name:     "fabricante em branco",
			input:    Product{ProductName: "Produto", Manufacturer: "", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "manufacturer",
		},
		{
			name:     "preço de custo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: -1, SalePrice: 20, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "cost_price",
		},
		{
			name:     "preço de venda negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: -5, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "sale_price",
		},
		{
			name:     "preço de venda menor que custo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 20, SalePrice: 10, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "sale_price",
		},
		{
			name:     "estoque negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, StockQuantity: -1, SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "stock_quantity",
		},
		{
			name:     "código de barras inválido",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, Barcode: func() *string { s := "abc123"; return &s }(), SupplierID: &validSupplierID, Status: true},
			wantErr:  true,
			errField: "barcode",
		},

		{
			name:     "fornecedor ausente",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, Status: true},
			wantErr:  true,
			errField: "supplier_id",
		},
		{
			name:     "produto inativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: false},
			wantErr:  true,
			errField: "status",
		},
		{
			name:     "desconto negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true, AllowDiscount: true, MinDiscountPercent: -5},
			wantErr:  true,
			errField: "discount",
		},
		{
			name:     "intervalo de desconto inválido",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true, AllowDiscount: true, MinDiscountPercent: 30, MaxDiscountPercent: 20},
			wantErr:  true,
			errField: "discount_range",
		},
		{
			name:     "desconto acima de 100%",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true, AllowDiscount: true, MinDiscountPercent: 10, MaxDiscountPercent: 120},
			wantErr:  true,
			errField: "max_discount_percent",
		},
		{
			name:     "estoque mínimo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: true, MinStock: -1},
			wantErr:  true,
			errField: "min_stock",
		},
		{
			name: "estoque máximo menor que mínimo",
			input: func() Product {
				min := 10
				max := 5
				return Product{
					ProductName:  "Produto",
					Manufacturer: "Fab",
					CostPrice:    10,
					SalePrice:    20,
					SupplierID:   &validSupplierID,
					Status:       true,
					MinStock:     min,
					MaxStock:     &max,
				}
			}(),
			wantErr:  true,
			errField: "max_stock",
		},
		{
			name: "produto válido",
			input: Product{
				ProductName:        "Produto",
				Manufacturer:       "Fab",
				CostPrice:          10,
				SalePrice:          20,
				StockQuantity:      5,
				MinStock:           1,
				MaxStock:           func() *int { v := 10; return &v }(),
				Barcode:            func() *string { s := "12345678"; return &s }(),
				SupplierID:         &validSupplierID,
				Status:             true,
				AllowDiscount:      true,
				MinDiscountPercent: 5,
				MaxDiscountPercent: 20,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()

			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)

			// valida se retornou ValidationErrors
			var vErrs validators.ValidationErrors
			if ok := errors.As(err, &vErrs); ok {
				found := false
				for _, e := range vErrs {
					if strings.EqualFold(e.Field, tt.errField) {
						found = true
						break
					}
				}
				assert.True(t, found, "esperava erro no campo %q, mas erros foram: %v", tt.errField, vErrs)
			} else {
				// fallback para string
				assert.Contains(t, err.Error(), tt.errField)
			}
		})
	}
}
