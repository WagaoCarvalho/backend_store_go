package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryRelationsDTO(t *testing.T) {
	// Modelo de teste
	createdAt := time.Now()
	model := &models.ProductCategoryRelation{
		ProductID:  1,
		CategoryID: 10,
		CreatedAt:  createdAt,
	}

	t.Run("To DTO", func(t *testing.T) {
		dto := ToProductCategoryRelationsDTO(model)

		assert.Equal(t, model.ProductID, dto.ProductID)
		assert.Equal(t, model.CategoryID, dto.CategoryID)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
	})

	t.Run("To DTO with nil model", func(t *testing.T) {
		dto := ToProductCategoryRelationsDTO(nil)
		assert.Equal(t, int64(0), dto.ProductID)
		assert.Equal(t, int64(0), dto.CategoryID)
		assert.Equal(t, "", dto.CreatedAt)
	})

	t.Run("To Model", func(t *testing.T) {
		dto := ProductCategoryRelationsDTO{
			ProductID:  2,
			CategoryID: 20,
		}

		model := ToProductCategoryRelationsModel(dto)

		assert.Equal(t, dto.ProductID, model.ProductID)
		assert.Equal(t, dto.CategoryID, model.CategoryID)
		assert.WithinDuration(t, time.Now(), model.CreatedAt, time.Second) // CreatedAt definido como now()
	})
}
