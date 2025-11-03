package mockLogout

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLogout struct {
	mock.Mock
}

func (m *MockLogout) Logout(ctx context.Context, tokenString string) error {
	args := m.Called(ctx, tokenString)
	return args.Error(0)
}
