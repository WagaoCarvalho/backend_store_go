package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	"github.com/stretchr/testify/mock"
)

type MockClientService struct {
	mock.Mock
}

func (m *MockClientService) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	args := m.Called(ctx, client)
	if result := args.Get(0); result != nil {
		return result.(*models.Client), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClientService) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*models.Client), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClientService) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	args := m.Called(ctx, name)
	if result := args.Get(0); result != nil {
		return result.([]*models.Client), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClientService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockClientService) GetAll(ctx context.Context, limit, offset int) ([]*models.Client, error) {
	args := m.Called(ctx, limit, offset)
	if result := args.Get(0); result != nil {
		return result.([]*models.Client), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockClientService) Update(ctx context.Context, client *models.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientService) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientService) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientService) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}
