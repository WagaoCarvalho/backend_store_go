package model

import (
	"testing"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestClientContactRelations_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ccr := &ClientContactRelation{
			ClientID:  1,
			ContactID: 2,
			CreatedAt: time.Now(),
		}

		err := ccr.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid ClientID", func(t *testing.T) {
		ccr := &ClientContactRelation{
			ClientID:  0,
			ContactID: 2,
		}

		err := ccr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "client_id", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})

	t.Run("invalid ContactID", func(t *testing.T) {
		ccr := &ClientContactRelation{
			ClientID:  1,
			ContactID: 0,
		}

		err := ccr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "contact_id", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})

	t.Run("both ClientID and ContactID invalid", func(t *testing.T) {
		ccr := &ClientContactRelation{
			ClientID:  0,
			ContactID: 0,
		}

		err := ccr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		// Retorna apenas o primeiro erro encontrado
		assert.Equal(t, "client_id", validationErr.Field)
	})
}
