package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockContact struct {
	mock.Mock
}

func (m *MockContact) Create(ctx context.Context, c *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, c)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Contact), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContact) CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, tx, contact)

	var result *models.Contact
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Contact)
	}

	return result, args.Error(1)
}

func (m *MockContact) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContact) Update(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContact) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
