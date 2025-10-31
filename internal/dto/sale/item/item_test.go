package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestToSaleItemModel(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	t.Run("conversão válida de DTO para Model", func(t *testing.T) {
		id := int64(1)
		dto := SaleItemDTO{
			ID:          &id,
			SaleID:      10,
			ProductID:   20,
			Quantity:    5,
			UnitPrice:   100.00,
			Discount:    10.00,
			Tax:         5.00,
			Subtotal:    475.00,
			Description: "Item válido",
			CreatedAt:   &now,
			UpdatedAt:   &now,
		}

		model := ToSaleItemModel(dto)
		assert.Equal(t, int64(1), model.ID)
		assert.Equal(t, int64(10), model.SaleID)
		assert.Equal(t, int64(20), model.ProductID)
		assert.Equal(t, 5, model.Quantity)
		assert.Equal(t, 100.00, model.UnitPrice)
		assert.Equal(t, 10.00, model.Discount)
		assert.Equal(t, 5.00, model.Tax)
		assert.Equal(t, 475.00, model.Subtotal)
		assert.Equal(t, "Item válido", model.Description)
		assert.False(t, model.CreatedAt.IsZero())
		assert.False(t, model.UpdatedAt.IsZero())
	})

	t.Run("datas inválidas devem ser ignoradas", func(t *testing.T) {
		id := int64(2)
		dto := SaleItemDTO{
			ID:        &id,
			SaleID:    11,
			ProductID: 22,
			Quantity:  2,
			UnitPrice: 50.00,
			CreatedAt: utils.StrToPtr("data inválida"),
			UpdatedAt: utils.StrToPtr("outra inválida"),
		}

		model := ToSaleItemModel(dto)
		assert.True(t, model.CreatedAt.IsZero())
		assert.True(t, model.UpdatedAt.IsZero())
	})
}

func TestToSaleItemDTO(t *testing.T) {
	now := time.Now()
	model := &models.SaleItem{
		ID:          1,
		SaleID:      10,
		ProductID:   20,
		Quantity:    5,
		UnitPrice:   100.00,
		Discount:    10.00,
		Tax:         5.00,
		Subtotal:    475.00,
		Description: "Item válido",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("conversão válida de Model para DTO", func(t *testing.T) {
		dto := ToSaleItemDTO(model)

		assert.Equal(t, *dto.ID, int64(1))
		assert.Equal(t, int64(10), dto.SaleID)
		assert.Equal(t, int64(20), dto.ProductID)
		assert.Equal(t, 5, dto.Quantity)
		assert.Equal(t, 100.00, dto.UnitPrice)
		assert.Equal(t, 10.00, dto.Discount)
		assert.Equal(t, 5.00, dto.Tax)
		assert.Equal(t, 475.00, dto.Subtotal)
		assert.Equal(t, "Item válido", dto.Description)
		assert.NotNil(t, dto.CreatedAt)
		assert.NotNil(t, dto.UpdatedAt)
	})
}

func TestToSaleItemDTOList(t *testing.T) {
	now := time.Now()
	modelsList := []*models.SaleItem{
		{ID: 1, SaleID: 10, ProductID: 20, Quantity: 1, UnitPrice: 50, CreatedAt: now, UpdatedAt: now},
		{ID: 2, SaleID: 11, ProductID: 21, Quantity: 2, UnitPrice: 75, CreatedAt: now, UpdatedAt: now},
	}

	t.Run("converter lista de models em lista de DTOs", func(t *testing.T) {
		dtos := ToSaleItemDTOList(modelsList)
		assert.Len(t, dtos, 2)
		assert.Equal(t, int64(1), *dtos[0].ID)
		assert.Equal(t, int64(2), *dtos[1].ID)
	})
}

func TestSaleItemDTOListToModelList(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	id1 := int64(1)
	id2 := int64(2)
	dtos := []*SaleItemDTO{
		{ID: &id1, SaleID: 10, ProductID: 20, Quantity: 1, UnitPrice: 50, CreatedAt: &now, UpdatedAt: &now},
		{ID: &id2, SaleID: 11, ProductID: 21, Quantity: 2, UnitPrice: 75, CreatedAt: &now, UpdatedAt: &now},
	}

	t.Run("converter lista de DTOs em lista de models", func(t *testing.T) {
		modelsList := SaleItemDTOListToModelList(dtos)
		assert.Len(t, modelsList, 2)
		assert.Equal(t, int64(1), modelsList[0].ID)
		assert.Equal(t, int64(2), modelsList[1].ID)
		assert.False(t, modelsList[0].CreatedAt.IsZero())
		assert.False(t, modelsList[1].UpdatedAt.IsZero())
	})
}
