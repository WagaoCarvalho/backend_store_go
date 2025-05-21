package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryRepository struct {
	mock.Mock
}

func (m *MockUserCategoryRepository) GetAll(ctx context.Context) ([]models.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) GetById(ctx context.Context, id int64) (models.UserCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) Create(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) Update(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(models.UserCategory), args.Error(1)
}

func (m *MockUserCategoryRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
