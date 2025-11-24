package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_Create(t *testing.T) {
	t.Run("successfully create user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user description",
			Status:      true,
			Version:     1,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(1),                // UID (id)
				2,                       // version (incrementado)
				"Test user description", // description
				createdAt,               // created_at
				updatedAt,               // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user, result) // O mesmo objeto é retornado
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, 2, user.Version) // Version foi incrementado
		assert.Equal(t, "Test user description", user.Description)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user description",
			Status:      true,
			Version:     1,
		}

		dbError := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully create user with false status", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "inactiveuser",
			Email:       "inactive@example.com",
			Password:    "hashedpassword456",
			Description: "Inactive user description",
			Status:      false,
			Version:     0,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(2),                    // UID
				1,                           // version
				"Inactive user description", // description
				createdAt,                   // created_at
				updatedAt,                   // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(2), user.UID)
		assert.Equal(t, 1, user.Version)
		assert.Equal(t, false, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully create user with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "user3",
			Email:       "user3@example.com",
			Password:    "hashedpassword789",
			Description: "",
			Status:      true,
			Version:     1,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(3),  // UID
				2,         // version
				"",        // empty description
				createdAt, // created_at
				updatedAt, // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(3), user.UID)
		assert.Equal(t, 2, user.Version)
		assert.Equal(t, "", user.Description)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully create user with version zero", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "versionuser",
			Email:       "version@example.com",
			Password:    "hashedpassword000",
			Description: "Version test user",
			Status:      true,
			Version:     0,
		}

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(4),            // UID
				1,                   // version (incrementado de 0 para 1)
				"Version test user", // description
				createdAt,           // created_at
				updatedAt,           // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(4), user.UID)
		assert.Equal(t, 1, user.Version) // Version foi incrementado
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_Update(t *testing.T) {
	t.Run("successfully update user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         1,
			Username:    "updateduser",
			Email:       "updated@example.com",
			Description: "Updated description",
			Status:      true,
			Version:     1,
		}

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
				2,         // version (incrementado)
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 2, user.Version) // Version foi incrementado
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrVersionConflict when no rows affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         1,
			Username:    "updateduser",
			Email:       "updated@example.com",
			Description: "Updated description",
			Status:      true,
			Version:     1,
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         1,
			Username:    "updateduser",
			Email:       "updated@example.com",
			Description: "Updated description",
			Status:      true,
			Version:     1,
		}

		dbError := errors.New("database error")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully update user with false status", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         2,
			Username:    "inactiveuser",
			Email:       "inactive@example.com",
			Description: "Inactive user description",
			Status:      false,
			Version:     3,
		}

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
				4,         // version (incrementado)
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 4, user.Version)
		assert.Equal(t, false, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully update user with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         3,
			Username:    "emptydescuser",
			Email:       "empty@example.com",
			Description: "",
			Status:      true,
			Version:     2,
		}

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
				3,         // version (incrementado)
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 3, user.Version)
		assert.Equal(t, "", user.Description)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrVersionConflict with outdated version", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         1,
			Username:    "outdateduser",
			Email:       "outdated@example.com",
			Description: "Outdated description",
			Status:      true,
			Version:     1, // Versão desatualizada no banco
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
			user.Version,
		}).Return(mockRow)

		err := repo.Update(ctx, user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_Delete(t *testing.T) {
	t.Run("successfully delete user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		mockResult := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockResult := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		dbError := errors.New("database error")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(nil, dbError)

		err := repo.Delete(ctx, userID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		mockResult := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockResult := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete multiple users with same ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(5)

		mockResult := pgconn.NewCommandTag("DELETE 1") // Apenas 1 linha afetada mesmo com múltiplas tentativas
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
