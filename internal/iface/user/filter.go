package iface

import (
	"context"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	user "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
)

type UserFilter interface {
	Filter(ctx context.Context, f *filter.UserFilter) ([]*user.User, error)
}
