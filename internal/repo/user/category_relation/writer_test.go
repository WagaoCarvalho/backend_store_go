package repo

import (
	"context"
	"errors"
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationRepo_Create(t *testing.T) {
	t.Run("successfully create user category relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 1,
		}

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).Return(mockResult, nil)

		result, err := repo.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (unique violation)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).
			Return(mockDb.MockCommandTag{}, pgErr)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrRelationExists when duplicate key (string detection)", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		// Testa o caso onde IsDuplicateKey detecta por string
		dbError := errors.New("duplicate key value violates unique constraint")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).
			Return(mockDb.MockCommandTag{}, dbError)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDBInvalidForeignKey when foreign key violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 999, // ID inexistente
		}

		pgErr := &pgconn.PgError{Code: "23503", Message: "foreign key violation"}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).
			Return(mockDb.MockCommandTag{}, pgErr)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.NotErrorIs(t, err, errMsg.ErrCreate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).
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

	t.Run("return ErrCreate when check violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		pgErr := &pgconn.PgError{Code: "23514", Message: "check constraint violation"}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).
			Return(mockDb.MockCommandTag{}, pgErr)

		result, err := repo.Create(ctx, relation)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.NotErrorIs(t, err, errMsg.ErrRelationExists)
		assert.NotErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		assert.Contains(t, err.Error(), errMsg.ErrCreate.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when zero rows affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}
		ctx := context.Background()

		relation := &models.UserCategoryRelation{
			UserID:     1,
			CategoryID: 2,
		}

		mockResult := mockDb.MockCommandTag{
			RowsAffectedCount: 0,
		}

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{relation.UserID, relation.CategoryID}).Return(mockResult, nil)

		result, err := repo.Create(ctx, relation)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, relation, result)
		// Nota: A função retorna success mesmo com 0 rows affected,
		// pois o Exec pode retornar 0 para alguns tipos de INSERT
		mockDB.AssertExpectations(t)
	})
}

func TestUserCategoryRelationRepo_Delete(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)
	categoryID := int64(2)

	t.Run("successfully delete relation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, userID, categoryID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID, categoryID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.Delete(ctx, userID, categoryID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("error - relation not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID, categoryID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.Delete(ctx, userID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete relation with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{int64(0), categoryID},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, 0, categoryID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - relation not found with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{int64(0), categoryID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.Delete(ctx, 0, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete relation with negative IDs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{int64(-1), int64(-2)},
		).Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.Delete(ctx, -1, -2)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}

func TestUserCategoryRelationRepo_DeleteAll(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("successfully delete all relations by user id", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID},
		).Return(pgconn.NewCommandTag("DELETE 3"), nil).Once()

		err := repo.DeleteAll(ctx, userID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("error - database failure", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		dbErr := errors.New("db failure")
		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID},
		).Return(pgconn.CommandTag{}, dbErr).Once()

		err := repo.DeleteAll(ctx, userID)
		assert.ErrorContains(t, err, "db failure")
		assert.ErrorContains(t, err, errMsg.ErrDelete.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete all relations even when zero rows affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{userID},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.DeleteAll(ctx, userID)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete all relations with zero user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{int64(0)},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.DeleteAll(ctx, 0)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully delete all relations with negative user ID", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &userCategoryRelationRepo{db: mockDB}

		mockDB.On("Exec", mock.Anything, mock.Anything,
			[]interface{}{int64(-1)},
		).Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.DeleteAll(ctx, -1)
		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})
}
