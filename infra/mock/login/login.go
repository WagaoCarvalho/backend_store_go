package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, email, password string) (*models.AuthResponse, error) {
	args := m.Called(ctx, email, password)

	var authResp *models.AuthResponse
	if args.Get(0) != nil {
		authResp = args.Get(0).(*models.AuthResponse)
	}

	return authResp, args.Error(1)
}
