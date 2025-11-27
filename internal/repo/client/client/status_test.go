package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func TestClientRepo_Disable(t *testing.T) {
	t.Run("successfully disable client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				2,         // version (int)
				updatedAt, // updated_at (time.Time)
			},
		}

		query := `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		query := `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		query := `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDisable)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrDisable.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully disable client with zero ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(0)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1,         // version
				updatedAt, // updated_at
			},
		}

		query := `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative client ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		query := `
		UPDATE clients
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Disable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})
}

func TestClientRepo_Enable(t *testing.T) {
	t.Run("successfully enable client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				3,         // version (int)
				updatedAt, // updated_at (time.Time)
			},
		}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when client does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbErr := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbErr}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrEnable)
		assert.Contains(t, err.Error(), dbErr.Error())
		assert.Contains(t, err.Error(), errMsg.ErrEnable.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully enable client with zero ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(0)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				1,         // version
				updatedAt, // updated_at
			},
		}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative client ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully enable client with high version number", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(5)

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				100,       // version (número alto)
				updatedAt, // updated_at
			},
		}

		query := `
		UPDATE clients
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`
		mockDB.On("QueryRow", ctx, query, []interface{}{clientID}).Return(mockRow)

		err := repo.Enable(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
