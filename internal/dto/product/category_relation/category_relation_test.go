package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryRelationDTO(t *testing.T) {
	// Modelo de teste
	createdAt := time.Now()
	model := &models.ProductCategoryRelation{
		ProductID:  1,
		CategoryID: 10,
		CreatedAt:  createdAt,
	}

	t.Run("ToDTO - conversão de Model para DTO", func(t *testing.T) {
		dto := ToDTO(model)

		assert.Equal(t, model.ProductID, dto.ProductID)
		assert.Equal(t, model.CategoryID, dto.CategoryID)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
	})

	t.Run("ToDTO - com CreatedAt zero", func(t *testing.T) {
		modelZero := &models.ProductCategoryRelation{
			ProductID:  2,
			CategoryID: 20,
			CreatedAt:  time.Time{}, // Zero value
		}

		dto := ToDTO(modelZero)

		assert.Equal(t, modelZero.ProductID, dto.ProductID)
		assert.Equal(t, modelZero.CategoryID, dto.CategoryID)
		assert.Equal(t, "", dto.CreatedAt) // Deve ser omitido
	})

	t.Run("ToDTO - com model nil", func(t *testing.T) {
		dto := ToDTO(nil)
		assert.Equal(t, int64(0), dto.ProductID)
		assert.Equal(t, int64(0), dto.CategoryID)
		assert.Equal(t, "", dto.CreatedAt)
	})

	t.Run("ToModel - conversão de DTO para Model", func(t *testing.T) {
		dto := ProductCategoryRelationDTO{
			ProductID:  2,
			CategoryID: 20,
		}

		model := ToModel(dto)

		assert.Equal(t, dto.ProductID, model.ProductID)
		assert.Equal(t, dto.CategoryID, model.CategoryID)
		assert.True(t, model.CreatedAt.IsZero(), "CreatedAt deve ser zero value, será definido pelo banco")
	})

	t.Run("ToDTOs - slice de models", func(t *testing.T) {
		models := []*models.ProductCategoryRelation{
			{ProductID: 1, CategoryID: 10, CreatedAt: createdAt},
			{ProductID: 2, CategoryID: 20, CreatedAt: createdAt},
			nil, // Deve ser ignorado
		}

		dtos := ToDTOs(models)

		assert.Len(t, dtos, 2)
		assert.Equal(t, int64(1), dtos[0].ProductID)
		assert.Equal(t, int64(10), dtos[0].CategoryID)
		assert.Equal(t, int64(2), dtos[1].ProductID)
		assert.Equal(t, int64(20), dtos[1].CategoryID)
	})

	t.Run("ToDTOs - slice vazio", func(t *testing.T) {
		dtos := ToDTOs([]*models.ProductCategoryRelation{})
		assert.Len(t, dtos, 0)
	})

	t.Run("ToDTOs - slice nil", func(t *testing.T) {
		dtos := ToDTOs(nil)
		assert.Len(t, dtos, 0)
	})
}
