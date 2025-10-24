package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupplierCategoryRelationsDTO_ToModel(t *testing.T) {
	now := time.Now()

	dto := SupplierCategoryRelationsDTO{
		SupplierID: utils.Int64Ptr(100),
		CategoryID: utils.Int64Ptr(200),
		Version:    utils.IntPtr(1),
		CreatedAt:  &now,
	}

	model := ToSupplierCategoryRelationsModel(dto)

	assert.Equal(t, int64(100), model.SupplierID)
	assert.Equal(t, int64(200), model.CategoryID)
	assert.Equal(t, 1, model.Version)
	assert.Equal(t, now, model.CreatedAt)
}

func TestToSupplierCategoryRelationsDTO(t *testing.T) {
	now := time.Now()

	model := models.SupplierCategoryRelation{
		SupplierID: 101,
		CategoryID: 202,
		Version:    2,
		CreatedAt:  now,
	}

	dto := ToSupplierCategoryRelationsDTO(&model)

	assert.Equal(t, int64(101), *dto.SupplierID)
	assert.Equal(t, int64(202), *dto.CategoryID)
	assert.Equal(t, 2, *dto.Version)
	assert.Equal(t, now, *dto.CreatedAt)
}
