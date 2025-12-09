package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client/filter"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_GetAll(t *testing.T) {

	t.Run("successfully get all clients", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
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
			mock.AnythingOfType("*int"),
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
			*args[7].(*int) = 1
			*args[8].(*time.Time) = now
			*args[9].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "João Silva", result[0].Name)
		assert.Equal(t, "Cliente de teste", result[0].Description)
		assert.Equal(t, true, result[0].Status)
		assert.Equal(t, 1, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrGet when query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
		ctx := context.Background()
		dbErr := errors.New("database failure")

		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		emptyRows := new(mockDb.MockRows)
		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(emptyRows, dbErr)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrScan when scan fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
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
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Return(scanErr).Once()
		mockRows.On("Close").Return()

		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrScan)
		assert.ErrorContains(t, err, scanErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("return ErrIterate when rows iteration fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
		ctx := context.Background()
		rowsErr := errors.New("iteration error")

		mockRows := new(mockDb.MockRows)
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(rowsErr)
		mockRows.On("Close").Return()

		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
		}

		mockDB.On("Query", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIterate)
		assert.ErrorContains(t, err, rowsErr.Error())
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters name, email and status correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
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
			mock.AnythingOfType("*int"),
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
			*args[7].(*int) = 2
			*args[8].(*time.Time) = now
			*args[9].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		status := true
		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 5,
			},
			Name:   "Maria",
			Email:  "maria@",
			Status: &status,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {

			if len(args) != 3 {
				// nome, email, status — limit/offset são concatenados direto na query
				return false
			}

			name := args[0].(string)
			email := args[1].(string)
			stat := args[2].(bool)

			return name == "Maria" &&
				email == "maria@" &&
				stat == true
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Maria Souza", result[0].Name)
		assert.Equal(t, true, result[0].Status)
		assert.Equal(t, 2, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters CPF and CNPJ correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
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
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 7
			*args[1].(*string) = "Carlos Lima"
			email := "carlos@email.com"
			cpf := "111.222.333-44"
			cnpj := "12.345.678/0001-99"
			*args[2].(**string) = &email
			*args[3].(**string) = &cpf
			*args[4].(**string) = &cnpj
			*args[5].(*string) = "Cliente com CPF e CNPJ"
			*args[6].(*bool) = false
			*args[7].(*int) = 3
			*args[8].(*time.Time) = now
			*args[9].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			CPF:  "111.222.333-44",
			CNPJ: "12.345.678/0001-99",
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			if len(args) != 2 {
				// cpf e cnpj
				return false
			}
			cpf := args[0].(string)
			cnpj := args[1].(string)
			return cpf == "111.222.333-44" && cnpj == "12.345.678/0001-99"
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Carlos Lima", result[0].Name)
		assert.Equal(t, false, result[0].Status)
		assert.Equal(t, 3, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("apply filters version and date ranges correctly", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientFilterRepo{db: mockDB}
		ctx := context.Background()

		now := time.Now()
		createdFrom := now.Add(-48 * time.Hour)
		createdTo := now.Add(-24 * time.Hour)
		updatedFrom := now.Add(-12 * time.Hour)
		updatedTo := now

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
			mock.AnythingOfType("*int"),
			mock.AnythingOfType("*time.Time"),
			mock.AnythingOfType("*time.Time"),
		).Run(func(args mock.Arguments) {
			*args[0].(*int64) = 99
			*args[1].(*string) = "Teste Data e Versão"
			email := "teste@teste.com"
			cpf := ""
			cnpj := ""
			*args[2].(**string) = &email
			*args[3].(**string) = &cpf
			*args[4].(**string) = &cnpj
			*args[5].(*string) = "Filtro de data e versão"
			*args[6].(*bool) = true
			*args[7].(*int) = 9
			*args[8].(*time.Time) = now
			*args[9].(*time.Time) = now
		}).Return(nil).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return()

		version := 9
		filter := &filterClient.ClientFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
			Version:     &version,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
		}

		mockDB.On("Query", ctx, mock.Anything, mock.MatchedBy(func(args []interface{}) bool {
			if len(args) != 5 {
				// version, created_from, created_to, updated_from, updated_to
				return false
			}

			v := args[0].(int)
			cf := args[1].(time.Time)
			ct := args[2].(time.Time)
			uf := args[3].(time.Time)
			ut := args[4].(time.Time)

			return v == 9 &&
				cf.Equal(createdFrom) &&
				ct.Equal(createdTo) &&
				uf.Equal(updatedFrom) &&
				ut.Equal(updatedTo)
		})).Return(mockRows, nil)

		result, err := repo.Filter(ctx, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Teste Data e Versão", result[0].Name)
		assert.Equal(t, 9, result[0].Version)
		mockDB.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
