package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSupplierContactRelationTx(t *testing.T) {
	t.Run("successfully create new SupplierContactRelationTx instance", func(t *testing.T) {
		var db repo.DBExecutor

		result := NewSupplierContactRelationTx(db)

		assert.NotNil(t, result)

		_, ok := result.(*supplierContactRelationTx)
		assert.True(t, ok, "Expected result to be of type *supplierContactRelationTx")
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		var db repo.DBExecutor

		result := NewSupplierContactRelationTx(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db repo.DBExecutor

		instance1 := NewSupplierContactRelationTx(db)
		instance2 := NewSupplierContactRelationTx(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestSupplierContactRelationTx_CreateTx(t *testing.T) {
	t.Run("successfully create supplier contact relation within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierContactRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		mockResult := pgconn.NewCommandTag("INSERT 0 1")
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).Return(mockResult, nil)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (unique violation)", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierContactRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		// Retornar CommandTag vazio junto com o erro
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(pgconn.CommandTag{}, pgErr)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (string detection)", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierContactRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		// Testa o caso onde IsDuplicateKey detecta por string
		dbError := errors.New("duplicate key value violates unique constraint")
		// Retornar CommandTag vazio junto com o erro
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(pgconn.CommandTag{}, dbError)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierContactRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		// Retornar CommandTag vazio junto com o erro
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(pgconn.CommandTag{}, pgErr)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &supplierContactRelationTx{db: nil}
		ctx := context.Background()

		relation := &models.SupplierContactRelation{
			SupplierID: 1,
			ContactID:  2,
		}

		dbError := errors.New("database connection failed")
		// Retornar CommandTag vazio junto com o erro
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.SupplierID, relation.ContactID}).
			Return(pgconn.CommandTag{}, dbError)

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

}
