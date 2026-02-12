package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRepo_Create(t *testing.T) {
	t.Run("successfully create product category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			Name:        "Test Category",
			Description: "Test Description",
		}

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description}).Return(mockRow)

		result, err := repo.Create(ctx, category)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			Name:        "Test Category",
			Description: "Test Description",
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description}).Return(mockRow)

		result, err := repo.Create(ctx, category)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate on duplicate key", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			Name:        "Duplicated",
			Description: "Already exists",
		}

		pgErr := &pgconn.PgError{Code: "23505"}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description}).
			Return(mockRow)

		result, err := repo.Create(ctx, category)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})
}

func TestProductCategoryRepo_Update(t *testing.T) {
	t.Run("successfully update product category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			ID:          int64(1), // Correção: int64
			Name:        "Updated Category",
			Description: "Updated Description",
		}

		mockRow := &mockDb.MockRow{}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when category does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			ID:          int64(999), // Correção: int64
			Name:        "Non-existent Category",
			Description: "Non-existent Description",
		}

		mockRow := &mockDb.MockRow{Err: pgx.ErrNoRows}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			ID:          int64(1), // Correção: int64
			Name:        "Test Category",
			Description: "Test Description",
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()

		category := &models.ProductCategory{
			ID:          int64(1),
			Name:        "Duplicate Category",
			Description: "Duplicate Description",
		}

		pgErr := &pgconn.PgError{Code: "23505"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{category.Name, category.Description, category.ID}).Return(mockRow)

		err := repo.Update(ctx, category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro já cadastrado") // Apenas se seu Update trata duplicate
		mockDB.AssertExpectations(t)
	})
}

func TestProductCategoryRepo_Delete(t *testing.T) {
	t.Run("successfully delete product category", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		cmdTag := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{categoryID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, categoryID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when category does not exist", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(999)

		cmdTag := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{categoryID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, categoryID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRepo{db: mockDB}
		ctx := context.Background()
		categoryID := int64(1)

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{categoryID}).Return(pgconn.CommandTag{}, dbError)

		err := repo.Delete(ctx, categoryID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
