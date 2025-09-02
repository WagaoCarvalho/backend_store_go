package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type UserDTO struct {
	UID      *int64 `json:"uid,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Status   bool   `json:"status"`
	Version  int    `json:"version"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func ToUserModel(dto UserDTO) *models.User {
	return &models.User{
		UID:      utils.NilToZero(dto.UID),
		Username: dto.Username,
		Email:    dto.Email,
		Password: dto.Password,
		Status:   dto.Status,
		Version:  dto.Version,
	}
}

func ToUserDTO(m *models.User) UserDTO {
	if m == nil {
		return UserDTO{}
	}

	return UserDTO{
		UID:       &m.UID,
		Username:  m.Username,
		Email:     m.Email,
		Password:  m.Password,
		Status:    m.Status,
		Version:   m.Version,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
