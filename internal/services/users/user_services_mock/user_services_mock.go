package services

import (
	"context"

	models_user "github.com/WagaoCarvalho/backend_store_go/internal/model/user"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço de usuário
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models_user.User), args.Error(1)
}

func (m *MockUserService) CreateFull(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models_user.User), args.Error(1)
}

func (m *MockUserService) GetAll(ctx context.Context) ([]*models_user.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models_user.User), args.Error(1)
}

func (m *MockUserService) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	args := m.Called(ctx, uid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models_user.User), args.Error(1)
}

func (m *MockUserService) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models_user.User), args.Error(1)
}

func (m *MockUserService) GetByName(ctx context.Context, name string) ([]*models_user.User, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models_user.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models_user.User), args.Error(1)
}

func (m *MockUserService) Disable(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserService) Enable(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserService) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
