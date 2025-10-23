package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relation"
)

type UserCategoryRelationsDTO struct {
	UserID     int64  `json:"user_id"`
	CategoryID int64  `json:"category_id"`
	CreatedAt  string `json:"created_at,omitempty"`
}

func ToUserCategoryRelationsModel(dto UserCategoryRelationsDTO) *models.UserCategoryRelation {
	return &models.UserCategoryRelation{
		UserID:     dto.UserID,
		CategoryID: dto.CategoryID,
		CreatedAt:  time.Now(),
	}
}

func ToUserCategoryRelationsDTO(m *models.UserCategoryRelation) UserCategoryRelationsDTO {
	if m == nil {
		return UserCategoryRelationsDTO{}
	}

	return UserCategoryRelationsDTO{
		UserID:     m.UserID,
		CategoryID: m.CategoryID,
		CreatedAt:  m.CreatedAt.Format(time.RFC3339),
	}
}
