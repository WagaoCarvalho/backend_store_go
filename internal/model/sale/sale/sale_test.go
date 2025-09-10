package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSale_Validate(t *testing.T) {
	t.Run("válido", func(t *testing.T) {
		s := &Sale{
			UserID:        1,
			TotalAmount:   100.0,
			TotalDiscount: 10.0,
			TotalTax:      5.0,
			PaymentType:   "cash",
			Status:        "active",
			Notes:         "Venda válida",
			SaleDate:      time.Now(),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		err := s.Validate()
		assert.NoError(t, err)
	})

	t.Run("user_id obrigatório", func(t *testing.T) {
		s := &Sale{
			UserID:      0,
			TotalAmount: 50,
			PaymentType: "cash",
			Status:      "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id")
	})

	t.Run("total_amount negativo", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: -50.0,
			PaymentType: "cash",
			Status:      "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_amount")
	})

	t.Run("total_discount negativo", func(t *testing.T) {
		s := &Sale{
			UserID:        1,
			TotalAmount:   100,
			TotalDiscount: -5,
			PaymentType:   "cash",
			Status:        "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_discount")
	})

	t.Run("total_tax negativo", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			TotalTax:    -3,
			PaymentType: "cash",
			Status:      "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "total_tax")
	})

	t.Run("payment_type vazio", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			Status:      "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment_type")
	})

	t.Run("payment_type maior que 50", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			PaymentType: string(make([]byte, 51)),
			Status:      "active",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "payment_type")
	})

	t.Run("status vazio", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			PaymentType: "cash",
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status")
	})

	t.Run("status maior que 50", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			PaymentType: "cash",
			Status:      string(make([]byte, 51)),
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status")
	})

	t.Run("notes maior que 500", func(t *testing.T) {
		s := &Sale{
			UserID:      1,
			TotalAmount: 100,
			PaymentType: "cash",
			Status:      "active",
			Notes:       string(make([]byte, 501)),
		}
		err := s.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "notes")
	})
}
