package model

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
		isUpdate bool
		wantErr  bool
		errField string
	}{
		{
			name:     "nome em branco (criação)",
			input:    Product{ProductName: "", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "product_name",
		},
		{
			name:     "nome muito longo",
			input:    Product{ProductName: strings.Repeat("a", 256), Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "product_name",
		},
		{
			name:     "fabricante em branco",
			input:    Product{ProductName: "Produto", Manufacturer: "", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "manufacturer",
		},
		{
			name:     "fabricante muito longo",
			input:    Product{ProductName: "Produto", Manufacturer: strings.Repeat("b", 256), CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "manufacturer",
		},
		{
			name:     "preço de custo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: -1, SalePrice: 20, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "cost_price",
		},
		{
			name:     "preço de venda negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: -5, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "sale_price",
		},
		{
			name:     "preço de venda menor que custo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 20, SalePrice: 10, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "sale_price",
		},
		{
			name:     "estoque negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, StockQuantity: -1, SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "stock_quantity",
		},
		{
			name:     "código de barras inválido",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, Barcode: func() *string { s := "abc123"; return &s }(), SupplierID: &validSupplierID},
			isUpdate: false,
			wantErr:  true,
			errField: "barcode",
		},
		{
			name:     "fornecedor ausente na criação",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20},
			isUpdate: false,
			wantErr:  true,
			errField: "supplier_id",
		},
		{
			name:     "fornecedor ausente na atualização (permitido)",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20},
			isUpdate: true,
			wantErr:  false,
		},
		{
			name:     "produto inativo (permitido)",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, Status: false},
			isUpdate: false,
			wantErr:  false,
		},
		{
			name:     "desconto mínimo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, AllowDiscount: true, MinDiscountPercent: -5},
			isUpdate: false,
			wantErr:  true,
			errField: "min_discount_percent",
		},
		{
			name:     "desconto máximo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, AllowDiscount: true, MaxDiscountPercent: -10},
			isUpdate: false,
			wantErr:  true,
			errField: "max_discount_percent",
		},
		{
			name:     "desconto mínimo > máximo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, AllowDiscount: true, MinDiscountPercent: 30, MaxDiscountPercent: 20},
			isUpdate: false,
			wantErr:  true,
			errField: "discount_range",
		},
		{
			name:     "desconto máximo > 100%",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, AllowDiscount: true, MinDiscountPercent: 10, MaxDiscountPercent: 120},
			isUpdate: false,
			wantErr:  true,
			errField: "max_discount_percent",
		},
		{
			name:     "estoque mínimo negativo",
			input:    Product{ProductName: "Produto", Manufacturer: "Fab", CostPrice: 10, SalePrice: 20, SupplierID: &validSupplierID, MinStock: -1},
			isUpdate: false,
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
					MinStock:     min,
					MaxStock:     &max,
				}
			}(),
			isUpdate: false,
			wantErr:  true,
			errField: "max_stock",
		},
		{
			name: "validação de desconto mesmo com AllowDiscount = false",
			input: Product{
				ProductName:        "Produto",
				Manufacturer:       "Fab",
				CostPrice:          10,
				SalePrice:          20,
				SupplierID:         &validSupplierID,
				AllowDiscount:      false,
				MinDiscountPercent: -5, // Deve ser rejeitado mesmo com AllowDiscount = false
				MaxDiscountPercent: 0,
			},
			isUpdate: false,
			wantErr:  true,
			errField: "min_discount_percent",
		},
		{
			name: "produto válido na criação",
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
				AllowDiscount:      true,
				MinDiscountPercent: 5,
				MaxDiscountPercent: 20,
			},
			isUpdate: false,
			wantErr:  false,
		},
		{
			name: "produto válido na atualização sem fornecedor",
			input: Product{
				ProductName:        "Produto",
				Manufacturer:       "Fab",
				CostPrice:          10,
				SalePrice:          20,
				StockQuantity:      5,
				AllowDiscount:      false,
				MinDiscountPercent: 0,
				MaxDiscountPercent: 0,
			},
			isUpdate: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate(tt.isUpdate)

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
