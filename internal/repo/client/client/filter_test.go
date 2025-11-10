package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetAll(t *testing.T) {
	t.Run("successfully get all clients", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &client{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 1
			*args[1].(*string) = "João Silva"
			email := "joao@email.com"
			cpf := "123.456.789-09"
			cnpj := ""
			*args[2].(**string) = &email
			*args[3].(**string) = &cpf
			*args[4].(**string) = &cnpj
			*args[5].(*string) = "Cliente de teste"
			*args[6].(*bool) = true
			*args[7].(*time.Time) = now
			*args[8].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &model.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "João Silva", result[0].Name)
		assert.Equal(t, "Cliente de teste", result[0].Description)
		assert.Equal(t, true, result[0].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &client{db: mockDB}
		ctx := context.Background()
		dbErr := errors.New("database failure")

		filter := &model.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		emptyRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(emptyRows, dbErr)

		result, err := repo.GetAll(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &client{db: mockDB}
		ctx := context.Background()
		scanErr := errors.New("failed to scan row")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		filter := &model.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &client{db: mockDB}
		ctx := context.Background()
		rowsErr := errors.New("iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		filter := &model.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.GetAll(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters name, email and status correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &client{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		mockRows := new(mockDb.MockRows)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan",
			mock.AnythingOfType("*int64"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("**string"),
			mock.AnythingOfType("*string"),
			mock.AnythingOfType("*bool"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 42
			*args[1].(*string) = "Maria Souza"
			email := "maria@email.com"
			cpf := "987.654.321-00"
			cnpj := ""
			*args[2].(**string) = &email
			*args[3].(**string) = &cpf
			*args[4].(**string) = &cnpj
			*args[5].(*string) = "Cliente filtro test"
			*args[6].(*bool) = true
			*args[7].(*time.Time) = now
			*args[8].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		status := true
		filter := &model.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 5,
			},
			Name:   "Maria",
			Email:  "maria@",
			Status: &status,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {

			if len(args) != 5 {
				return false
			}

			name := args[0].(string)
			email := args[1].(string)
			stat := args[2].(bool)
			limit := args[3].(int)
			offset := args[4].(int)

			return name == "%Maria%" &&
				email == "%maria@%" &&
				stat == true &&
				limit == 10 &&
				offset == 5
		})).Return(mockRows, nil)

		result, err := repo.GetAll(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Maria Souza", result[0].Name)
		assert.Equal(t, true, result[0].Status)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}
