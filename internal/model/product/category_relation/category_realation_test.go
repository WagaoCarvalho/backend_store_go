package model

import (
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestProductCategoryRelations_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		ucr := &ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}
		err := ucr.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid ProductID", func(t *testing.T) {
		ucr := &ProductCategoryRelation{
			ProductID:  0,
			CategoryID: 2,
		}
		err := ucr.Validate()
		assert.Error(t, err)
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "ProductID", verr.Field)
	})

	t.Run("Invalid CategoryID", func(t *testing.T) {
		ucr := &ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 0,
		}
		err := ucr.Validate()
		assert.Error(t, err)
		verr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CategoryID", verr.Field)
	})
}
