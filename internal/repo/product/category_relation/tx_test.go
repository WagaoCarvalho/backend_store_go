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

func TestNewProductCategoryRelationRepoTx(t *testing.T) {
	t.Run("successfully create new ProductCategoryRelationRepoTx instance", func(t *testing.T) {
		result := NewProductCategoryRelationRepoTx()

		assert.NotNil(t, result)
		assert.IsType(t, &ProductCategoryRelationRepoTx{}, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		instance1 := NewProductCategoryRelationRepoTx()
		instance2 := NewProductCategoryRelationRepoTx()

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}

func TestProductCategoryRelationRepoTx_CreateTx(t *testing.T) {
	t.Run("successfully create product category relation within transaction", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &ProductCategoryRelationRepoTx{}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		mockRow := &mockDb.MockRow{
			Value: time.Now(),
		}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		err := repo.CreateTx(ctx, mockTx, relation)

		assert.NoError(t, err)
		assert.NotZero(t, relation.CreatedAt)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &ProductCategoryRelationRepoTx{}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		err := repo.CreateTx(ctx, mockTx, relation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &ProductCategoryRelationRepoTx{}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 999,
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		err := repo.CreateTx(ctx, mockTx, relation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockTx.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockTx := new(mockDb.MockTx)
		repo := &ProductCategoryRelationRepoTx{}
		ctx := context.Background()

		relation := &models.ProductCategoryRelation{
			ProductID:  1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbError}

		mockTx.On("QueryRow", ctx, mock.Anything, []interface{}{relation.ProductID, relation.CategoryID}).Return(mockRow)

		err := repo.CreateTx(ctx, mockTx, relation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbError.Error())
		mockTx.AssertExpectations(t)
	})
}
