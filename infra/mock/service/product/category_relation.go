package services

import (
	"context"

	product_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/stretchr/testify/mock"
)

type MockProductCategoryRelationService struct {
	mock.Mock
}

func (m *MockProductCategoryRelationService) Create(ctx context.Context, productID, categoryID int64) (*product_category_relations.ProductCategoryRelation, bool, error) {
	args := m.Called(ctx, productID, categoryID)

	var relation *product_category_relations.ProductCategoryRelation
	if rel, ok := args.Get(0).(*product_category_relations.ProductCategoryRelation); ok {
		relation = rel
	}

	created := false
	if val, ok := args.Get(1).(bool); ok {
		created = val
	}

	return relation, created, args.Error(2)
}

func (m *MockProductCategoryRelationService) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*product_category_relations.ProductCategoryRelation, error) {
	args := m.Called(ctx, productID)
	if rels, ok := args.Get(0).([]*product_category_relations.ProductCategoryRelation); ok {
		return rels, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductCategoryRelationService) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	args := m.Called(ctx, productID, categoryID)
	if exists, ok := args.Get(0).(bool); ok {
		return exists, args.Error(1)
	}
	return false, args.Error(1)
}

func (m *MockProductCategoryRelationService) Update(ctx context.Context, relation *product_category_relations.ProductCategoryRelation) (*product_category_relations.ProductCategoryRelation, error) {
	args := m.Called(ctx, relation)
	if updated, ok := args.Get(0).(*product_category_relations.ProductCategoryRelation); ok {
		return updated, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductCategoryRelationService) Delete(ctx context.Context, productID, categoryID int64) error {
	args := m.Called(ctx, productID, categoryID)
	return args.Error(0)
}

func (m *MockProductCategoryRelationService) DeleteAll(ctx context.Context, productID int64) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}
