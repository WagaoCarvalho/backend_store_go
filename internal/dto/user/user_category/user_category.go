package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category"
)

type UserCategoryDTO struct {
	ID          *uint  `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// Converte DTO para Model
func ToUserCategoryModel(dto UserCategoryDTO) *models.UserCategory {
	id := uint(0)
	if dto.ID != nil {
		id = *dto.ID
	}

	return &models.UserCategory{
		ID:          id,
		Name:        dto.Name,
		Description: dto.Description,
	}
}

// Converte Model para DTO
func ToUserCategoryDTO(m *models.UserCategory) UserCategoryDTO {
	if m == nil {
		return UserCategoryDTO{}
	}

	return UserCategoryDTO{
		ID:          &m.ID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

func ToUserCategoryDTOs(models []*models.UserCategory) []UserCategoryDTO {
	if len(models) == 0 {
		return []UserCategoryDTO{}
	}

	dtos := make([]UserCategoryDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToUserCategoryDTO(m))
		}
	}
	return dtos
}
