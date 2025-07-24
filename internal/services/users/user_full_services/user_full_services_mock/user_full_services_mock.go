package services

import (
	"context"

	model_user_full "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_full"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço de usuário
type MockUserFullService struct {
	mock.Mock
}

func (m *MockUserFullService) CreateFull(ctx context.Context, user *model_user_full.UserFull) (*model_user_full.UserFull, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model_user_full.UserFull), args.Error(1)
}
