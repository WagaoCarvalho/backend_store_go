package model

import (
	"testing"
	"time"

	errval "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestSaleFilter_Validate(t *testing.T) {
	t.Run("valid filter with only base fields", func(t *testing.T) {
		f := SaleFilter{}
		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid limit from base filter", func(t *testing.T) {
		f := SaleFilter{}
		f.Limit = -5

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("valid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now().Add(-24 * time.Hour)
		to := time.Now()

		f := SaleFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-24 * time.Hour)

		f := SaleFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.IsType(t, &errval.ValidationError{}, err)
		assert.Contains(t, err.Error(), "intervalo de criação inválido")
	})

	t.Run("invalid UpdatedFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-1 * time.Hour)

		f := SaleFilter{
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de atualização inválido")
	})

	t.Run("valid SaleDateFrom/To range", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)

		f := SaleFilter{
			SaleDateFrom: &from,
			SaleDateTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid SaleDateFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-24 * time.Hour)

		f := SaleFilter{
			SaleDateFrom: &from,
			SaleDateTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de data da venda inválido")
	})

	t.Run("SaleDateFrom in future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		f := SaleFilter{
			SaleDateFrom: &future,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data da venda não pode estar no futuro")
	})

	t.Run("SaleDateTo in future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)
		f := SaleFilter{
			SaleDateTo: &future,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "data da venda não pode estar no futuro")
	})

	t.Run("invalid MinTotalItemsAmount > MaxTotalItemsAmount", func(t *testing.T) {
		min := 200.0
		max := 100.0
		f := SaleFilter{
			MinTotalItemsAmount: &min,
			MaxTotalItemsAmount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de valor total dos itens inválido")
	})

	t.Run("invalid MinTotalItemsDiscount > MaxTotalItemsDiscount", func(t *testing.T) {
		min := 50.0
		max := 20.0
		f := SaleFilter{
			MinTotalItemsDiscount: &min,
			MaxTotalItemsDiscount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de desconto dos itens inválido")
	})

	t.Run("invalid MinTotalSaleDiscount > MaxTotalSaleDiscount", func(t *testing.T) {
		min := 30.0
		max := 10.0
		f := SaleFilter{
			MinTotalSaleDiscount: &min,
			MaxTotalSaleDiscount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de desconto da venda inválido")
	})

	t.Run("invalid MinTotalAmount > MaxTotalAmount", func(t *testing.T) {
		min := 500.0
		max := 300.0
		f := SaleFilter{
			MinTotalAmount: &min,
			MaxTotalAmount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de valor total da venda inválido")
	})

	t.Run("negative MinTotalItemsAmount", func(t *testing.T) {
		min := -10.0
		f := SaleFilter{
			MinTotalItemsAmount: &min,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MaxTotalItemsAmount", func(t *testing.T) {
		max := -5.0
		f := SaleFilter{
			MaxTotalItemsAmount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MinTotalItemsDiscount", func(t *testing.T) {
		min := -15.0
		f := SaleFilter{
			MinTotalItemsDiscount: &min,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MaxTotalItemsDiscount", func(t *testing.T) {
		max := -5.0
		f := SaleFilter{
			MaxTotalItemsDiscount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MinTotalSaleDiscount", func(t *testing.T) {
		min := -8.0
		f := SaleFilter{
			MinTotalSaleDiscount: &min,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MaxTotalSaleDiscount", func(t *testing.T) {
		max := -3.0
		f := SaleFilter{
			MaxTotalSaleDiscount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MinTotalAmount", func(t *testing.T) {
		min := -100.0
		f := SaleFilter{
			MinTotalAmount: &min,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("negative MaxTotalAmount", func(t *testing.T) {
		max := -50.0
		f := SaleFilter{
			MaxTotalAmount: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não pode ser negativo")
	})

	t.Run("MinTotalItemsDiscount greater than MaxTotalAmount", func(t *testing.T) {
		minDiscount := 200.0
		maxAmount := 150.0
		f := SaleFilter{
			MinTotalItemsDiscount: &minDiscount,
			MaxTotalAmount:        &maxAmount,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "desconto dos itens não pode ser maior que o valor total")
	})

	t.Run("MinTotalSaleDiscount greater than MaxTotalAmount", func(t *testing.T) {
		minDiscount := 100.0
		maxAmount := 80.0
		f := SaleFilter{
			MinTotalSaleDiscount: &minDiscount,
			MaxTotalAmount:       &maxAmount,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "desconto da venda não pode ser maior que o valor total")
	})

	t.Run("invalid PaymentType", func(t *testing.T) {
		f := SaleFilter{
			PaymentType: "invalid_payment",
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tipo de pagamento inválido")
	})

	t.Run("valid PaymentType", func(t *testing.T) {
		paymentTypes := []string{"cash", "card", "credit", "pix"}
		for _, paymentType := range paymentTypes {
			t.Run(paymentType, func(t *testing.T) {
				f := SaleFilter{
					PaymentType: paymentType,
				}
				err := f.Validate()
				assert.NoError(t, err)
			})
		}
	})

	t.Run("invalid Status", func(t *testing.T) {
		f := SaleFilter{
			Status: "invalid_status",
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status inválido")
	})

	t.Run("valid Status", func(t *testing.T) {
		statuses := []string{"active", "canceled", "returned", "completed"}
		for _, status := range statuses {
			t.Run(status, func(t *testing.T) {
				f := SaleFilter{
					Status: status,
				}
				err := f.Validate()
				assert.NoError(t, err)
			})
		}
	})

	t.Run("valid full filter", func(t *testing.T) {
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)
		dayBefore := now.Add(-48 * time.Hour)
		clientID := int64(1)
		userID := int64(2)
		minItemsAmount := 100.0
		maxItemsAmount := 500.0
		minItemsDiscount := 10.0
		maxItemsDiscount := 50.0
		minSaleDiscount := 5.0
		maxSaleDiscount := 20.0
		minTotal := 150.0
		maxTotal := 450.0

		f := SaleFilter{
			ClientID:              &clientID,
			UserID:                &userID,
			PaymentType:           "card",
			Status:                "completed",
			MinTotalItemsAmount:   &minItemsAmount,
			MaxTotalItemsAmount:   &maxItemsAmount,
			MinTotalItemsDiscount: &minItemsDiscount,
			MaxTotalItemsDiscount: &maxItemsDiscount,
			MinTotalSaleDiscount:  &minSaleDiscount,
			MaxTotalSaleDiscount:  &maxSaleDiscount,
			MinTotalAmount:        &minTotal,
			MaxTotalAmount:        &maxTotal,
			Notes:                 "Test note",
			SaleDateFrom:          &dayBefore,
			SaleDateTo:            &yesterday,
			CreatedFrom:           &dayBefore,
			CreatedTo:             &yesterday,
			UpdatedFrom:           &dayBefore,
			UpdatedTo:             &yesterday,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with partial ranges", func(t *testing.T) {
		minAmount := 100.0
		maxDiscount := 50.0
		fromDate := time.Now().Add(-72 * time.Hour)

		f := SaleFilter{
			MinTotalAmount:       &minAmount,
			MaxTotalSaleDiscount: &maxDiscount,
			SaleDateFrom:         &fromDate,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with only text fields", func(t *testing.T) {
		f := SaleFilter{
			PaymentType: "pix",
			Status:      "active",
			Notes:       "Venda com desconto",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty string PaymentType and Status are valid", func(t *testing.T) {
		f := SaleFilter{
			PaymentType: "",
			Status:      "",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})
}
