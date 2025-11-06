package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
)

type UserContactRelationWriter interface {
	Create(ctx context.Context, relation *models.UserContactRelation) (*models.UserContactRelation, error)
	Delete(ctx context.Context, userID, contactID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type UserContactRelationReader interface {
	HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error)
}
