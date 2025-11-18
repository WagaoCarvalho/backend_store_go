package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierCategoryRelationRepo_Create(t *testing.T) {
	t.Run("successfully create supplier category relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
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

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		assert.Equal(t, createdAt, relation.CreatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.NotErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.Contains(t, err.Error(), dbError.Error())
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key with string message", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		// Testa o caso onde IsDuplicateKey detecta por string
		dbError := errors.New("duplicate key value violates unique constraint")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockDB.On("QueryRow", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.CategoryID}).Return(mockRow)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationRepo_Delete(t *testing.T) {
	ctx := context.Background()
	supplierID := int64(1)
	categoryID := int64(2)

	t.Run("successfully delete relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, supplierID, categoryID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - no rows affected (not found)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.Delete(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, categoryID},
		).Return(pgconn.CommandTag{}, pgErr).Once()

		err := repo.Delete(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrDelete)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, categoryID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.Delete(ctx, supplierID, categoryID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationRepo_DeleteAllBySupplierID(t *testing.T) {
	ctx := context.Background()
	supplierID := int64(1)

	t.Run("successfully delete all relations by supplier id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID},
		).Return(pgconn.NewCommandTag("DELETE 3"), nil).Once()

		err := repo.DeleteAllBySupplierID(ctx, supplierID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID},
		).Return(pgconn.CommandTag{}, pgErr).Once()

		err := repo.DeleteAllBySupplierID(ctx, supplierID)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrDelete)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.DeleteAllBySupplierID(ctx, supplierID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})
}
