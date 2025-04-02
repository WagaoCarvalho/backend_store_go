package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(models.User); ok {
		return user, args.Error(1)
	}
	return models.User{}, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(models.User); ok {
		return user, args.Error(1)
	}
	return models.User{}, args.Error(1)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user models.User, categoryID int64) (models.User, error) {
	args := m.Called(ctx, user, categoryID)
	if createdUser, ok := args.Get(0).(models.User); ok {
		return createdUser, args.Error(1)
	}
	return models.User{}, args.Error(1)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	if updatedUser, ok := args.Get(0).(models.User); ok {
		return updatedUser, args.Error(1)
	}
	return models.User{}, args.Error(1)
}

func (m *MockUserRepository) DeleteUserById(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
