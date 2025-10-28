package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockProductCategoryRelationRepo struct {
	mock.Mock
}

func (m *MockProductCategoryRelationRepo) Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelationRepo) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, tx, relation) // <-- aqui passa os 3 argumentos
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelationRepo) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]*models.ProductCategoryRelation), args.Error(1)
}

func (m *MockProductCategoryRelationRepo) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	args := m.Called(ctx, productID, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductCategoryRelationRepo) Delete(ctx context.Context, productID, categoryID int64) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockProductCategoryRelationRepo) DeleteAll(ctx context.Context, productID int64) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}
