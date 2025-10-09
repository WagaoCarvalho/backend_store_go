package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func TestModelUserContactRelations_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ucr := &UserContactRelations{
			UserID:    1,
			ContactID: 2,
		}

		err := ucr.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid UserID", func(t *testing.T) {
		ucr := &UserContactRelations{
			UserID:    0,
			ContactID: 2,
		}

		err := ucr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UserID", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})

	t.Run("invalid ContactID", func(t *testing.T) {
		ucr := &UserContactRelations{
			UserID:    1,
			ContactID: 0,
		}

		err := ucr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "ContactID", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})
}
