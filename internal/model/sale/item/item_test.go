package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaleItem_ValidateStructural(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		si := &SaleItem{
			SaleID:      1,
			ProductID:   2,
			Quantity:    5,
			UnitPrice:   10.5,
			Discount:    1.0,
			Tax:         0.5,
			Subtotal:    52.0,
			Description: "Item válido",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		err := si.ValidateStructural()
		assert.NoError(t, err)
	})

	t.Run("sale_id inválido", func(t *testing.T) {
		si := &SaleItem{SaleID: 0, ProductID: 1, Quantity: 1, UnitPrice: 10, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("product_id inválido", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 0, Quantity: 1, UnitPrice: 10, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("quantity inválida", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 0, UnitPrice: 10, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("unit_price negativo", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 1, UnitPrice: -10, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("discount negativo", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 1, UnitPrice: 10, Discount: -5, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("tax negativo", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 1, UnitPrice: 10, Tax: -2, Subtotal: 10}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("subtotal negativo", func(t *testing.T) {
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 1, UnitPrice: 10, Subtotal: -5}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})

	t.Run("descrição muito longa", func(t *testing.T) {
		desc := make([]byte, 501)
		for i := range desc {
			desc[i] = 'a'
		}
		si := &SaleItem{SaleID: 1, ProductID: 1, Quantity: 1, UnitPrice: 10, Subtotal: 10, Description: string(desc)}
		err := si.ValidateStructural()
		assert.Error(t, err)
	})
}

func TestSaleItem_ValidateBusinessRules(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		si := &SaleItem{
			Quantity:  2,
			UnitPrice: 10,
			Discount:  2,
			Tax:       1,
			Subtotal:  19,
		}
		err := si.ValidateBusinessRules()
		assert.NoError(t, err)
	})

	t.Run("subtotal inconsistente", func(t *testing.T) {
		si := &SaleItem{
			Quantity:  2,
			UnitPrice: 10,
			Discount:  1,
			Tax:       1,
			Subtotal:  15, // deveria ser 19
		}
		err := si.ValidateBusinessRules()
		assert.Error(t, err)
	})
}
