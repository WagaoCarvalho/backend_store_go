package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_ClientCpfExists(t *testing.T) {
	t.Run("successfully check client exists", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		// Usar bool diretamente
		mockRow := &mockDb.MockRow{Value: true}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		exists, err := repo.ClientCpfExists(ctx, 1)

		assert.NoError(t, err)
		assert.True(t, exists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrGet when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan fail")
		mockRow := &mockDb.MockRow{Err: scanErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{int64(1)}).Return(mockRow)

		exists, err := repo.ClientCpfExists(ctx, 1)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
	})
}
