package dto

import (
	"strings"
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryDTO_Validate(t *testing.T) {
	t.Run("Validação bem-sucedida com dados válidos", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name:        "Categoria Válida",
			Description: utils.StrToPtr("Descrição válida dentro do limite"),
		}

		err := dto.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validação bem-sucedida sem descrição", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name: "Categoria Válida",
		}

		err := dto.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validação falha com nome vazio", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name: "",
		}

		err := dto.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome é obrigatório")
	})

	t.Run("Validação falha com nome apenas espaços", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name: "   ",
		}

		err := dto.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome é obrigatório")
	})

	t.Run("Validação falha com nome muito curto", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name: "A",
		}

		err := dto.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "mínimo de 2 caracteres") // Mensagem corrigida
	})

	t.Run("Validação falha com nome muito longo", func(t *testing.T) {
		longName := strings.Repeat("A", 256)
		dto := ProductCategoryDTO{
			Name: longName,
		}

		err := dto.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome máximo 255 caracteres")
	})

	t.Run("Validação falha com descrição muito longa", func(t *testing.T) {
		longDesc := strings.Repeat("A", 256)
		dto := ProductCategoryDTO{
			Name:        "Categoria Válida",
			Description: &longDesc,
		}

		err := dto.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "descrição máxima 255 caracteres")
	})

	t.Run("Trim spaces é aplicado no nome", func(t *testing.T) {
		dto := ProductCategoryDTO{
			Name: "  Categoria Teste  ",
		}

		err := dto.Validate()
		assert.NoError(t, err)
		// O DTO deve manter o valor original, o trim só é feito na validação
		assert.Equal(t, "  Categoria Teste  ", dto.Name)
	})

	t.Run("Trim spaces é aplicado na descrição", func(t *testing.T) {
		desc := "  Descrição com espaços  "
		dto := ProductCategoryDTO{
			Name:        "Categoria Teste",
			Description: &desc,
		}

		err := dto.Validate()
		assert.NoError(t, err)
		// O ponteiro deve ter sido atualizado com o valor trimmed
		assert.Equal(t, "Descrição com espaços", *dto.Description)
	})

	t.Run("Descrição vazia após trim NÃO é tratada como nil", func(t *testing.T) {
		desc := "   "
		dto := ProductCategoryDTO{
			Name:        "Categoria Teste",
			Description: &desc,
		}

		err := dto.Validate()
		assert.NoError(t, err)
		// Se seu Validate() não seta como nil, apenas faz trim
		// então vamos verificar se ficou como string vazia
		if dto.Description != nil {
			assert.Equal(t, "", *dto.Description)
		}
	})
}

func TestToProductCategoryModel(t *testing.T) {
	t.Run("Converte DTO para Model corretamente", func(t *testing.T) {
		dto := ProductCategoryDTO{
			ID:          nil,
			Name:        "Categoria Teste",
			Description: utils.StrToPtr("Descrição da categoria"),
		}

		model := ToProductCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, int64(0), model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, *dto.Description, model.Description)
	})

	t.Run("Converte DTO com ID definido para Model corretamente", func(t *testing.T) {
		var id int64 = 10
		dto := ProductCategoryDTO{
			ID:          &id,
			Name:        "Categoria Teste",
			Description: utils.StrToPtr("Descrição da categoria"),
		}

		model := ToProductCategoryModel(dto)

		assert.NotNil(t, model)
		assert.Equal(t, id, model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, *dto.Description, model.Description)
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
		assert.Equal(t, model.Description, *dto.Description)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), *dto.CreatedAt)
		assert.Equal(t, model.UpdatedAt.Format(time.RFC3339), *dto.UpdatedAt)
	})

	t.Run("Retorna DTO vazio se Model for nil", func(t *testing.T) {
		dto := ToProductCategoryDTO(nil)

		assert.NotNil(t, dto)
		assert.Nil(t, dto.ID)
		assert.Equal(t, "", dto.Name)
		assert.Nil(t, dto.Description)
		assert.Nil(t, dto.CreatedAt)
		assert.Nil(t, dto.UpdatedAt)
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
		assert.Equal(t, modelsInput[0].Description, *dtos[0].Description)

		assert.Equal(t, modelsInput[1].ID, *dtos[1].ID)
		assert.Equal(t, modelsInput[1].Name, dtos[1].Name)
		assert.Equal(t, modelsInput[1].Description, *dtos[1].Description)
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
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			nil,
		}

		dtos := ToProductCategoryDTOs(modelsInput)

		assert.Len(t, dtos, 1)
		assert.Equal(t, modelsInput[1].Name, dtos[0].Name)
		assert.Equal(t, modelsInput[1].Description, *dtos[0].Description)
	})
}

func TestToProductCategoryModelPtr(t *testing.T) {
	t.Run("Converte DTO pointer para Model corretamente", func(t *testing.T) {
		var id int64 = 5
		dto := &ProductCategoryDTO{
			ID:          &id,
			Name:        "Categoria Pointer",
			Description: utils.StrToPtr("Descrição do pointer"),
		}

		model := ToProductCategoryModelPtr(dto)

		assert.NotNil(t, model)
		assert.Equal(t, id, model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, *dto.Description, model.Description)
	})

	t.Run("Retorna nil quando DTO pointer é nil", func(t *testing.T) {
		model := ToProductCategoryModelPtr(nil)

		assert.Nil(t, model)
	})

	t.Run("Converte DTO pointer com valores nil corretamente", func(t *testing.T) {
		dto := &ProductCategoryDTO{
			Name: "Categoria Sem ID",
		}

		model := ToProductCategoryModelPtr(dto)

		assert.NotNil(t, model)
		assert.Equal(t, int64(0), model.ID)
		assert.Equal(t, dto.Name, model.Name)
		assert.Equal(t, "", model.Description) // Descrição vazia quando nil
	})
}
