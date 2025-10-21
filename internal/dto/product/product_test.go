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
	barcode := "1234567890123"
	created := time.Now()
	updated := created.Add(time.Hour)
	maxStock := 20

	dto := ProductDTO{
		ID:                 &id,
		SupplierID:         &supplierID,
		ProductName:        "Produto X",
		Manufacturer:       "Fabricante Y",
		Description:        "Descrição do produto",
		CostPrice:          10.5,
		SalePrice:          15.0,
		StockQuantity:      100,
		MinStock:           5,
		MaxStock:           &maxStock,
		Barcode:            &barcode,
		Status:             true,
		Version:            2,
		AllowDiscount:      true,
		MinDiscountPercent: 5,
		MaxDiscountPercent: 20,
		CreatedAt:          &created,
		UpdatedAt:          &updated,
	}

	model := ToProductModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, supplierID, *model.SupplierID)
	assert.Equal(t, "Produto X", model.ProductName)
	assert.Equal(t, "Fabricante Y", model.Manufacturer)
	assert.Equal(t, "Descrição do produto", model.Description)
	assert.Equal(t, 10.5, model.CostPrice)
	assert.Equal(t, 15.0, model.SalePrice)
	assert.Equal(t, 100, model.StockQuantity)
	assert.Equal(t, 5, model.MinStock)
	assert.Equal(t, 20, *model.MaxStock)
	assert.Equal(t, barcode, *model.Barcode)
	assert.True(t, model.Status)
	assert.Equal(t, 2, model.Version)
	assert.True(t, model.AllowDiscount)
	assert.Equal(t, 5.0, model.MinDiscountPercent)
	assert.Equal(t, 20.0, model.MaxDiscountPercent)
	assert.Equal(t, created, model.CreatedAt)
	assert.Equal(t, updated, model.UpdatedAt)
}

func TestToProductDTO(t *testing.T) {
	created := time.Now()
	updated := created.Add(time.Hour)

	supplierID := int64(10)
	barcode := "1234567890123"
	maxStock := 20

	model := &models.Product{
		ID:                 1,
		SupplierID:         &supplierID,
		ProductName:        "Produto X",
		Manufacturer:       "Fabricante Y",
		Description:        "Descrição do produto",
		CostPrice:          10.5,
		SalePrice:          15.0,
		StockQuantity:      100,
		MinStock:           5,
		MaxStock:           &maxStock,
		Barcode:            &barcode,
		Status:             true,
		Version:            2,
		AllowDiscount:      true,
		MinDiscountPercent: 5,
		MaxDiscountPercent: 20,
		CreatedAt:          created,
		UpdatedAt:          updated,
	}

	dto := ToProductDTO(model)

	assert.NotNil(t, dto.ID)
	assert.Equal(t, int64(1), *dto.ID)
	assert.Equal(t, model.SupplierID, dto.SupplierID)
	assert.Equal(t, "Produto X", dto.ProductName)
	assert.Equal(t, "Fabricante Y", dto.Manufacturer)
	assert.Equal(t, "Descrição do produto", dto.Description)
	assert.Equal(t, 10.5, dto.CostPrice)
	assert.Equal(t, 15.0, dto.SalePrice)
	assert.Equal(t, 100, dto.StockQuantity)
	assert.Equal(t, 5, dto.MinStock)
	assert.Equal(t, 20, *dto.MaxStock)
	assert.Equal(t, barcode, *dto.Barcode)
	assert.True(t, dto.Status)
	assert.Equal(t, 2, dto.Version)
	assert.True(t, dto.AllowDiscount)
	assert.Equal(t, 5.0, dto.MinDiscountPercent)
	assert.Equal(t, 20.0, dto.MaxDiscountPercent)
	assert.Equal(t, created, *dto.CreatedAt)
	assert.Equal(t, updated, *dto.UpdatedAt)
}

func TestToProductModel_NilValues(t *testing.T) {
	dto := ProductDTO{
		ID:         nil,
		SupplierID: nil,
		Barcode:    nil,
	}

	model := ToProductModel(dto)

	assert.Equal(t, int64(0), model.ID)
	assert.Nil(t, model.SupplierID)
	assert.Nil(t, model.Barcode)
}

func TestToProductDTO_NilAndZeroValues(t *testing.T) {
	model := &models.Product{
		ID:         0, // deve gerar ID nil no DTO
		SupplierID: nil,
		Barcode:    nil,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	dto := ToProductDTO(model)

	// ID deve ser nil porque model.ID == 0
	assert.Nil(t, dto.ID)
	assert.Nil(t, dto.SupplierID)
	assert.Nil(t, dto.Barcode)

	// createdAt e updatedAt devem ser retornados como ponteiros válidos
	assert.NotNil(t, dto.CreatedAt)
	assert.Equal(t, model.CreatedAt, *dto.CreatedAt)
	assert.NotNil(t, dto.UpdatedAt)
	assert.Equal(t, model.UpdatedAt, *dto.UpdatedAt)
}

func TestToProductDTOs(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	t.Run("Convert multiple products", func(t *testing.T) {
		products := []*models.Product{
			{
				ID:            1,
				ProductName:   "Produto A",
				SalePrice:     100,
				StockQuantity: 10,
				Status:        true,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
			{
				ID:            2,
				ProductName:   "Produto B",
				SalePrice:     200,
				StockQuantity: 20,
				Status:        false,
				CreatedAt:     createdAt,
				UpdatedAt:     updatedAt,
			},
		}

		dtos := ToProductDTOs(products)

		assert.Len(t, dtos, 2)
		assert.Equal(t, products[0].ID, *dtos[0].ID)
		assert.Equal(t, products[0].ProductName, dtos[0].ProductName)
		assert.Equal(t, products[1].ID, *dtos[1].ID)
		assert.Equal(t, products[1].ProductName, dtos[1].ProductName)
	})

	t.Run("Handle empty slice", func(t *testing.T) {
		var products []*models.Product
		dtos := ToProductDTOs(products)
		assert.Empty(t, dtos)
	})

	t.Run("Handle slice with nil elements", func(t *testing.T) {
		products := []*models.Product{nil}
		dtos := ToProductDTOs(products)
		assert.Empty(t, dtos)
	})
}
