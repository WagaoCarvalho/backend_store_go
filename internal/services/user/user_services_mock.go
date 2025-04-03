package services

import (
	"context"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models_user.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]models_user.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, uid int64) (models_user.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models_user.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user models_user.User, categoryID int64, address models_address.Address) (models_user.User, error) {
	args := m.Called(ctx, user, categoryID, address)
	if createdUser, ok := args.Get(0).(models_user.User); ok {
		return createdUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user models_user.User) (models_user.User, error) {
	args := m.Called(ctx, user)
	if updatedUser, ok := args.Get(0).(models_user.User); ok {
		return updatedUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) DeleteUserById(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
