package repositories

import (
	"context"

	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// MockUserRepository
func (m *MockUserRepository) Create(
	ctx context.Context,
	user *models_user.User,
) (*models_user.User, error) {
	args := m.Called(ctx, user)
	if createdUser, ok := args.Get(0).(*models_user.User); ok { // Mudar para ponteiro
		return createdUser, args.Error(1)
	}
	return nil, args.Error(1) // Retornar nil em caso de erro
}

func (m *MockUserRepository) Update(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	args := m.Called(ctx, user)

	var usr *models_user.User
	if val := args.Get(0); val != nil {
		usr = val.(*models_user.User)
	}

	return usr, args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models_user.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]*models_user.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(*models_user.User); ok {
		return user, args.Error(1)
	}
	return &models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) GetVersionByID(ctx context.Context, uid int64) (int64, error) {
	args := m.Called(ctx, uid)

	// Garante segurança ao extrair o valor
	if version, ok := args.Get(0).(int64); ok {
		return version, args.Error(1)
	}

	return 0, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(*models_user.User); ok {
		return user, args.Error(1)
	}
	return &models_user.User{}, args.Error(1)
}
