package model

import (
	"testing"
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestClientCpfFilter_Validate(t *testing.T) {

	t.Run("valid filter with no dates", func(t *testing.T) {
		f := &ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with valid created date range", func(t *testing.T) {
		from := time.Now().Add(-24 * time.Hour)
		to := time.Now()

		f := &ClientCpfFilter{
			BaseFilter:  filter.BaseFilter{Limit: 10},
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with valid updated date range", func(t *testing.T) {
		from := time.Now().Add(-48 * time.Hour)
		to := time.Now().Add(-24 * time.Hour)

		f := &ClientCpfFilter{
			BaseFilter:  filter.BaseFilter{Limit: 10},
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("return error when CreatedFrom is after CreatedTo", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-1 * time.Hour)

		f := &ClientCpfFilter{
			BaseFilter:  filter.BaseFilter{Limit: 10},
			CreatedFrom: &from,
			CreatedTo:   &to,
		}

		err := f.Validate()

		assert.Error(t, err)

		var vErr *validators.ValidationError
		assert.ErrorAs(t, err, &vErr)
		assert.Equal(t, "CreatedFrom/CreatedTo", vErr.Field)
	})

	t.Run("return error when UpdatedFrom is after UpdatedTo", func(t *testing.T) {
		from := time.Now()
		to := time.Now().Add(-2 * time.Hour)

		f := &ClientCpfFilter{
			BaseFilter:  filter.BaseFilter{Limit: 10},
			UpdatedFrom: &from,
			UpdatedTo:   &to,
		}

		err := f.Validate()

		assert.Error(t, err)

		var vErr *validators.ValidationError
		assert.ErrorAs(t, err, &vErr)
		assert.Equal(t, "UpdatedFrom/UpdatedTo", vErr.Field)
	})

	t.Run("return error when BaseFilter.Validate fails", func(t *testing.T) {
		f := &ClientCpfFilter{
			BaseFilter: filter.BaseFilter{
				Limit: -1, // inválido
			},
		}

		err := f.Validate()
		assert.Error(t, err)
	})
}
