package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryRelationRepo struct {
	mock.Mock
}

func (m *MockUserCategoryRelationRepo) Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	args := m.Called(ctx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	args := m.Called(ctx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.UserCategoryRelations), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	args := m.Called(ctx, userID, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserCategoryRelationRepo) Delete(ctx context.Context, userID, categoryID int64) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationRepo) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
