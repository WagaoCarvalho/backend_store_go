package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	"github.com/stretchr/testify/assert"
)

func TestToUserCategoryModel(t *testing.T) {
	t.Run("Converte DTO para Model corretamente", func(t *testing.T) {
		dto := UserCategoryDTO{
			ID:          nil,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
		}

		model := ToUserCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, uint(0), model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, dto.Description, model.Description)
	})

	t.Run("Converte DTO com ID definido para Model corretamente", func(t *testing.T) {
		id := uint(10)
		dto := UserCategoryDTO{
			ID:          &id,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
		}

		model := ToUserCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, id, model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, dto.Description, model.Description)
	})
}

func TestToUserCategoryDTO(t *testing.T) {
	t.Run("Converte Model para DTO corretamente", func(t *testing.T) {
		created := time.Now().Add(-1 * time.Hour)
		updated := time.Now()
		model := &models.UserCategory{
			ID:          5,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
			CreatedAt:   created,
			UpdatedAt:   updated,
		}

		dto := ToUserCategoryDTO(model)

		assert.NotNil(t, dto.ID)
		assert.Equal(t, model.ID, *dto.ID)
		assert.Equal(t, model.Name, dto.Name)
		assert.Equal(t, model.Description, dto.Description)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
		assert.Equal(t, model.UpdatedAt.Format(time.RFC3339), dto.UpdatedAt)
	})

	t.Run("Retorna DTO vazio se Model for nil", func(t *testing.T) {
		dto := ToUserCategoryDTO(nil)

		assert.NotNil(t, dto)
		assert.Nil(t, dto.ID)
		assert.Equal(t, "", dto.Name)
		assert.Equal(t, "", dto.Description)
		assert.Equal(t, "", dto.CreatedAt)
		assert.Equal(t, "", dto.UpdatedAt)
	})
}
