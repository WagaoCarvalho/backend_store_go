package handlers

import (
	"context"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockLoginService representa um mock do LoginService
type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials models.LoginCredentials) (string, error) {
	args := m.Called(ctx, credentials)
	return args.String(0), args.Error(1)
}
