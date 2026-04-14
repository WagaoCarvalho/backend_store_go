package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
)

type UserReader interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type UserWriter interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) error

	Delete(ctx context.Context, id int64) error
}

type UserStatus interface {
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
}

type UserVersion interface {
	GetVersionByID(ctx context.Context, id int64) (int64, error)
}
