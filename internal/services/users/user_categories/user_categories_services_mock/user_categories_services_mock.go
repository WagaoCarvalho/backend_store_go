package services

import (
	"context"

	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryService struct {
	mock.Mock
}

func (m *MockUserCategoryService) GetAll(ctx context.Context) ([]*models_user_categories.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models_user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) GetByID(ctx context.Context, id int64) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, id)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockUserCategoryService) Create(ctx context.Context, cat *models_user_categories.UserCategory) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) Update(ctx context.Context, cat *models_user_categories.UserCategory) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
