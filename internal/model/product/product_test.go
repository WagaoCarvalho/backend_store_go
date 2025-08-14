package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProduct_Validate(t *testing.T) {
	setup := func() *Product {
		return &Product{
			ID:                 1,
			SupplierID:         func() *int64 { v := int64(1); return &v }(),
			ProductName:        "Produto Teste",
			Manufacturer:       "Fabricante X",
			Description:        "Descrição",
			CostPrice:          10.0,
			SalePrice:          15.0,
			StockQuantity:      5,
			Barcode:            "12345678",
			Status:             true,
			Version:            1,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			AllowDiscount:      true,
			MinDiscountPercent: 5,
			MaxDiscountPercent: 20,
		}
	}

	t.Run("erro nome do produto inválido", func(t *testing.T) {
		p := setup()
		p.ProductName = ""
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidProductName)
	})

	t.Run("erro fabricante inválido", func(t *testing.T) {
		p := setup()
		p.Manufacturer = ""
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidManufacturer)
	})

	t.Run("erro preço de custo negativo", func(t *testing.T) {
		p := setup()
		p.CostPrice = -1
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidCostPrice)
	})

	t.Run("erro preço de venda negativo", func(t *testing.T) {
		p := setup()
		p.SalePrice = -1
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidSalePrice)
	})

	t.Run("erro preço de venda menor que o custo", func(t *testing.T) {
		p := setup()
		p.SalePrice = 5
		err := p.Validate()
		assert.ErrorIs(t, err, ErrSalePriceBelowCost)
	})

	t.Run("erro estoque negativo", func(t *testing.T) {
		p := setup()
		p.StockQuantity = -1
		err := p.Validate()
		assert.ErrorIs(t, err, ErrNegativeStock)
	})

	t.Run("erro barcode inválido", func(t *testing.T) {
		p := setup()
		p.Barcode = "ABC123"
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidBarcode)
	})

	t.Run("erro supplier ausente", func(t *testing.T) {
		p := setup()
		p.SupplierID = nil
		err := p.Validate()
		assert.ErrorIs(t, err, ErrSupplierRequired)
	})

	t.Run("erro produto inativo", func(t *testing.T) {
		p := setup()
		p.Status = false
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInactiveProductNotAllowed)
	})

	t.Run("erro desconto negativo", func(t *testing.T) {
		p := setup()
		p.MinDiscountPercent = -5
		err := p.Validate()
		assert.ErrorIs(t, err, ErrNegativeDiscount)
	})

	t.Run("erro desconto acima de 100%", func(t *testing.T) {
		p := setup()
		p.MaxDiscountPercent = 150
		err := p.Validate()
		assert.ErrorIs(t, err, ErrDiscountAboveLimit)
	})

	t.Run("erro range de desconto inválido", func(t *testing.T) {
		p := setup()
		p.MinDiscountPercent = 30
		p.MaxDiscountPercent = 10
		err := p.Validate()
		assert.ErrorIs(t, err, ErrInvalidDiscountRange)
	})

	t.Run("validação bem-sucedida", func(t *testing.T) {
		p := setup()
		err := p.Validate()
		assert.NoError(t, err)
	})
}
