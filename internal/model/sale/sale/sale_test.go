package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSale_ValidateStructural(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   100.0,
			TotalItemsDiscount: 10.0,
			TotalSaleDiscount:  5.0,
			TotalAmount:        95.0,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "Venda válida",
			SaleDate:           time.Now(),
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
			Version:            1,
		}

		err := s.ValidateStructural()
		assert.NoError(t, err)
	})

	t.Run("total_items_amount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:           &userID,
			TotalItemsAmount: -50.0,
			PaymentType:      "cash",
			Status:           "active",
			Version:          1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_items_amount")
	})

	t.Run("total_items_discount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   100,
			TotalItemsDiscount: -5,
			PaymentType:        "cash",
			Status:             "active",
			Version:            1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_items_discount")
	})

	t.Run("total_sale_discount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:            &userID,
			TotalItemsAmount:  100,
			TotalSaleDiscount: -5,
			PaymentType:       "cash",
			Status:            "active",
			Version:           1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_sale_discount")
	})

	t.Run("total_amount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: -10,
			PaymentType: "cash",
			Status:      "active",
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_amount")
	})

	t.Run("payment_type vazio", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			Status:      "active",
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment_type")
	})

	t.Run("payment_type maior que 50", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			PaymentType: string(make([]byte, 51)),
			Status:      "active",
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment_type")
	})

	t.Run("status vazio", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			PaymentType: "cash",
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status")
	})

	t.Run("status maior que 50", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			PaymentType: "cash",
			Status:      string(make([]byte, 51)),
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status")
	})

	t.Run("notes maior que 500", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			PaymentType: "cash",
			Status:      "active",
			Notes:       string(make([]byte, 501)),
			Version:     1,
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notes")
	})

	t.Run("version menor que 1", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   100,
			TotalItemsDiscount: 10,
			TotalSaleDiscount:  5,
			TotalAmount:        95,
			PaymentType:        "cash",
			Status:             "active",
			Version:            0, // inválido
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "version")
	})

}

func TestSale_ValidateBusinessRules(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   200.0,
			TotalItemsDiscount: 30.0,
			TotalSaleDiscount:  20.0,
			TotalAmount:        150.0,
			SaleDate:           time.Now(),
			Version:            1,
		}
		err := s.ValidateBusinessRules()
		assert.NoError(t, err)
	})

	t.Run("soma dos descontos maior que total_amount", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   100,
			TotalItemsDiscount: 60,
			TotalSaleDiscount:  50,
			TotalAmount:        100,
			SaleDate:           time.Now(),
			Version:            1,
		}
		err := s.ValidateBusinessRules()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sum of discounts")
	})

	t.Run("sale_date vazio", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:             &userID,
			TotalItemsAmount:   100,
			TotalItemsDiscount: 0,
			TotalSaleDiscount:  0,
			TotalAmount:        100,
			SaleDate:           time.Time{},
			Version:            1,
		}
		err := s.ValidateBusinessRules()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sale_date")
	})
}
