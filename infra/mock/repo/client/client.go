package repositories

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
)

type MockClientRepository struct {
	mock.Mock
}

func (m *MockClientRepository) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	args := m.Called(ctx, client)
	var result *models.Client
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClientRepository) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	args := m.Called(ctx, id)
	var result *models.Client
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClientRepository) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	args := m.Called(ctx, name)
	var result []*models.Client
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClientRepository) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockClientRepository) GetAll(ctx context.Context) ([]*models.Client, error) {
	args := m.Called(ctx)
	var result []*models.Client
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClientRepository) Update(ctx context.Context, client *models.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientRepository) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientRepository) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientRepository) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}
