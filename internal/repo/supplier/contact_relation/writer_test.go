package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierContactRelationRepo_Create(t *testing.T) {
	t.Run("successfully create supplier contact relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1,
		}

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).Return(mockResult, nil)

		result, err := repo.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (unique violation)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(mockDb.MockCommandTag{}, pgErr)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (string detection)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		// Testa o caso onde IsDuplicateKey detecta por string
		dbError := errors.New("duplicate key value violates unique constraint")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(mockDb.MockCommandTag{}, dbError)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(mockDb.MockCommandTag{}, pgErr)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(mockDb.MockCommandTag{}, dbError)

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
}

func TestSupplierContactRelationRepo_Delete(t *testing.T) {
	ctx := context.Background()
	supplierID := int64(1)
	contactID := int64(2)

	t.Run("successfully delete relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, contactID},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, supplierID, contactID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID, contactID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.Delete(ctx, supplierID, contactID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestSupplierContactRelationRepo_DeleteAll(t *testing.T) {
	ctx := context.Background()
	supplierID := int64(1)

	t.Run("successfully delete all relations by supplier id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID},
		).Return(pgconn.NewCommandTag("DELETE 3"), nil).Once()

		err := repo.DeleteAll(ctx, supplierID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &supplierContactRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{supplierID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.DeleteAll(ctx, supplierID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})
}
