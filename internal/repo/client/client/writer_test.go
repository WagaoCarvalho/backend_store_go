package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClient_Create(t *testing.T) {
	t.Run("successfully create client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			Name:        "Test Client",
			Email:       utils.StrToPtr("test@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr("12.345.678/0001-90"),
			Description: "Test description",
			Status:      true,
			Version:     1,
		}

		now := time.Now()
		mockRow := &mockDb.MockRowWithID{
			IDValue:   1,
			TimeValue: now,
		}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.ID)
		assert.Equal(t, now, result.CreatedAt)
		assert.Equal(t, now, result.UpdatedAt)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "23505", Message: "duplicate key value violates unique constraint"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check constraint violation", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "23514", Message: "check constraint violation"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		dbErr := errors.New("database connection failed")
		mockRow := &mockDb.MockRow{Err: dbErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, dbErr.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrCreate when other pg error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		client := &models.Client{}

		pgErr := &pgconn.PgError{Code: "99999", Message: "generic pg error"}
		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.ErrorContains(t, err, pgErr.Message)
		mockDB.AssertExpectations(t)
	})

	t.Run("retorna ErrDBInvalidForeignKey quando ocorre violação de chave estrangeira", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{}

		// Erro de FK → código 23503
		pgErr := &pgconn.PgError{
			Code:    "23503",
			Message: "foreign key violation",
		}

		mockRow := &mockDb.MockRow{Err: pgErr}

		mockDB.
			On("QueryRow", ctx, mock.Anything, mock.AnythingOfType("[]interface {}")).
			Return(mockRow)

		result, err := repo.Create(ctx, client)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
		mockDB.AssertExpectations(t)
	})

}

func TestClientRepo_Update(t *testing.T) {
	t.Run("successfully update client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("updated@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "Updated description",
			Version:     1,
		}

		// Mock da PRIMEIRA chamada (SELECT version)
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1}, // currentVersion = 1 (igual ao client.Version)
		}

		// Mock da SEGUNDA chamada (UPDATE)
		updatedAt := time.Now()
		mockRowUpdate := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt, // updated_at
				2,         // version (incrementado)
			},
		}

		// Primeira chamada: SELECT version
		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		// Segunda chamada: UPDATE
		updateQuery := `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			client.Name,
			client.Email,
			client.CPF,
			client.CNPJ,
			client.Status,
			client.Description,
			client.ID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, client)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, client.UpdatedAt)
		assert.Equal(t, 2, client.Version)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when client not found", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:      999,
			Version: 1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - cliente não existe
		mockRowSelect := &mockDb.MockRow{
			Err: pgx.ErrNoRows,
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when version conflict occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("updated@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "Updated description",
			Version:     1, // Versão local
		}

		// Mock da PRIMEIRA chamada (SELECT version) - versão diferente no banco
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{2}, // currentVersion = 2 (diferente da local)
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrZeroVersion)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when SELECT query fails", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:      1,
			Version: 1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - erro no banco
		dbError := errors.New("database connection error")
		mockRowSelect := &mockDb.MockRow{
			Err: dbError,
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		err := repo.Update(ctx, client)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "erro ao consultar cliente")
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDuplicate when unique constraint violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("duplicate@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "Updated description",
			Version:     1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - sucesso
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1},
		}

		// Mock da SEGUNDA chamada (UPDATE) - erro de constraint única
		pgErr := &pgconn.PgError{
			Code: "23505", // unique_violation
		}
		mockRowUpdate := &mockDb.MockRow{
			Err: pgErr,
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			client.Name,
			client.Email,
			client.CPF,
			client.CNPJ,
			client.Status,
			client.Description,
			client.ID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrInvalidData when check constraint violation occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("invalid@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "Updated description",
			Version:     1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - sucesso
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1},
		}

		// Mock da SEGUNDA chamada (UPDATE) - erro de constraint de verificação
		pgErr := &pgconn.PgError{
			Code: "23514", // check_violation
		}
		mockRowUpdate := &mockDb.MockRow{
			Err: pgErr,
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			client.Name,
			client.Email,
			client.CPF,
			client.CNPJ,
			client.Status,
			client.Description,
			client.ID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, client)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockDB.AssertExpectations(t)
	})

	t.Run("return error when UPDATE query fails with generic error", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          1,
			Name:        "Updated Client",
			Email:       utils.StrToPtr("updated@example.com"),
			CPF:         utils.StrToPtr("123.456.789-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "Updated description",
			Version:     1,
		}

		// Mock da PRIMEIRA chamada (SELECT version) - sucesso
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{1},
		}

		// Mock da SEGUNDA chamada (UPDATE) - erro genérico no banco
		dbError := errors.New("generic database error")
		mockRowUpdate := &mockDb.MockRow{
			Err: dbError,
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			client.Name,
			client.Email,
			client.CPF,
			client.CNPJ,
			client.Status,
			client.Description,
			client.ID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, client)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "erro ao atualizar cliente")
		assert.Contains(t, err.Error(), dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("successfully update client with empty description", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()

		client := &models.Client{
			ID:          3,
			Name:        "Client No Description",
			Email:       utils.StrToPtr("nodesc@example.com"),
			CPF:         utils.StrToPtr("987.654.321-00"),
			CNPJ:        utils.StrToPtr(""),
			Status:      true,
			Description: "",
			Version:     2,
		}

		// Mock da PRIMEIRA chamada (SELECT version)
		mockRowSelect := &mockDb.MockRow{
			Values: []interface{}{2},
		}

		// Mock da SEGUNDA chamada (UPDATE)
		updatedAt := time.Now()
		mockRowUpdate := &mockDb.MockRow{
			Values: []interface{}{
				updatedAt,
				3,
			},
		}

		selectQuery := `
		SELECT version
		FROM clients
		WHERE id = $1
	`
		mockDB.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).Return(mockRowSelect)

		updateQuery := `
		UPDATE clients
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			cnpj = $4,
			status = $5,
			description = $6,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at, version
	`
		mockDB.On("QueryRow", ctx, updateQuery, []interface{}{
			client.Name,
			client.Email,
			client.CPF,
			client.CNPJ,
			client.Status,
			client.Description,
			client.ID,
		}).Return(mockRowUpdate)

		err := repo.Update(ctx, client)

		assert.NoError(t, err)
		assert.Equal(t, updatedAt, client.UpdatedAt)
		assert.Equal(t, 3, client.Version)
		assert.Empty(t, client.Description)
		mockDB.AssertExpectations(t)
	})
}

func TestClient_Delete(t *testing.T) {
	t.Run("successfully delete client", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		cmdTag := pgconn.NewCommandTag("DELETE 1")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{clientID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, clientID)

		assert.NoError(t, err)
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrDelete when database error occurs", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		dbError := errors.New("database connection failed")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{clientID}).Return(pgconn.CommandTag{}, dbError)

		err := repo.Delete(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.ErrorContains(t, err, dbError.Error())
		mockDB.AssertExpectations(t)
	})

	t.Run("return ErrNotFound when no rows are affected", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)
		repo := &clientRepo{db: mockDB}
		ctx := context.Background()
		clientID := int64(1)

		// DELETE 0 => nenhuma linha removida
		cmdTag := pgconn.NewCommandTag("DELETE 0")
		mockDB.On("Exec", ctx, mock.Anything, []interface{}{clientID}).Return(cmdTag, nil)

		err := repo.Delete(ctx, clientID)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockDB.AssertExpectations(t)
	})

}
