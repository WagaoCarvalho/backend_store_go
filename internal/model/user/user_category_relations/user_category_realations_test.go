package models

import (
	"testing"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserCategoryRelations_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		ucr := &UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}
		err := ucr.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid UserID", func(t *testing.T) {
		ucr := &UserCategoryRelations{
			UserID:     0,
			CategoryID: 2,
		}
		err := ucr.Validate()
		assert.Error(t, err)
		verr, ok := err.(*utils_errors.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UserID", verr.Field)
	})

	t.Run("Invalid CategoryID", func(t *testing.T) {
		ucr := &UserCategoryRelations{
			UserID:     1,
			CategoryID: 0,
		}
		err := ucr.Validate()
		assert.Error(t, err)
		verr, ok := err.(*utils_errors.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CategoryID", verr.Field)
	})
}
