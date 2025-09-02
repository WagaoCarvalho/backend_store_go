package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupplierCategoryDTO_ToModel(t *testing.T) {
	now := time.Now()

	dto := SupplierCategoryDTO{
		ID:          utils.Int64Ptr(1),
		Name:        "Categoria X",
		Description: "Descrição da categoria",
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	model := ToSupplierCategoryModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, "Categoria X", model.Name)
	assert.Equal(t, "Descrição da categoria", model.Description)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
}

func TestToSupplierCategoryDTO(t *testing.T) {
	now := time.Now()

	model := models.SupplierCategory{
		ID:          10,
		Name:        "Categoria Y",
		Description: "Outra descrição",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	dto := ToSupplierCategoryDTO(&model)

	assert.Equal(t, int64(10), *dto.ID)
	assert.Equal(t, "Categoria Y", dto.Name)
	assert.Equal(t, "Outra descrição", dto.Description)
	assert.Equal(t, now, *dto.CreatedAt)
	assert.Equal(t, now, *dto.UpdatedAt)
}
