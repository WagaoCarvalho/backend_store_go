package services

import (
	"context"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]models_user.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]models_user.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) Create(
	ctx context.Context,
	user models_user.User,
	categoryID int64,
	address models_address.Address,
	contact models_contact.Contact,
) (models_user.User, error) {

	args := m.Called(ctx, user, categoryID, address, contact)
	if createdUser, ok := args.Get(0).(models_user.User); ok {
		return createdUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error) {
	args := m.Called(ctx, user, contact)
	if updatedUser, ok := args.Get(0).(models_user.User); ok {
		return updatedUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}
