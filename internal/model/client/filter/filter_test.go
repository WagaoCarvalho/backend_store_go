package model

import (
	"testing"
	"time"

	errval "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestClientFilter_Validate(t *testing.T) {
	t.Run("valid filter with only base fields", func(t *testing.T) {
		f := ClientFilter{}
		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid limit from base filter", func(t *testing.T) {
		f := ClientFilter{}
		f.Limit = -10

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("valid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now().Add(-24 * time.Hour)
		to := time.Now()

		f := ClientFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid CreatedFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-24 * time.Hour)

		f := ClientFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.IsType(t, &errval.ValidationError{}, err)
		assert.Contains(t, err.Error(), "intervalo de criação inválido")
	})

	t.Run("invalid UpdatedFrom/To range", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-1 * time.Hour)

		f := ClientFilter{
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "intervalo de atualização inválido")
	})
}
