package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	"github.com/stretchr/testify/assert"
)

func TestUserCategoryRelationsDTO(t *testing.T) {
	// Modelo de teste
	createdAt := time.Now()
	model := &models.UserCategoryRelations{
		UserID:     1,
		CategoryID: 10,
		CreatedAt:  createdAt,
	}

	t.Run("To DTO", func(t *testing.T) {
		dto := ToUserCategoryRelationsDTO(model)

		assert.Equal(t, model.UserID, dto.UserID)
		assert.Equal(t, model.CategoryID, dto.CategoryID)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
	})

	t.Run("To DTO with nil model", func(t *testing.T) {
		dto := ToUserCategoryRelationsDTO(nil)
		assert.Equal(t, int64(0), dto.UserID)
		assert.Equal(t, int64(0), dto.CategoryID)
		assert.Equal(t, "", dto.CreatedAt)
	})

	t.Run("To Model", func(t *testing.T) {
		dto := UserCategoryRelationsDTO{
			UserID:     2,
			CategoryID: 20,
		}

		model := ToUserCategoryRelationsModel(dto)

		assert.Equal(t, dto.UserID, model.UserID)
		assert.Equal(t, dto.CategoryID, model.CategoryID)
		assert.WithinDuration(t, time.Now(), model.CreatedAt, time.Second) // CreatedAt definido como now()
	})
}
