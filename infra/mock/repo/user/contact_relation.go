package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockUserContactRelationRepo struct {
	mock.Mock
}

func (m *MockUserContactRelationRepo) Create(ctx context.Context, relation *models.UserContactRelation) (*models.UserContactRelation, error) {
	args := m.Called(ctx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.UserContactRelation), args.Error(1)
}

func (m *MockUserContactRelationRepo) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserContactRelation) (*models.UserContactRelation, error) {
	args := m.Called(ctx, tx, relation)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.UserContactRelation), args.Error(1)
}

func (m *MockUserContactRelationRepo) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	args := m.Called(ctx, userID, contactID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserContactRelationRepo) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error) {
	args := m.Called(ctx, userID)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*models.UserContactRelation), args.Error(1)
}

func (m *MockUserContactRelationRepo) Delete(ctx context.Context, userID, contactID int64) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func (m *MockUserContactRelationRepo) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
