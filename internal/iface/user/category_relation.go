package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
)

type UserCategoryRelationWriter interface {
	Create(ctx context.Context, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type UserCategoryRelationReader interface {
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error)
}
