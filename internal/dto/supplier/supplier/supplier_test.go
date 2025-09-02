package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupplierDTO_ToModel(t *testing.T) {
	cnpj := "12345678000190"
	cpf := "12345678901"

	dto := SupplierDTO{
		ID:      utils.Int64Ptr(1),
		Name:    "Fornecedor X",
		CNPJ:    &cnpj,
		CPF:     &cpf,
		Version: 2,
		Status:  true,
	}

	model := ToSupplierModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, "Fornecedor X", model.Name)
	assert.Equal(t, &cnpj, model.CNPJ)
	assert.Equal(t, &cpf, model.CPF)
	assert.Equal(t, 2, model.Version)
	assert.True(t, model.Status)
}

func TestToSupplierDTO(t *testing.T) {
	now := time.Now()
	cnpj := "12345678000190"

	model := models.Supplier{
		ID:        10,
		Name:      "Fornecedor Y",
		CNPJ:      &cnpj,
		Version:   1,
		Status:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	dto := ToSupplierDTO(&model)

	assert.Equal(t, int64(10), *dto.ID)
	assert.Equal(t, "Fornecedor Y", dto.Name)
	assert.Equal(t, &cnpj, dto.CNPJ)
	assert.Equal(t, 1, dto.Version)
	assert.False(t, dto.Status)
	assert.Equal(t, now.Format(time.RFC3339), dto.CreatedAt)
	assert.Equal(t, now.Format(time.RFC3339), dto.UpdatedAt)
}

func TestToSupplierDTO_Nil(t *testing.T) {
	dto := ToSupplierDTO(nil)

	assert.Equal(t, SupplierDTO{}, dto)
}
