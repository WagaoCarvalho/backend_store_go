package services

import (
	"context"

	user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryRelationService struct {
	mock.Mock
}

func (m *MockUserCategoryRelationService) Create(ctx context.Context, userID, categoryID int64) (*user_category_relations.UserCategoryRelations, bool, error) {
	args := m.Called(ctx, userID, categoryID)

	var relation *user_category_relations.UserCategoryRelations
	if rel, ok := args.Get(0).(*user_category_relations.UserCategoryRelations); ok {
		relation = rel
	}

	created := false
	if val, ok := args.Get(1).(bool); ok {
		created = val
	}

	return relation, created, args.Error(2)
}

func (m *MockUserCategoryRelationService) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, userID)
	if rels, ok := args.Get(0).([]*user_category_relations.UserCategoryRelations); ok {
		return rels, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryRelationService) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	args := m.Called(ctx, userID, categoryID)
	if exists, ok := args.Get(0).(bool); ok {
		return exists, args.Error(1)
	}
	return false, args.Error(1)
}

func (m *MockUserCategoryRelationService) Update(ctx context.Context, relation *user_category_relations.UserCategoryRelations) (*user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, relation)
	if updated, ok := args.Get(0).(*user_category_relations.UserCategoryRelations); ok {
		return updated, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryRelationService) Delete(ctx context.Context, userID, categoryID int64) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationService) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
