package models_test

import (
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func TestSupplierCategoryRelations_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		rel := &models.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 2,
		}
		err := rel.Validate()
		assert.Nil(t, err)
	})

	t.Run("missing supplier ID", func(t *testing.T) {
		rel := &models.SupplierCategoryRelations{
			SupplierID: 0,
			CategoryID: 2,
		}
		err := rel.Validate()
		assert.NotNil(t, err)
		verr, ok := err.(*utils_errors.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "SupplierID", verr.Field)
	})

	t.Run("missing category ID", func(t *testing.T) {
		rel := &models.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 0,
		}
		err := rel.Validate()
		assert.NotNil(t, err)
		verr, ok := err.(*utils_errors.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CategoryID", verr.Field)
	})
}
