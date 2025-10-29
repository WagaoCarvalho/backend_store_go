package repo_test

import (
	"context"
	"errors"
	"testing"

	mock_repo "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSale_Create_Error(t *testing.T) {
	mockDB := new(mock_repo.MockDB)
	r := repo.NewSale(mockDB)

	ctx := context.Background()
	s := &model.Sale{}

	mockRow := new(MockRow)
	mockRow.On("Scan", mock.Anything).Return(errors.New("db error"))
	mockDB.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)

	res, err := r.Create(ctx, s)
	assert.Nil(t, res)
	assert.Error(t, err)
}

type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...any) error {
	args := m.Called(dest)
	return args.Error(0)
}
