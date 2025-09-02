package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
)

// UserCategoryRelationsDTO representa o DTO de UserCategoryRelations
type UserCategoryRelationsDTO struct {
	UserID     int64  `json:"user_id"`
	CategoryID int64  `json:"category_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

// Converte DTO para Model
func ToUserCategoryRelationsModel(dto UserCategoryRelationsDTO) *models.UserCategoryRelations {
	return &models.UserCategoryRelations{
		UserID:     dto.UserID,
		CategoryID: dto.CategoryID,
		CreatedAt:  time.Now(), // ou use dto.CreatedAt parse se necessário
	}
}

// Converte Model para DTO
func ToUserCategoryRelationsDTO(m *models.UserCategoryRelations) UserCategoryRelationsDTO {
	if m == nil {
		return UserCategoryRelationsDTO{}
	}

	return UserCategoryRelationsDTO{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}
