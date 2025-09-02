package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	"github.com/stretchr/testify/assert"
)

func TestToProductModel(t *testing.T) {
	id := int64(1)
	supplierID := int64(10)
	dto := ProductDTO{
		ID:                 &id,
		SupplierID:         &supplierID,
		ProductName:        "Produto X",
		Manufacturer:       "Fabricante Y",
		Description:        "Descrição do produto",
		CostPrice:          10.5,
		SalePrice:          15.0,
		StockQuantity:      100,
		Barcode:            "1234567890123",
		Status:             true,
		Version:            2,
		AllowDiscount:      true,
		MinDiscountPercent: 5,
		MaxDiscountPercent: 20,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	model := ToProductModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, &supplierID, model.SupplierID)
	assert.Equal(t, "Produto X", model.ProductName)
	assert.Equal(t, "Fabricante Y", model.Manufacturer)
	assert.Equal(t, "Descrição do produto", model.Description)
	assert.Equal(t, 10.5, model.CostPrice)
	assert.Equal(t, 15.0, model.SalePrice)
	assert.Equal(t, 100, model.StockQuantity)
	assert.Equal(t, "1234567890123", model.Barcode)
	assert.True(t, model.Status)
	assert.Equal(t, 2, model.Version)
	assert.True(t, model.AllowDiscount)
	assert.Equal(t, 5.0, model.MinDiscountPercent)
	assert.Equal(t, 20.0, model.MaxDiscountPercent)
	assert.Equal(t, dto.CreatedAt, model.CreatedAt)
	assert.Equal(t, dto.UpdatedAt, model.UpdatedAt)
}

func TestToProductDTO(t *testing.T) {
	created := time.Now()
	updated := created.Add(time.Hour)

	supplierID := int64(10)
	model := &models.Product{
		ID:                 1,
		SupplierID:         &supplierID,
		ProductName:        "Produto X",
		Manufacturer:       "Fabricante Y",
		Description:        "Descrição do produto",
		CostPrice:          10.5,
		SalePrice:          15.0,
		StockQuantity:      100,
		Barcode:            "1234567890123",
		Status:             true,
		Version:            2,
		AllowDiscount:      true,
		MinDiscountPercent: 5,
		MaxDiscountPercent: 20,
		CreatedAt:          created,
		UpdatedAt:          updated,
	}

	dto := ToProductDTO(model)

	assert.Equal(t, &model.ID, dto.ID)
	assert.Equal(t, model.SupplierID, dto.SupplierID)
	assert.Equal(t, model.ProductName, dto.ProductName)
	assert.Equal(t, model.Manufacturer, dto.Manufacturer)
	assert.Equal(t, model.Description, dto.Description)
	assert.Equal(t, model.CostPrice, dto.CostPrice)
	assert.Equal(t, model.SalePrice, dto.SalePrice)
	assert.Equal(t, model.StockQuantity, dto.StockQuantity)
	assert.Equal(t, model.Barcode, dto.Barcode)
	assert.Equal(t, model.Status, dto.Status)
	assert.Equal(t, model.Version, dto.Version)
	assert.Equal(t, model.AllowDiscount, dto.AllowDiscount)
	assert.Equal(t, model.MinDiscountPercent, dto.MinDiscountPercent)
	assert.Equal(t, model.MaxDiscountPercent, dto.MaxDiscountPercent)
	assert.Equal(t, model.CreatedAt, dto.CreatedAt)
	assert.Equal(t, model.UpdatedAt, dto.UpdatedAt)
}

func TestToProductModel_NilID(t *testing.T) {
	dto := ProductDTO{
		ID: nil,
	}

	model := ToProductModel(dto)
	assert.Equal(t, int64(0), model.ID)
}
