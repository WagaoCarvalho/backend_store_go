package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
)

type MockClientCpf struct {
	mock.Mock
}

func (m *MockClientCpf) Create(ctx context.Context, clientCpf *models.ClientCpf) (*models.ClientCpf, error) {
	args := m.Called(ctx, clientCpf)
	var result *models.ClientCpf
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ClientCpf)
	}
	return result, args.Error(1)
}

func (m *MockClientCpf) GetByID(ctx context.Context, id int64) (*models.ClientCpf, error) {
	args := m.Called(ctx, id)
	var result *models.ClientCpf
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ClientCpf)
	}
	return result, args.Error(1)
}

func (m *MockClientCpf) GetByName(ctx context.Context, name string) ([]*models.ClientCpf, error) {
	args := m.Called(ctx, name)
	var result []*models.ClientCpf
	if args.Get(0) != nil {
		result = args.Get(0).([]*models.ClientCpf)
	}
	return result, args.Error(1)
}

func (m *MockClientCpf) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockClientCpf) Filter(ctx context.Context, f *filter.ClientCpfFilter) ([]*models.ClientCpf, error) {
	args := m.Called(ctx, f)

	var result []*models.ClientCpf
	if res := args.Get(0); res != nil {
		result = res.([]*models.ClientCpf)
	}
	return result, args.Error(1)
}

func (m *MockClientCpf) Update(ctx context.Context, client *models.ClientCpf) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

func (m *MockClientCpf) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientCpf) Disable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientCpf) Enable(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockClientCpf) ClientExists(ctx context.Context, clientID int64) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}

func (m *MockClientCpf) ClientCpfExists(ctx context.Context, clientID int64) (bool, error) {
	args := m.Called(ctx, clientID)
	return args.Bool(0), args.Error(1)
}
