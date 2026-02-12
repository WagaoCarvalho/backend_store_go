package model

import (
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryRelation_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}
		err := pcr.Validate()
		assert.NoError(t, err)
	})

	t.Run("ProductID zero", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  0,
			CategoryID: 2,
		}
		err := pcr.Validate()
		assert.Error(t, err)

		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "product_id", verr.Field)
		assert.Contains(t, verr.Message, "maior que zero")
	})

	t.Run("ProductID negativo", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  -1,
			CategoryID: 2,
		}
		err := pcr.Validate()
		assert.Error(t, err)

		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "product_id", verr.Field)
	})

	t.Run("CategoryID zero", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 0,
		}
		err := pcr.Validate()
		assert.Error(t, err)

		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "category_id", verr.Field)
	})

	t.Run("CategoryID negativo", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  1,
			CategoryID: -5,
		}
		err := pcr.Validate()
		assert.Error(t, err)

		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "category_id", verr.Field)
	})

	t.Run("Ambos inválidos - retorna primeiro erro", func(t *testing.T) {
		pcr := &ProductCategoryRelation{
			ProductID:  0,
			CategoryID: 0,
		}
		err := pcr.Validate()
		assert.Error(t, err)

		// Versão simplificada retorna apenas o primeiro erro
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "product_id", verr.Field)
	})
}
