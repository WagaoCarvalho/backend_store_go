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
			UserID:        &userID,
			TotalAmount:   100.0,
			TotalDiscount: 10.0,
			PaymentType:   "cash",
			Status:        "active",
			Notes:         "Venda válida",
			SaleDate:      time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err := s.ValidateStructural()
		assert.NoError(t, err)
	})

	t.Run("total_amount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: -50.0,
			PaymentType: "cash",
			Status:      "active",
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_amount")
	})

	t.Run("total_discount negativo", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:        &userID,
			TotalAmount:   100,
			TotalDiscount: -5,
			PaymentType:   "cash",
			Status:        "active",
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_discount")
	})

	t.Run("payment_type vazio", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:      &userID,
			TotalAmount: 100,
			Status:      "active",
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
		}
		err := s.ValidateStructural()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notes")
	})
}

func TestSale_ValidateBusinessRules(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:        &userID,
			TotalAmount:   200.0,
			TotalDiscount: 50.0,
			SaleDate:      time.Now(),
		}
		err := s.ValidateBusinessRules()
		assert.NoError(t, err)
	})

	t.Run("desconto maior que total", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:        &userID,
			TotalAmount:   100,
			TotalDiscount: 150,
			SaleDate:      time.Now(),
		}
		err := s.ValidateBusinessRules()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "discount cannot exceed total amount")
	})

	t.Run("sale_date vazio", func(t *testing.T) {
		userID := int64(1)
		s := &Sale{
			UserID:        &userID,
			TotalAmount:   100,
			TotalDiscount: 0,
			SaleDate:      time.Time{}, // zero
		}
		err := s.ValidateBusinessRules()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sale_date")
	})
}
