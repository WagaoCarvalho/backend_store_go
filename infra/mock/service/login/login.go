package services

import (
	"context"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/login"
	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials dto.LoginCredentialsDTO) (*dto.AuthResponseDTO, error) {
	args := m.Called(ctx, credentials)
	authResp, _ := args.Get(0).(*dto.AuthResponseDTO)
	return authResp, args.Error(1)
}
