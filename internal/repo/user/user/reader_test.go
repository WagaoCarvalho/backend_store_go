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
	"github.com/stretchr/testify/mock"
)

func TestUserRepo_GetByID(t *testing.T) {
	t.Run("successfully get user by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(1),           // UID (id na query)
				"testuser",         // username
				"test@example.com", // email
				"Test description", // description
				true,               // status
				createdAt,          // created_at
				updatedAt,          // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "Test description", user.Description)
		assert.Equal(t, true, user.Status)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, user)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "erro ao buscar")
		assert.Contains(t, err.Error(), errMsg.ErrGet.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with false status", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(2)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(2),                    // UID
				"inactiveuser",              // username
				"inactive@example.com",      // email
				"Inactive user description", // description
				false,                       // status
				createdAt,                   // created_at
				updatedAt,                   // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(2), user.UID)
		assert.Equal(t, "inactiveuser", user.Username)
		assert.Equal(t, "inactive@example.com", user.Email)
		assert.Equal(t, "Inactive user description", user.Description)
		assert.Equal(t, false, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(3)

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(3),            // UID
				"user3",             // username
				"user3@example.com", // email
				"",                  // empty description
				true,                // status
				createdAt,           // created_at
				updatedAt,           // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		user, err := repo.GetByID(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(3), user.UID)
		assert.Equal(t, "user3", user.Username)
		assert.Equal(t, "user3@example.com", user.Email)
		assert.Equal(t, "", user.Description)
		assert.Equal(t, true, user.Status)
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)
		expectedVersion := int64(5)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(999)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(1)

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrGetVersion)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get version zero", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(2)
		expectedVersion := int64(0)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get version with high number", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(3)
		expectedVersion := int64(100)

		mockRow := &mockDb.MockRow{
			Value: expectedVersion,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedVersion, version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(0)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		userID := int64(-1)

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{userID}).Return(mockRow)

		version, err := repo.GetVersionByID(ctx, userID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_GetAll(t *testing.T) {
	t.Run("successfully get all users", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		// Primeiro usuário
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &user.UID
			mock.AnythingOfType("*string"),    // &user.Username
			mock.AnythingOfType("*string"),    // &user.Email
			mock.AnythingOfType("*string"),    // &user.Description
			mock.AnythingOfType("*bool"),      // &user.Status
			mock.AnythingOfType("*time.Time"), // &user.CreatedAt
			mock.AnythingOfType("*time.Time"), // &user.UpdatedAt
		).Run(func(args mock.Arguments) {
			// Simular preenchimento do primeiro usuário
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "user1"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "user1@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Description 1"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()

		// Segundo usuário
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			// Simular preenchimento do segundo usuário
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(2)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "user2"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "user2@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Description 2"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = false
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
		}).Return(nil).Once()

		// Fim dos resultados
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		users, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		assert.Equal(t, int64(1), users[0].UID)
		assert.Equal(t, "user1", users[0].Username)
		assert.Equal(t, int64(2), users[1].UID)
		assert.Equal(t, "user2", users[1].Username)
		assert.Equal(t, false, users[1].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return nil when no users found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		users, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.Nil(t, users) // SUA FUNÇÃO RETORNA nil PARA LISTA VAZIA
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		dbErr := errors.New("database error")
		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(nil, dbErr)

		users, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		users, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		rowsErr := errors.New("rows error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		users, err := repo.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get single user", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "singleuser"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "single@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Single user description"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, mock.Anything).Return(mockRows, nil)

		users, err := repo.GetAll(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 1)
		assert.Equal(t, "singleuser", users[0].Username)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

func TestUserRepo_GetByEmail(t *testing.T) {
	t.Run("successfully get user by email", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "test@example.com"

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(1),           // UID (id na query)
				"testuser",         // username
				"test@example.com", // email
				"hashedpassword",   // password_hash
				true,               // status
				createdAt,          // created_at
				updatedAt,          // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@example.com", user.Email)
		assert.Equal(t, "hashedpassword", user.Password)
		assert.Equal(t, true, user.Status)
		assert.Equal(t, createdAt, user.CreatedAt)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when user does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "nonexistent@example.com"

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, user)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when database scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "test@example.com"

		scanErr := errors.New("scan error")
		mockRow := &mockDb.MockRow{Err: scanErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with false status", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "inactive@example.com"

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(2),               // UID
				"inactiveuser",         // username
				"inactive@example.com", // email
				"hashedpassword123",    // password_hash
				false,                  // status
				createdAt,              // created_at
				updatedAt,              // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(2), user.UID)
		assert.Equal(t, "inactiveuser", user.Username)
		assert.Equal(t, "inactive@example.com", user.Email)
		assert.Equal(t, "hashedpassword123", user.Password)
		assert.Equal(t, false, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully get user with empty password", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "empty@example.com"

		createdAt := time.Now()
		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				int64(3),            // UID
				"emptyuser",         // username
				"empty@example.com", // email
				"",                  // password_hash vazio
				true,                // status
				createdAt,           // created_at
				updatedAt,           // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(3), user.UID)
		assert.Equal(t, "emptyuser", user.Username)
		assert.Equal(t, "empty@example.com", user.Email)
		assert.Equal(t, "", user.Password)
		assert.Equal(t, true, user.Status)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error with empty email", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := ""

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, user)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error with invalid email format", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		email := "invalid-email"

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{email}).Return(mockRow)

		user, err := repo.GetByEmail(ctx, email)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Nil(t, user)
		mockDB.AssertExpectations(t)
	})
}

func TestUserRepo_GetByName(t *testing.T) {
	t.Run("successfully get users by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		mockRows := new(mockDb.MockRows)
		// Primeiro usuário
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),     // &user.UID
			mock.AnythingOfType("*string"),    // &user.Username
			mock.AnythingOfType("*string"),    // &user.Email
			mock.AnythingOfType("*string"),    // &user.Description
			mock.AnythingOfType("*bool"),      // &user.Status
			mock.AnythingOfType("*time.Time"), // &user.CreatedAt
			mock.AnythingOfType("*time.Time"), // &user.UpdatedAt
		).Run(func(args mock.Arguments) {
			// Simular preenchimento do primeiro usuário
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "testuser1"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "test1@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Description 1"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()

		// Segundo usuário
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			// Simular preenchimento do segundo usuário
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(2)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "testuser2"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "test2@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Description 2"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = false
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now().Add(-24 * time.Hour)
			}
		}).Return(nil).Once()

		// Fim dos resultados
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		assert.Equal(t, "testuser1", users[0].Username)
		assert.Equal(t, "testuser2", users[1].Username)
		assert.Equal(t, false, users[1].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no users found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "nonexistent"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.ErrorIs(t, err, errMsg.ErrNotFound) // ❌ Comportamento inconsistente
		assert.Nil(t, users)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		dbErr := errors.New("database error")
		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(nil, dbErr)

		users, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		scanErr := errors.New("scan error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrScan) // ✅ ErrScan, não ErrGet
		assert.Contains(t, err.Error(), scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "test"

		rowsErr := errors.New("rows error")
		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.ErrorIs(t, err, errMsg.ErrIterate) // ✅ ErrIterate, não ErrGet
		assert.Contains(t, err.Error(), rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("successfully get single user by name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := "single"

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			if ptr, ok := args.Get(0).(*int64); ok {
				*ptr = int64(1)
			}
			if ptr, ok := args.Get(1).(*string); ok {
				*ptr = "singleuser"
			}
			if ptr, ok := args.Get(2).(*string); ok {
				*ptr = "single@example.com"
			}
			if ptr, ok := args.Get(3).(*string); ok {
				*ptr = "Single user description"
			}
			if ptr, ok := args.Get(4).(*bool); ok {
				*ptr = true
			}
			if ptr, ok := args.Get(5).(*time.Time); ok {
				*ptr = time.Now()
			}
			if ptr, ok := args.Get(6).(*time.Time); ok {
				*ptr = time.Now()
			}
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%" + name + "%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 1)
		assert.Equal(t, "singleuser", users[0].Username)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrNotFound with empty name", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()
		name := ""

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		mockDB.On("Query", ctx, mock.Anything, []interface{}{"%%"}).Return(mockRows, nil)

		users, err := repo.GetByName(ctx, name)

		assert.ErrorIs(t, err, errMsg.ErrNotFound) // ❌ Comportamento inconsistente
		assert.Nil(t, users)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
