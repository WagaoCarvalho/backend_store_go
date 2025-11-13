package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductRepo_EnableProduct(t *testing.T) {
	t.Run("successfully enable product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		// Usando pgconn.NewCommandTag para criar o CommandTag
		mockResult := pgconn.NewCommandTag("UPDATE 1")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(mockResult, nil)

		err := repo.EnableProduct(ctx, productID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		// CommandTag com 0 rows affected
		mockResult := pgconn.NewCommandTag("UPDATE 0")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(mockResult, nil)

		err := repo.EnableProduct(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database exec fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		execErr := errors.New("database error")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(nil, execErr)

		err := repo.EnableProduct(ctx, productID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar")
		mockDB.AssertExpectations(t)
	})
}

func TestProductRepo_DisableProduct(t *testing.T) {
	t.Run("successfully disable product", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		mockResult := pgconn.NewCommandTag("UPDATE 1")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(mockResult, nil)

		err := repo.DisableProduct(ctx, productID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when product does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(999)

		mockResult := pgconn.NewCommandTag("UPDATE 0")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(mockResult, nil)

		err := repo.DisableProduct(ctx, productID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database exec fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productRepo{db: mockDB}
		ctx := context.Background()
		productID := int64(1)

		execErr := errors.New("database error")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{productID}).Return(nil, execErr)

		err := repo.DisableProduct(ctx, productID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar")
		mockDB.AssertExpectations(t)
	})
}
