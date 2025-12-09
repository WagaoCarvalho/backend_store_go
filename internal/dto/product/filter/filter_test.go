package dto

import (
	"testing"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	"github.com/stretchr/testify/assert"
)

func TestProductFilterDTO_ToModel(t *testing.T) {
	t.Run("Converte campos preenchidos corretamente", func(t *testing.T) {
		status := true
		allowDiscount := true
		version := 2
		supplierID := int64(15)
		minCost := "10.5"
		maxCost := "25.7"
		minSale := "12.0"
		maxSale := "30.0"
		minStock := "5"
		maxStock := "50"
		minDiscount := "1.5"
		maxDiscount := "10.0"
		createdFrom := "2024-01-01"
		createdTo := "2024-12-31"
		updatedFrom := "2024-06-01"
		updatedTo := "2024-06-30"

		dto := ProductFilterDTO{
			ProductName:        "Produto Teste",
			Manufacturer:       "Fabricante X",
			Barcode:            "ABC123",
			Status:             &status,
			SupplierID:         &supplierID,
			Version:            &version,
			MinCostPrice:       &minCost,
			MaxCostPrice:       &maxCost,
			MinSalePrice:       &minSale,
			MaxSalePrice:       &maxSale,
			MinStockQuantity:   &minStock,
			MaxStockQuantity:   &maxStock,
			AllowDiscount:      &allowDiscount,
			MinDiscountPercent: &minDiscount,
			MaxDiscountPercent: &maxDiscount,
			CreatedFrom:        &createdFrom,
			CreatedTo:          &createdTo,
			UpdatedFrom:        &updatedFrom,
			UpdatedTo:          &updatedTo,
			Limit:              10,
			Offset:             5,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		assert.Equal(t, dto.ProductName, model.ProductName)
		assert.Equal(t, dto.Manufacturer, model.Manufacturer)
		assert.Equal(t, dto.Barcode, model.Barcode)
		assert.Equal(t, dto.Status, model.Status)
		assert.NotNil(t, model.SupplierID)
		assert.Equal(t, int64(15), *model.SupplierID)
		assert.Equal(t, version, *model.Version)
		assert.Equal(t, float64(10.5), *model.MinCostPrice)
		assert.Equal(t, float64(25.7), *model.MaxCostPrice)
		assert.Equal(t, float64(12.0), *model.MinSalePrice)
		assert.Equal(t, float64(30.0), *model.MaxSalePrice)
		assert.Equal(t, 5, *model.MinStockQuantity)
		assert.Equal(t, 50, *model.MaxStockQuantity)
		assert.Equal(t, float64(1.5), *model.MinDiscountPercent)
		assert.Equal(t, float64(10.0), *model.MaxDiscountPercent)
		assert.Equal(t, allowDiscount, *model.AllowDiscount)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 5}, model.BaseFilter)

		assert.NotNil(t, model.CreatedFrom)
		assert.NotNil(t, model.CreatedTo)
		assert.NotNil(t, model.UpdatedFrom)
		assert.NotNil(t, model.UpdatedTo)

		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-12-31", model.CreatedTo.Format("2006-01-02"))
		assert.Equal(t, "2024-06-01", model.UpdatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-06-30", model.UpdatedTo.Format("2006-01-02"))
	})

	t.Run("Retorna nil para campos vazios", func(t *testing.T) {
		dto := ProductFilterDTO{}
		model, err := dto.ToModel()
		assert.NoError(t, err)

		assert.Nil(t, model.SupplierID)
		assert.Nil(t, model.MinCostPrice)
		assert.Nil(t, model.MaxCostPrice)
		assert.Nil(t, model.MinSalePrice)
		assert.Nil(t, model.MaxSalePrice)
		assert.Nil(t, model.MinStockQuantity)
		assert.Nil(t, model.MaxStockQuantity)
		assert.Nil(t, model.MinDiscountPercent)
		assert.Nil(t, model.MaxDiscountPercent)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
	})

	t.Run("Ignora valores inválidos", func(t *testing.T) {
		invalidInt := "abc"
		invalidFloat := "xyz"
		invalidDate := "31/12/2024"
		dto := ProductFilterDTO{
			MinCostPrice:     &invalidFloat,
			MinStockQuantity: &invalidInt,
			CreatedFrom:      &invalidDate,
			SupplierID:       nil,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		assert.Nil(t, model.MinCostPrice)
		assert.Nil(t, model.MinStockQuantity)
		assert.Nil(t, model.CreatedFrom)
	})
}
