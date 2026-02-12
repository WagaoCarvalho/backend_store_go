package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationRepo_Create(t *testing.T) {
	t.Run("successfully create product category relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()
		now := time.Now()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		mockRow := &mockDb.MockRow{
			Value: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, []any{relation.ProductID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		assert.Equal(t, now, relation.CreatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.Contains(t, err.Error(), "relação já existe")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 999,
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.Contains(t, err.Error(), "chave estrangeira inválida")
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}
		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestProductCategoryRelationRepo_Delete(t *testing.T) {
	ctx := context.Background()
	productID := int64(1)
	categoryID := int64(2)

	t.Run("successfully delete relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, productID, categoryID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - no rows affected (not found)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.Delete(ctx, productID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Contains(t, err.Error(), "relação não encontrada")
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID, categoryID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.Delete(ctx, productID, categoryID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		mockDB.AssertExpectations(t)
	})
}

func TestProductCategoryRelationRepo_DeleteAll(t *testing.T) {
	ctx := context.Background()
	productID := int64(1)

	t.Run("successfully delete all relations with rows", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID},
		).Return(pgconn.NewCommandTag("DELETE 3"), nil).Once()

		err := repo.DeleteAll(ctx, productID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete all relations with zero rows (no error)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.DeleteAll(ctx, productID)
		assert.NoError(t, err) // Não retorna ErrNotFound
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &productCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{productID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.DeleteAll(ctx, productID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		mockDB.AssertExpectations(t)
	})
}
