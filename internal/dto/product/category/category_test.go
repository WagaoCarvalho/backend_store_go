package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	"github.com/stretchr/testify/assert"
)

func TestToProductCategoryModel(t *testing.T) {
	t.Run("Converte DTO para Model corretamente", func(t *testing.T) {
		dto := ProductCategoryDTO{
			ID:          nil,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
		}

		model := ToProductCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, uint(0), model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, dto.Description, model.Description)
	})

	t.Run("Converte DTO com ID definido para Model corretamente", func(t *testing.T) {
		id := uint(10)
		dto := ProductCategoryDTO{
			ID:          &id,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
		}

		model := ToProductCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, id, model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, dto.Description, model.Description)
	})
}

func TestToProductCategoryDTO(t *testing.T) {
	t.Run("Converte Model para DTO corretamente", func(t *testing.T) {
		created := time.Now().Add(-1 * time.Hour)
		updated := time.Now()
		model := &models.ProductCategory{
			ID:          5,
			Name:        "Categoria Teste",
			Description: "Descrição da categoria",
			CreatedAt:   created,
			UpdatedAt:   updated,
		}

		dto := ToProductCategoryDTO(model)

		assert.NotNil(t, dto.ID)
		assert.Equal(t, model.ID, *dto.ID)
		assert.Equal(t, model.Name, dto.Name)
		assert.Equal(t, model.Description, dto.Description)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
		assert.Equal(t, model.UpdatedAt.Format(time.RFC3339), dto.UpdatedAt)
	})

	t.Run("Retorna DTO vazio se Model for nil", func(t *testing.T) {
		dto := ToProductCategoryDTO(nil)

		assert.NotNil(t, dto)
		assert.Nil(t, dto.ID)
		assert.Equal(t, "", dto.Name)
		assert.Equal(t, "", dto.Description)
		assert.Equal(t, "", dto.CreatedAt)
		assert.Equal(t, "", dto.UpdatedAt)
	})
}

func TestToProductCategoryDTOs(t *testing.T) {
	t.Run("Converte slice de Models para DTOs corretamente", func(t *testing.T) {
		now := time.Now()
		modelsInput := []*models.ProductCategory{
			{
				ID:          1,
				Name:        "Categoria A",
				Description: "Descrição A",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				Name:        "Categoria B",
				Description: "Descrição B",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		dtos := ToProductCategoryDTOs(modelsInput)

		assert.Len(t, dtos, 2)
		assert.Equal(t, modelsInput[0].ID, *dtos[0].ID)
		assert.Equal(t, modelsInput[0].Name, dtos[0].Name)
		assert.Equal(t, modelsInput[0].Description, dtos[0].Description)

		assert.Equal(t, modelsInput[1].ID, *dtos[1].ID)
		assert.Equal(t, modelsInput[1].Name, dtos[1].Name)
		assert.Equal(t, modelsInput[1].Description, dtos[1].Description)
	})

	t.Run("Retorna slice vazio quando lista de models é vazia", func(t *testing.T) {
		var modelsInput []*models.ProductCategory

		dtos := ToProductCategoryDTOs(modelsInput)

		assert.NotNil(t, dtos)
		assert.Empty(t, dtos)
	})

	t.Run("Ignora elementos nulos no slice", func(t *testing.T) {
		modelsInput := []*models.ProductCategory{
			nil,
			{
				ID:          3,
				Name:        "Categoria Válida",
				Description: "Descrição válida",
			},
			nil,
		}

		dtos := ToProductCategoryDTOs(modelsInput)

		assert.Len(t, dtos, 1)
		assert.Equal(t, modelsInput[1].Name, dtos[0].Name)
		assert.Equal(t, modelsInput[1].Description, dtos[0].Description)
	})
}
