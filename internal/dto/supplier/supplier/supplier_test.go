package dto

import (
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupplierDTO_ToSupplierModel(t *testing.T) {
	cnpj := "12345678000190"
	cpf := "12345678901"

	t.Run("With IsActive provided", func(t *testing.T) {
		isActive := true
		dto := SupplierDTO{
			ID:          utils.Int64Ptr(1),
			Name:        "Fornecedor X",
			CNPJ:        &cnpj,
			CPF:         &cpf,
			Description: "Descrição do fornecedor",
			IsActive:    &isActive,
		}

		model := ToSupplierModel(dto)

		assert.Equal(t, int64(1), model.ID)
		assert.Equal(t, "Fornecedor X", model.Name)
		assert.Equal(t, &cnpj, model.CNPJ)
		assert.Equal(t, &cpf, model.CPF)
		assert.Equal(t, "Descrição do fornecedor", model.Description)
		assert.True(t, model.Status)
	})

	t.Run("With IsActive not provided - should default to true", func(t *testing.T) {
		dto := SupplierDTO{
			ID:   utils.Int64Ptr(2),
			Name: "Fornecedor Y",
			CNPJ: &cnpj,
		}

		model := ToSupplierModel(dto)

		assert.Equal(t, int64(2), model.ID)
		assert.Equal(t, "Fornecedor Y", model.Name)
		assert.Equal(t, &cnpj, model.CNPJ)
		assert.True(t, model.Status, "Status should default to true when IsActive is nil")
	})

	t.Run("With IsActive set to false", func(t *testing.T) {
		isActive := false
		dto := SupplierDTO{
			ID:       utils.Int64Ptr(3),
			Name:     "Fornecedor Z",
			CPF:      &cpf,
			IsActive: &isActive,
		}

		model := ToSupplierModel(dto)

		assert.Equal(t, int64(3), model.ID)
		assert.Equal(t, "Fornecedor Z", model.Name)
		assert.Equal(t, &cpf, model.CPF)
		assert.False(t, model.Status)
	})
}

func TestToSupplierDTO(t *testing.T) {
	cnpj := "12345678000190"
	cpf := "12345678901"

	t.Run("Complete supplier with all fields", func(t *testing.T) {
		model := &models.Supplier{
			ID:          10,
			Name:        "Fornecedor Completo",
			CNPJ:        &cnpj,
			CPF:         nil,
			Description: "Descrição completa",
			Status:      true,
		}

		dto := ToSupplierDTO(model)

		assert.Equal(t, int64(10), *dto.ID)
		assert.Equal(t, "Fornecedor Completo", dto.Name)
		assert.Equal(t, &cnpj, dto.CNPJ)
		assert.Nil(t, dto.CPF)
		assert.Equal(t, "Descrição completa", dto.Description)
		assert.Equal(t, true, *dto.IsActive)
	})

	t.Run("Supplier with CPF and inactive status", func(t *testing.T) {
		model := &models.Supplier{
			ID:          11,
			Name:        "Fornecedor PF",
			CNPJ:        nil,
			CPF:         &cpf,
			Description: "",
			Status:      false,
		}

		dto := ToSupplierDTO(model)

		assert.Equal(t, int64(11), *dto.ID)
		assert.Equal(t, "Fornecedor PF", dto.Name)
		assert.Nil(t, dto.CNPJ)
		assert.Equal(t, &cpf, dto.CPF)
		assert.Equal(t, "", dto.Description)
		assert.Equal(t, false, *dto.IsActive)
	})

	t.Run("Supplier without CPF/CNPJ", func(t *testing.T) {
		model := &models.Supplier{
			ID:     12,
			Name:   "Fornecedor Sem Documento",
			Status: true,
		}

		dto := ToSupplierDTO(model)

		assert.Equal(t, int64(12), *dto.ID)
		assert.Equal(t, "Fornecedor Sem Documento", dto.Name)
		assert.Nil(t, dto.CNPJ)
		assert.Nil(t, dto.CPF)
		assert.Equal(t, true, *dto.IsActive)
	})
}

func TestToSupplierDTO_Nil(t *testing.T) {
	dto := ToSupplierDTO(nil)

	assert.Equal(t, SupplierDTO{}, dto)
	assert.Nil(t, dto.ID)
	assert.Nil(t, dto.IsActive)
	assert.Empty(t, dto.Name)
}

func TestToSupplierDTOs(t *testing.T) {
	cnpj := "12345678000190"
	cpf := "12345678901"

	t.Run("Multiple suppliers", func(t *testing.T) {
		models := []*models.Supplier{
			{
				ID:     1,
				Name:   "Fornecedor 1",
				CNPJ:   &cnpj,
				Status: true,
			},
			{
				ID:     2,
				Name:   "Fornecedor 2",
				CPF:    &cpf,
				Status: false,
			},
			{
				ID:     3,
				Name:   "Fornecedor 3",
				Status: true,
			},
		}

		dtos := ToSupplierDTOs(models)

		assert.Len(t, dtos, 3)

		assert.Equal(t, int64(1), *dtos[0].ID)
		assert.Equal(t, "Fornecedor 1", dtos[0].Name)
		assert.Equal(t, &cnpj, dtos[0].CNPJ)
		assert.Equal(t, true, *dtos[0].IsActive)

		assert.Equal(t, int64(2), *dtos[1].ID)
		assert.Equal(t, "Fornecedor 2", dtos[1].Name)
		assert.Equal(t, &cpf, dtos[1].CPF)
		assert.Equal(t, false, *dtos[1].IsActive)

		assert.Equal(t, int64(3), *dtos[2].ID)
		assert.Equal(t, "Fornecedor 3", dtos[2].Name)
		assert.Equal(t, true, *dtos[2].IsActive)
	})

	t.Run("Empty slice", func(t *testing.T) {
		dtos := ToSupplierDTOs([]*models.Supplier{})
		assert.Empty(t, dtos)
	})

	t.Run("Nil slice", func(t *testing.T) {
		dtos := ToSupplierDTOs(nil)
		assert.Empty(t, dtos)
	})

	t.Run("Slice with nil elements", func(t *testing.T) {
		models := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor 1",
			},
			nil,
			{
				ID:   2,
				Name: "Fornecedor 2",
			},
		}

		dtos := ToSupplierDTOs(models)

		// Deve pular os elementos nil
		assert.Len(t, dtos, 2)
		assert.Equal(t, int64(1), *dtos[0].ID)
		assert.Equal(t, int64(2), *dtos[1].ID)
	})
}

// Teste adicional para verificar consistência entre ToModel e ToDTO
func TestSupplierDTO_Consistency(t *testing.T) {
	originalDTO := SupplierDTO{
		ID:          utils.Int64Ptr(100),
		Name:        "Fornecedor Consistente",
		CNPJ:        utils.StrToPtr("12345678000190"),
		Description: "Teste de consistência",
		IsActive:    utils.BoolPtr(true),
	}

	// Converter DTO -> Model -> DTO
	model := ToSupplierModel(originalDTO)
	resultDTO := ToSupplierDTO(model)

	assert.Equal(t, *originalDTO.ID, *resultDTO.ID)
	assert.Equal(t, originalDTO.Name, resultDTO.Name)
	assert.Equal(t, originalDTO.CNPJ, resultDTO.CNPJ)
	assert.Equal(t, originalDTO.CPF, resultDTO.CPF)
	assert.Equal(t, originalDTO.Description, resultDTO.Description)
	assert.Equal(t, *originalDTO.IsActive, *resultDTO.IsActive)
}
