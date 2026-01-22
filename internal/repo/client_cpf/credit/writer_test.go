package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/credit"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientCredit_Create(t *testing.T) {
	t.Run("successfully create client credit", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		credit := &models.ClientCredit{
			ClientID:      1,
			AllowCredit:   true,
			CreditLimit:   5000.0,
			CreditBalance: 1000.0,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   42,
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, credit)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(42), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		credit := &models.ClientCredit{
			ClientID: 1,
		}

		dbErr := errors.New("failed to insert client credit")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, credit)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestClientCredit_Update(t *testing.T) {
	t.Run("successfully update client credit", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		credit := &models.ClientCredit{
			ID:            1,
			AllowCredit:   true,
			CreditLimit:   10000.0,
			CreditBalance: 2500.0,
		}

		cmdTag := mockDb.MockCommandTag{RowsAffectedCount: 1}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{
			credit.AllowCredit,
			credit.CreditLimit,
			credit.CreditBalance,
			credit.ID,
		}).Return(cmdTag, nil)

		err := repo.Update(ctx, credit)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrUpdate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()

		credit := &models.ClientCredit{
			ID:            1,
			AllowCredit:   false,
			CreditLimit:   5000.0,
			CreditBalance: 1500.0,
		}

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{
			credit.AllowCredit,
			credit.CreditLimit,
			credit.CreditBalance,
			credit.ID,
		}).Return(nil, dbError)

		err := repo.Update(ctx, credit)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}

func TestClientCredit_Delete(t *testing.T) {
	t.Run("successfully delete client credit", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)

		cmdTag := mockDb.MockCommandTag{RowsAffectedCount: 1}
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, creditID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientCreditRepo{db: mockDB}
		ctx := context.Background()
		creditID := int64(1)
		dbError := errors.New("database connection failed")

		mockDB.On("Exec", ctx, mock.Anything, []interface{}{creditID}).Return(nil, dbError)

		err := repo.Delete(ctx, creditID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})
}
