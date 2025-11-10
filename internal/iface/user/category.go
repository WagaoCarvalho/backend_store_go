package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
)

type UserCategoryReader interface {
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
}

type UserCategoryWriter interface {
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) error
	Delete(ctx context.Context, id int64) error
}
