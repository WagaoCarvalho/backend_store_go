package mock

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	"github.com/stretchr/testify/mock"
)

type MockUserCategory struct {
	mock.Mock
}

func (m *MockUserCategory) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.UserCategory), args.Error(1)
}

func (m *MockUserCategory) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.UserCategory), args.Error(1)
}

func (m *MockUserCategory) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	args := m.Called(ctx, category)

	var result *models.UserCategory
	if args.Get(0) != nil {
		result = args.Get(0).(*models.UserCategory)
	}
	return result, args.Error(1)
}

func (m *MockUserCategory) Update(ctx context.Context, category *models.UserCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockUserCategory) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
