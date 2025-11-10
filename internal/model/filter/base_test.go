package model

import (
	"testing"

	validator "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestBaseFilter_Validate(t *testing.T) {
	t.Run("valid filter", func(t *testing.T) {
		f := BaseFilter{
			Limit:     100,
			Offset:    0,
			SortBy:    "name",
			SortOrder: "ASC",
		}

		err := f.Validate()
		assert.NoError(t, err)
		assert.Equal(t, "asc", f.SortOrder)
	})

	t.Run("negative limit", func(t *testing.T) {
		f := BaseFilter{Limit: -1}
		err := f.Validate()

		assert.Error(t, err)
		assert.IsType(t, &validator.ValidationError{}, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("limit above maximum", func(t *testing.T) {
		f := BaseFilter{Limit: 2000}
		err := f.Validate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "máximo permitido é 1000")
	})

	t.Run("negative offset", func(t *testing.T) {
		f := BaseFilter{Offset: -5}
		err := f.Validate()

		assert.Error(t, err)
		assert.IsType(t, &validator.ValidationError{}, err)
		assert.Contains(t, err.Error(), "Offset")
	})

	t.Run("invalid sort order", func(t *testing.T) {
		f := BaseFilter{SortOrder: "ascending"}
		err := f.Validate()

		assert.Error(t, err)
		assert.IsType(t, &validator.ValidationError{}, err)
		assert.Contains(t, err.Error(), "SortOrder")
	})
}
