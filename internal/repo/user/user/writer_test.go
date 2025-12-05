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
				int64(1),
				2,
				"Test user description",
				createdAt,
				updatedAt,
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, user, result)
		assert.Equal(t, int64(1), user.UID)
		assert.Equal(t, 2, user.Version)
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
				int64(2),
				1,
				"Inactive user description",
				createdAt,
				updatedAt,
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
				int64(3),
				2,
				"",
				createdAt,
				updatedAt,
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
				int64(4),
				1,
				"Version test user",
				createdAt,
				updatedAt,
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(4), user.UID)
		assert.Equal(t, 1, user.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData for duplicate key (username)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "existinguser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Simular erro de chave duplicada do PostgreSQL
		dbError := &pgconn.PgError{
			Code:    "23505", // unique_violation
			Message: "duplicate key value violates unique constraint \"users_username_key\"",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "dados já existem no sistema")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData for duplicate key (email)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "newuser",
			Email:       "existing@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Simular erro de chave duplicada para email
		dbError := &pgconn.PgError{
			Code:    "23505", // unique_violation
			Message: "duplicate key value violates unique constraint \"users_email_key\"",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "dados já existem no sistema")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData for check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "invalid-email",
			Password:    "short",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Simular erro de constraint check
		dbError := &pgconn.PgError{
			Code:    "23514", // check_violation
			Message: "new row for relation \"users\" violates check constraint \"users_email_check\"",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "dados inválidos")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate for other database errors", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Erro genérico que não é duplicate key nem check violation
		dbError := &pgconn.PgError{
			Code:    "08000", // connection_exception
			Message: "connection to server lost",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.NotContains(t, err.Error(), "dados já existem no sistema")
		assert.NotContains(t, err.Error(), "dados inválidos")
		assert.Contains(t, err.Error(), "connection to server lost")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate for non-pgconn errors", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Erro que não é pgconn.PgError
		dbError := errors.New("network error")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), "network error")
		mockDB.AssertExpectations(t)
	})

	t.Run("handle multiple duplicate key constraints", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "existinguser",
			Email:       "existing@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Erro com constraint name diferente
		dbError := &pgconn.PgError{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint \"users_some_other_constraint\"",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "dados já existem no sistema")
		mockDB.AssertExpectations(t)
	})

	t.Run("handle different check constraint violations", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword123",
			Description: "Test user",
			Status:      true,
			Version:     1,
		}

		// Outro tipo de constraint check
		dbError := &pgconn.PgError{
			Code:    "23514",
			Message: "new row for relation \"users\" violates check constraint \"users_status_check\"",
		}
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, user)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		assert.Contains(t, err.Error(), "dados inválidos")
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

		// Mock da PRIMEIRA chamada (SELECT version)
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1}, // currentVersion = 1 (igual ao user.Version)
		}

		// Mock da SEGUNDA chamada (UPDATE)
		updatedAt := time.Now()
		mockRowUpdate := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
				2,         // version (incrementado)
			},
		}

		// Primeira chamada: SELECT version
		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		// Segunda chamada: UPDATE
		updateQuery := `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 2, user.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when user not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:     999,
			Version: 1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - usuário não existe
		mockRowSelect := &mockDb.MockRow{
			Err: pgx.ErrNoRows,
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		err := repo.Update(ctx, user)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when version conflict occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:         1,
			Username:    "updateduser",
			Email:       "updated@example.com",
			Description: "Updated description",
			Status:      true,
			Version:     1, // Versão local
		}

		// Mock da PRIMEIRA chamada (SELECT version) - versão diferente no banco
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{2}, // currentVersion = 2 (diferente da local)
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		err := repo.Update(ctx, user)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when SELECT query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userRepo{db: mockDB}
		ctx := context.Background()

		user := &models.User{
			UID:     1,
			Version: 1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - erro no banco
		dbError := errors.New("database connection error")
		mockRowSelect := &mockDb.MockRow{
			Err: dbError,
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		err := repo.Update(ctx, user)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "erro ao consultar usuário")
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when UPDATE query fails", func(t *testing.T) {
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

		// Mock da PRIMEIRA chamada (SELECT version) - sucesso
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1},
		}

		// Mock da SEGUNDA chamada (UPDATE) - erro no banco
		dbError := errors.New("update constraint violation")
		mockRowUpdate := &mockDb.MockRow{
			Err: dbError,
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, user)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "erro ao atualizar usuário")
		assert.Contains(t, err.Error(), dbError.Error())
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

		// Mock da PRIMEIRA chamada (SELECT version)
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{3},
		}

		// Mock da SEGUNDA chamada (UPDATE)
		updatedAt := time.Now()
		mockRowUpdate := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt,
				4,
			},
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 4, user.Version)
		assert.False(t, user.Status)
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

		// Mock da PRIMEIRA chamada (SELECT version)
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{2},
		}

		// Mock da SEGUNDA chamada (UPDATE)
		updatedAt := time.Now()
		mockRowUpdate := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt,
				3,
			},
		}

		selectQuery := `
		SELECT version
		FROM users
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{user.UID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			user.Username,
			user.Email,
			user.Description,
			user.Status,
			user.UID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, user)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
		assert.Equal(t, 3, user.Version)
		assert.Empty(t, user.Description)
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

		mockResult := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{userID}).Return(mockResult, nil)

		err := repo.Delete(ctx, userID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
