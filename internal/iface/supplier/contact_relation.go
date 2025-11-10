package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
)

type SupplierContactRelationReader interface {
	HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error)
	GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error)
}

type SupplierContactRelationWriter interface {
	Create(ctx context.Context, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error)
	Delete(ctx context.Context, supplierID, contactID int64) error
	DeleteAll(ctx context.Context, supplierID int64) error
}
