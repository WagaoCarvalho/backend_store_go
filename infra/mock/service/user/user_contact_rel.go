package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
	"github.com/stretchr/testify/mock"
)

type MockUserContactRelationService struct {
	mock.Mock
}

func (m *MockUserContactRelationService) Create(ctx context.Context, userID, contactID int64) (*models.UserContactRelations, bool, error) {
	args := m.Called(ctx, userID, contactID)
	var result *models.UserContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).(*models.UserContactRelations)
	}
	return result, args.Bool(1), args.Error(2)
}

func (m *MockUserContactRelationService) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelations, error) {
	args := m.Called(ctx, userID)
	var result []*models.UserContactRelations
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.UserContactRelations)
	}
	return result, args.Error(1)
}

func (m *MockUserContactRelationService) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	args := m.Called(ctx, userID, contactID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserContactRelationService) Delete(ctx context.Context, userID, contactID int64) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func (m *MockUserContactRelationService) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
