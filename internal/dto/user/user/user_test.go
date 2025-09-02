package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestToUserModel(t *testing.T) {
	t.Run("Convert UserDTO to User model", func(t *testing.T) {
		c := "user@example.com"
		dto := UserDTO{
			UID:      utils.Int64Ptr(1),
			Username: "usuario1",
			Email:    c,
			Password: "Senha123",
			Status:   true,
			Version:  2,
		}

		model := ToUserModel(dto)

		assert.Equal(t, int64(1), model.UID)
		assert.Equal(t, dto.Username, model.Username)
		assert.Equal(t, dto.Email, model.Email)
		assert.Equal(t, dto.Password, model.Password)
		assert.Equal(t, dto.Status, model.Status)
		assert.Equal(t, dto.Version, model.Version)
	})
}

func TestToUserDTO(t *testing.T) {
	t.Run("Convert User model to UserDTO", func(t *testing.T) {
		now := time.Now()
		model := &models.User{
			UID:       1,
			Username:  "usuario1",
			Email:     "user@example.com",
			Password:  "Senha123",
			Status:    true,
			Version:   2,
			CreatedAt: now,
			UpdatedAt: now,
		}

		dto := ToUserDTO(model)

		assert.Equal(t, &model.UID, dto.UID)
		assert.Equal(t, model.Username, dto.Username)
		assert.Equal(t, model.Email, dto.Email)
		assert.Equal(t, model.Password, dto.Password)
		assert.Equal(t, model.Status, dto.Status)
		assert.Equal(t, model.Version, dto.Version)
		assert.Equal(t, model.CreatedAt.Format(time.RFC3339), dto.CreatedAt)
		assert.Equal(t, model.UpdatedAt.Format(time.RFC3339), dto.UpdatedAt)
	})

	t.Run("Return empty DTO if model is nil", func(t *testing.T) {
		var model *models.User = nil
		dto := ToUserDTO(model)

		assert.Nil(t, dto.UID)
		assert.Equal(t, "", dto.Username)
		assert.Equal(t, "", dto.Email)
		assert.Equal(t, "", dto.Password)
		assert.False(t, dto.Status)
		assert.Equal(t, 0, dto.Version)
		assert.Equal(t, "", dto.CreatedAt)
		assert.Equal(t, "", dto.UpdatedAt)
	})
}
