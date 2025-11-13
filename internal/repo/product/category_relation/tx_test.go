package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewProductCategoryRelationRepoTx(t *testing.T) {
	t.Run("successfully create new ProductCategoryRelationRepoTx instance", func(t *testing.T) {
		var db *pgxpool.Pool

		result := NewProductCategoryRelationRepoTx(db)

		assert.NotNil(t, result)

		_, ok := result.(*productCategoryRelationRepoTx)
		assert.True(t, ok, "Expected result to be of type *productCategoryRelationRepoTx")
	})

	t.Run("return instance with provided db pool", func(t *testing.T) {
		var db *pgxpool.Pool

		result := NewProductCategoryRelationRepoTx(db)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		var db *pgxpool.Pool

		instance1 := NewProductCategoryRelationRepoTx(db)
		instance2 := NewProductCategoryRelationRepoTx(db)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestProductCategoryRelationRepoTx_CreateTx(t *testing.T) {
	t.Run("successfully create product category relation within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &productCategoryRelationRepoTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		cmdTag := pgconn.NewCommandTag("INSERT 1")
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(cmdTag, nil)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &productCategoryRelationRepoTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(pgconn.CommandTag{}, pgErr)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &productCategoryRelationRepoTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(pgconn.CommandTag{}, pgErr)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &productCategoryRelationRepoTx{db: &pgxpool.Pool{}}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockTx.On("Exec", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(pgconn.CommandTag{}, dbError)

		result, err := repo.CreateTx(ctx, mockTx, relation)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), dbError.Error())
		mockTx.AssertExpectations(t)
	})
}
