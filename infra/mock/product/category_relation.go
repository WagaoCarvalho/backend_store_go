package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockProductCategoryRelation struct {
	mock.Mock
}

func (m *MockProductCategoryRelation) Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, relation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, tx, relation) // <-- aqui passa os 3 argumentos
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelation) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelation) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	args := m.Called(ctx, productID, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductCategoryRelation) Delete(ctx context.Context, productID, categoryID int64) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockProductCategoryRelation) DeleteAll(ctx context.Context, productID int64) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}
