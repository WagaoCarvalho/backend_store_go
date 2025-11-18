package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSupplierCategoryRelationTx(t *testing.T) {
	t.Run("successfully create new SupplierCategoryRelationTx instance", func(t *testing.T) {
		var db repo.DBExecutor

		result := NewSupplierCategoryRelationTx(db)

		assert.NotNil(t, result)

		_, ok := result.(*supplierCategoryRelationTx)
		assert.True(t, ok, "Expected result to be of type *supplierCategoryRelationTx")
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		var db repo.DBExecutor

		result := NewSupplierCategoryRelationTx(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db repo.DBExecutor

		instance1 := NewSupplierCategoryRelationTx(db)
		instance2 := NewSupplierCategoryRelationTx(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestSupplierCategoryRelationTx_CreateTx(t *testing.T) {
	t.Run("successfully create supplier category relation within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierCategoryRelationTx{db: nil} // db não é usado no CreateTx
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		createdAt := time.Now()
		mockRow := &mockDb.MockRow{
			Values: []interface{}{
				createdAt, // created_at
			},
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		assert.Equal(t, createdAt, relation.CreatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierCategoryRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierCategoryRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierCategoryRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.NotErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key with string message", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierCategoryRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		// Testa o caso onde IsDuplicateKey detecta por string
		dbError := errors.New("duplicate key value violates unique constraint")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockTx.AssertExpectations(t)
	})
}
