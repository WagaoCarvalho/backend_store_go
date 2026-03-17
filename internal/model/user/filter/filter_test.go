package model

import (
	"testing"
	"time"

	errval "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestUserFilter_Validate(t *testing.T) {
	t.Run("valid filter with only base fields", func(t *testing.T) {
		f := UserFilter{}
		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid limit from base filter", func(t *testing.T) {
		f := UserFilter{}
		f.Limit = -10

		err := f.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Limit")
	})

	t.Run("invalid email format", func(t *testing.T) {
		f := UserFilter{
			Email: "invalid-email",
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Email", ve.Field)
	})

	t.Run("valid email format", func(t *testing.T) {
		f := UserFilter{
			Email: "test@example.com",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("username too short", func(t *testing.T) {
		f := UserFilter{
			Username: "ab",
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "Username", ve.Field)
	})

	t.Run("valid username", func(t *testing.T) {
		f := UserFilter{
			Username: "john_doe",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with username only", func(t *testing.T) {
		f := UserFilter{
			Username: "john_doe",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with email only", func(t *testing.T) {
		f := UserFilter{
			Email: "john@example.com",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with status only", func(t *testing.T) {
		status := true
		f := UserFilter{
			Status: &status,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid filter with all fields", func(t *testing.T) {
		status := true
		now := time.Now()
		yesterday := now.Add(-24 * time.Hour)

		f := UserFilter{
			Username:    "john_doe",
			Email:       "john@example.com",
			Status:      &status,
			CreatedFrom: &yesterday,
			CreatedTo:   &now,
			UpdatedFrom: &yesterday,
			UpdatedTo:   &now,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid created date range", func(t *testing.T) {
		from := time.Now()
		to := from.Add(-time.Hour)

		f := UserFilter{
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

		f := UserFilter{
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

		f := UserFilter{
			CreatedFrom: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "CreatedFrom", ve.Field)
	})

	t.Run("created to date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := UserFilter{
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

		f := UserFilter{
			UpdatedFrom: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UpdatedFrom", ve.Field)
	})

	t.Run("updated to date in the future", func(t *testing.T) {
		future := time.Now().Add(24 * time.Hour)

		f := UserFilter{
			UpdatedTo: &future,
		}

		err := f.Validate()
		assert.Error(t, err)

		ve, ok := err.(*errval.ValidationError)
		assert.True(t, ok)
		assert.Equal(t, "UpdatedTo", ve.Field)
	})

	t.Run("created from and created to both in the past", func(t *testing.T) {
		past := time.Now().Add(-48 * time.Hour)
		morePast := time.Now().Add(-72 * time.Hour)

		f := UserFilter{
			CreatedFrom: &morePast,
			CreatedTo:   &past,
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("username with minimum valid length", func(t *testing.T) {
		f := UserFilter{
			Username: "abc",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("username with special characters", func(t *testing.T) {
		f := UserFilter{
			Username: "john.doe_123",
		}

		err := f.Validate()
		assert.NoError(t, err) // Apenas valida tamanho, não caracteres especiais
	})

	t.Run("email with subdomain", func(t *testing.T) {
		f := UserFilter{
			Email: "user@sub.example.com",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})

	t.Run("email with plus sign", func(t *testing.T) {
		f := UserFilter{
			Email: "user+label@example.com",
		}

		err := f.Validate()
		assert.NoError(t, err)
	})
}
