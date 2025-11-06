package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
)

type UserWriter interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) error

	Delete(ctx context.Context, id int64) error
}

type UserReader interface {
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByName(ctx context.Context, name string) ([]*models.User, error)
	UserExists(ctx context.Context, userID int64) (bool, error)
}

type UserStatus interface {
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
}
