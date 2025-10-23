package model

import (
	"testing"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestContactSupplierRelations_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		csr := &SupplierContactRelation{
			ContactID:  1,
			SupplierID: 2,
			CreatedAt:  time.Now(),
		}

		err := csr.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid ContactID", func(t *testing.T) {
		csr := &SupplierContactRelation{
			ContactID:  0,
			SupplierID: 2,
		}

		err := csr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "contact_id", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})

	t.Run("invalid SupplierID", func(t *testing.T) {
		csr := &SupplierContactRelation{
			ContactID:  1,
			SupplierID: 0,
		}

		err := csr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "supplier_id", validationErr.Field)
		assert.Contains(t, validationErr.Message, "maior que zero")
	})

	t.Run("both ContactID and SupplierID invalid", func(t *testing.T) {
		csr := &SupplierContactRelation{
			ContactID:  0,
			SupplierID: 0,
		}

		err := csr.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*validators.ValidationError)
		assert.True(t, ok)

		assert.Equal(t, "contact_id", validationErr.Field)
	})
}
