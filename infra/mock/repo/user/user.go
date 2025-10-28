package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if createdUser, ok := args.Get(0).(*models.User); ok { // Mudar para ponteiro
		return createdUser, args.Error(1)
	}
	return nil, args.Error(1) // Retornar nil em caso de erro
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]*models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return &models.User{}, args.Error(1)
}

func (m *MockUserRepository) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	args := m.Called(ctx, uid)

	// Garante segurança ao extrair o valor
	if version, ok := args.Get(0).(int64); ok {
		return version, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models.User); ok {
		return user, args.Error(1)
	}
	return &models.User{}, args.Error(1)
}

func (m *MockUserRepository) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	args := m.Called(ctx, name)
	if users, ok := args.Get(0).([]*models.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)

	return args.Error(0)
}

func (m *MockUserRepository) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}
