package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewUserFull(t *testing.T) {
	t.Run("successfully create new UserFull instance", func(t *testing.T) {
		var db repo.DBTransactor

		result := NewUserFull(db)

		assert.NotNil(t, result)

		_, ok := result.(*userFullRepo)
		assert.True(t, ok, "Expected result to be of type *userFullRepo")
	})

	t.Run("return instance with provided db transactor", func(t *testing.T) {
		var db repo.DBTransactor

		result := NewUserFull(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db repo.DBTransactor

		instance1 := NewUserFull(db)
		instance2 := NewUserFull(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestUserFullRepo_BeginTx(t *testing.T) {
	t.Run("successfully begin transaction", func(t *testing.T) {
		mockDB := new(mockDb.MockDBTransactor)
		repo := &userFullRepo{db: mockDB}
		ctx := context.Background()

		mockTx := new(mockDb.MockTx)
		mockDB.On("BeginTx", ctx, pgx.TxOptions{}).Return(mockTx, nil)

		tx, err := repo.BeginTx(ctx)

		assert.NoError(t, err)
		assert.Equal(t, mockTx, tx)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when begin transaction fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDBTransactor)
		repo := &userFullRepo{db: mockDB}
		ctx := context.Background()

		dbError := errors.New("transaction failed")
		// Retornar MockTx vazio junto com o erro
		mockDB.On("BeginTx", ctx, pgx.TxOptions{}).Return(&mockDb.MockTx{}, dbError)

		tx, err := repo.BeginTx(ctx)

		assert.NotNil(t, tx) // Agora retorna um MockTx (mesmo que vazio)
		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		mockDB.AssertExpectations(t)
	})
}

func TestUserFullRepo_CreateTx(t *testing.T) {
	t.Run("successfully create user within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &userFullRepo{db: nil} // db não é usado no CreateTx
		ctx := context.Background()

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
			Status:   true,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []interface{}{
				int64(1),  // UID
				createdAt, // created_at
				updatedAt, // updated_at
			},
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Password,
			user.Status,
		}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user, result)
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &userFullRepo{db: nil}
		ctx := context.Background()

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
			Status:   true,
		}

		dbError := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbError}
		mockTx.On("QueryRow", ctx, mock.Anything, mock.Anything).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockTx.AssertExpectations(t)
	})

	t.Run("successfully create user with false status", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &userFullRepo{db: nil}
		ctx := context.Background()

		user := &models.User{
			Username: "inactiveuser",
			Email:    "inactive@example.com",
			Password: "Password123",
			Status:   false,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []interface{}{
				int64(2),  // UID
				createdAt, // created_at
				updatedAt, // updated_at
			},
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Password,
			user.Status,
		}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user, result)
		assert.Equal(t, int64(2), user.UID)
		assert.Equal(t, false, user.Status)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when scan fails due to wrong number of values", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &userFullRepo{db: nil}
		ctx := context.Background()

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
			Status:   true,
		}

		// MockRow com número incorreto de valores (apenas 2 em vez de 3)
		mockRow := &mockDb.MockRowWithIDArgs{
			Values: []interface{}{
				int64(1), // UID
				// Faltam created_at e updated_at
			},
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Password,
			user.Status,
		}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when scan fails with custom error", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &userFullRepo{db: nil}
		ctx := context.Background()

		user := &models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
			Status:   true,
		}

		// Usando MockRow com erro customizado
		scanError := errors.New("scan failed: invalid data format")
		mockRow := &mockDb.MockRowWithIDArgs{
			Err: scanError,
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Password,
			user.Status,
		}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), scanError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockTx.AssertExpectations(t)
	})
}
