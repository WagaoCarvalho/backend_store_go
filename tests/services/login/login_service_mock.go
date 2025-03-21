package services

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockLoginService é um mock da interface LoginService para testes
type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}
