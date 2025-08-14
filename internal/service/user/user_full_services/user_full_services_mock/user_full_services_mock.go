package services

import (
	"context"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço de usuário
type MockUserFullService struct {
	mock.Mock
}

func (m *MockUserFullService) CreateFull(ctx context.Context, user *model.UserFull) (*model.UserFull, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.UserFull), args.Error(1)
}
