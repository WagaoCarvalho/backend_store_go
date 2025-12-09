package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/client/filter"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Create(ctx context.Context, client *models.Client) (*models.Client, error) {
	args := m.Called(ctx, client)
	var result *models.Client
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClient) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	args := m.Called(ctx, id)
	var result *models.Client
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClient) GetByName(ctx context.Context, name string) ([]*models.Client, error) {
	args := m.Called(ctx, name)
	var result []*models.Client
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClient) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockClient) GetAll(ctx context.Context, f *filter.ClientFilter) ([]*models.Client, error) {
	args := m.Called(ctx, f)

	var result []*models.Client
	if res := args.Get(0); res != nil {
		result = res.([]*models.Client)
	}
	return result, args.Error(1)
}

func (m *MockClient) Update(ctx context.Context, client *models.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClient) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClient) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClient) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClient) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}
