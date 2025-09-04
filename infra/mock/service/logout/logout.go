package mockLogout

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogoutService struct {
	mock.Mock
}

func (m *MockLogoutService) Logout(ctx context.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}
