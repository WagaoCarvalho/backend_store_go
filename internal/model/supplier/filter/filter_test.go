package model

import (
	"testing"
	"time"

	errval "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestSupplierFilter_Validate(t *testing.T) {
	t.Run("valid filter with only base fields", func(t *testing.T) {
		f := SupplierFilter{}
		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid limit from base filter", func(t *testing.T) {
		f := SupplierFilter{}
		f.Limit = -10

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("invalid when CPF and CNPJ are provided together", func(t *testing.T) {
		f := SupplierFilter{
			CPF:  "123.456.789-00",
			CNPJ: "12.345.678/0001-00",
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CPF/CNPJ", ve.Field)
	})

	t.Run("invalid CPF format", func(t *testing.T) {
		f := SupplierFilter{
			CPF: "123",
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CPF", ve.Field)
	})

	t.Run("invalid CNPJ format", func(t *testing.T) {
		f := SupplierFilter{
			CNPJ: "123",
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CNPJ", ve.Field)
	})

	t.Run("invalid created date range", func(t *testing.T) {
		from := time.Now()
		to := from.Add(-time.Hour)

		f := SupplierFilter{
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CreatedFrom/CreatedTo", ve.Field)
	})

	t.Run("invalid updated date range", func(t *testing.T) {
		from := time.Now()
		to := from.Add(-time.Hour)

		f := SupplierFilter{
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UpdatedFrom/UpdatedTo", ve.Field)
	})

	t.Run("created from date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := SupplierFilter{
			CreatedFrom: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CreatedFrom", ve.Field)
	})

	t.Run("updated to date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := SupplierFilter{
			UpdatedTo: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UpdatedTo", ve.Field)
	})

	t.Run("valid filter with CPF only", func(t *testing.T) {
		f := SupplierFilter{
			CPF: "123.456.789-00",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with CNPJ only", func(t *testing.T) {
		f := SupplierFilter{
			CNPJ: "12.345.678/0001-00",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("created to date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := SupplierFilter{
			CreatedTo: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CreatedTo", ve.Field)
	})

	t.Run("updated from date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := SupplierFilter{
			UpdatedFrom: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UpdatedFrom", ve.Field)
	})

}
