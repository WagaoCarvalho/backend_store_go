package model

import (
	"testing"
	"time"

	errval "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestProductFilter_Validate(t *testing.T) {
	t.Run("valid filter with only base fields", func(t *testing.T) {
		f := ProductFilter{}
		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid limit from base filter", func(t *testing.T) {
		f := ProductFilter{}
		f.Limit = -5

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("valid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now().Add(-24 * time.Hour)
		to := time.Now()

		f := ProductFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-24 * time.Hour)

		f := ProductFilter{
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

		f := ProductFilter{
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de atualização inválido")
	})

	t.Run("invalid MinCostPrice > MaxCostPrice", func(t *testing.T) {
		min := 200.0
		max := 100.0
		f := ProductFilter{
			MinCostPrice: &min,
			MaxCostPrice: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de preço de custo inválido")
	})

	t.Run("invalid MinSalePrice > MaxSalePrice", func(t *testing.T) {
		min := 300.0
		max := 200.0
		f := ProductFilter{
			MinSalePrice: &min,
			MaxSalePrice: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de preço de venda inválido")
	})

	t.Run("invalid MinStockQuantity > MaxStockQuantity", func(t *testing.T) {
		min := 10
		max := 5
		f := ProductFilter{
			MinStockQuantity: &min,
			MaxStockQuantity: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de estoque inválido")
	})

	t.Run("invalid MinDiscountPercent > MaxDiscountPercent", func(t *testing.T) {
		min := 20.0
		max := 10.0
		f := ProductFilter{
			MinDiscountPercent: &min,
			MaxDiscountPercent: &max,
		}
		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de desconto inválido")
	})

	t.Run("valid full filter", func(t *testing.T) {
		from := time.Now().Add(-24 * time.Hour)
		to := time.Now()
		minCost := 10.0
		maxCost := 20.0
		minSale := 30.0
		maxSale := 50.0
		minStock := 5
		maxStock := 15
		minDisc := 5.0
		maxDisc := 10.0
		status := true
		supplier := int64(1)

		f := ProductFilter{
			ProductName:        "Produto",
			Manufacturer:       "Marca",
			Status:             &status,
			SupplierID:         &supplier,
			MinCostPrice:       &minCost,
			MaxCostPrice:       &maxCost,
			MinSalePrice:       &minSale,
			MaxSalePrice:       &maxSale,
			MinStockQuantity:   &minStock,
			MaxStockQuantity:   &maxStock,
			MinDiscountPercent: &minDisc,
			MaxDiscountPercent: &maxDisc,
			CreatedFrom:        &from,
			CreatedTo:          &to,
			UpdatedFrom:        &from,
			UpdatedTo:          &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})
}
