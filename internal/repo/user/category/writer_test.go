package repo

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRepo_Create(t *testing.T) {
	t.Run("successfully create user category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			Name:        "Test User Category",
			Description: "Test User Description",
		}

		expectedID := uint(1)
		createdAt := time.Now()
		updatedAt := time.Now()

		// Criando um mock row customizado que lida corretamente com uint
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				uint64(expectedID), // Usar uint64 em vez de int64
				createdAt,
				updatedAt,
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description}).Return(mockRow)

		result, err := repo.Create(ctx, category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedID, result.ID)
		assert.Equal(t, createdAt, result.CreatedAt)
		assert.Equal(t, updatedAt, result.UpdatedAt)
		assert.Equal(t, "Test User Category", result.Name)
		assert.Equal(t, "Test User Description", result.Description)
		mockDB.AssertExpectations(t)
	})

	t.Run("error creating user category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			Name:        "Test User Category",
			Description: "Test User Description",
		}

		mockErr := fmt.Errorf("scan failed")

		// Mockando a row que falha no Scan
		mockRow := &mockDb.MockRow{
			Err: mockErr,
		}

		// QueryRow deve devolver a mockRow que retorna erro no Scan
		mockDB.
			On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description}).
			Return(mockRow)

		result, err := repo.Create(ctx, category)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), "scan failed")

		mockDB.AssertExpectations(t)
	})

}

func TestUserCategoryRepo_Update(t *testing.T) {
	t.Run("successfully update user category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			ID:          1,
			Name:        "Updated User Category",
			Description: "Updated Description",
		}

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, category.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when category does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			ID:          999,
			Name:        "Non-existent Category",
			Description: "None",
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.NotErrorIs(t, err, errMsg.ErrUpdate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			ID:          1,
			Name:        "Test Category",
			Description: "Test Description",
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.NotErrorIs(t, err, errMsg.ErrNotFound)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when constraint violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			ID:          1,
			Name:        "Duplicate Category",
			Description: "Test Description",
		}

		pgError := &pgconn.PgError{
			Code:    "23505",
			Message: "duplicate key value violates unique constraint",
		}
		mockRow := &mockDb.MockRow{Err: pgError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.NotErrorIs(t, err, errMsg.ErrNotFound)
		assert.Contains(t, err.Error(), pgError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())

		mockDB.AssertExpectations(t)
	})

	t.Run("successfully update user category with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &model.UserCategory{
			ID:          1,
			Name:        "Updated User Category",
			Description: "",
		}

		updatedAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt,
			},
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, category.UpdatedAt)
		assert.Equal(t, "", category.Description)
		mockDB.AssertExpectations(t)
	})
}

func TestUserCategoryRepo_Delete(t *testing.T) {
	t.Run("successfully delete user category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		cmdTag := pgconn.NewCommandTag("DELETE 1")

		mockDB.
			On("Exec", ctx, mock.Anything, []interface{}{categoryID}).
			Return(cmdTag, nil)

		err := repo.Delete(ctx, categoryID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when category does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(999)

		cmdTag := pgconn.NewCommandTag("DELETE 0")

		mockDB.
			On("Exec", ctx, mock.Anything, []interface{}{categoryID}).
			Return(cmdTag, nil)

		err := repo.Delete(ctx, categoryID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		dbError := errors.New("database connection failed")

		mockDB.
			On("Exec", ctx, mock.Anything, []interface{}{categoryID}).
			Return(pgconn.CommandTag{}, dbError)

		err := repo.Delete(ctx, categoryID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
