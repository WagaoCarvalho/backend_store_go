package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

/* ===================== CREATE ===================== */

func TestClientRepo_Create_AllPaths(t *testing.T) {
	ctx := context.Background()

	baseClient := &models.ClientCpf{
		Name:        "Client",
		Email:       *utils.StrToPtr("client@test.com"),
		CPF:         *utils.StrToPtr("12345678900"),
		Description: "desc",
		Status:      true,
		Version:     1,
	}

	t.Run("success", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		now := time.Now()
		row := &mockDb.MockRowWithID{IDValue: 10, TimeValue: now}

		db.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(row)

		res, err := repo.Create(ctx, baseClient)

		assert.NoError(t, err)
		assert.Equal(t, int64(10), res.ID)
		assert.Equal(t, now, res.CreatedAt)
		assert.Equal(t, now, res.UpdatedAt)
	})

	t.Run("foreign key violation", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		row := &mockDb.MockRow{
			Err: &pgconn.PgError{Code: "23503"},
		}

		db.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(row)

		res, err := repo.Create(ctx, baseClient)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, errMsg.ErrDBInvalidForeignKey)
	})

	t.Run("unique violation", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		row := &mockDb.MockRow{
			Err: &pgconn.PgError{Code: "23505", ConstraintName: "clients_cpf_cpf_key"},
		}

		db.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(row)

		res, err := repo.Create(ctx, baseClient)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
	})

	t.Run("check violation", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		row := &mockDb.MockRow{
			Err: &pgconn.PgError{Code: "23514"},
		}

		db.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(row)

		res, err := repo.Create(ctx, baseClient)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("generic error", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		dbErr := errors.New("db down")
		row := &mockDb.MockRow{Err: dbErr}

		db.On("QueryRow", ctx, mock.Anything, mock.Anything).
			Return(row)

		res, err := repo.Create(ctx, baseClient)

		assert.Nil(t, res)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), dbErr.Error())
	})
}

/* ===================== UPDATE ===================== */

func TestClientRepo_Update_AllPaths(t *testing.T) {
	ctx := context.Background()

	// Helper para criar um novo cliente em cada teste
	createClient := func() *models.ClientCpf {
		return &models.ClientCpf{
			ID:          1,
			Name:        "Updated",
			Email:       *utils.StrToPtr("u@test.com"),
			CPF:         *utils.StrToPtr("12345678900"),
			Status:      true,
			Description: "desc",
			Version:     1,
		}
	}

	selectQuery := `
		SELECT version
		FROM clients_cpf
		WHERE id = $1
	`

	updateQuery := `
		UPDATE clients_cpf
		SET 
			name = $1,
			email = $2,
			cpf = $3,
			status = $4,
			description = $5,
			version = version + 1,
			updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at, version
	`

	t.Run("not found", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Err: pgx.ErrNoRows})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
	})

	t.Run("select error", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		dbErr := errors.New("select failed")
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Err: dbErr})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "select failed")
	})

	t.Run("version conflict", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		// Retornar versão 2 (diferente da versão do cliente que é 1)
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{2}})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
	})

	t.Run("successful update", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		now := time.Now()
		newVersion := 2

		// Retornar versão 1 (igual à versão do cliente)
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{1}})

		db.On("QueryRow", ctx, updateQuery, mock.Anything).
			Return(&mockDb.MockRow{
				Values: []interface{}{now, newVersion},
			})

		err := repo.Update(ctx, client)
		assert.NoError(t, err)
		assert.Equal(t, now, client.UpdatedAt)
		assert.Equal(t, newVersion, client.Version)
	})

	t.Run("update unique violation", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		// Retornar versão 1 (igual) para passar pela verificação de versão
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{1}})

		// Mock da query UPDATE com erro de violação única
		db.On("QueryRow", ctx, updateQuery, mock.Anything).
			Return(&mockDb.MockRow{Err: &pgconn.PgError{Code: "23505", Message: "duplicate key"}})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrDuplicate)
	})

	t.Run("update check violation", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		// Retornar versão 1 (igual)
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{1}})

		// Mock da query UPDATE com erro de violação de check
		db.On("QueryRow", ctx, updateQuery, mock.Anything).
			Return(&mockDb.MockRow{Err: &pgconn.PgError{Code: "23514", Message: "check constraint"}})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("update generic error", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		// Retornar versão 1 (igual)
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{1}})

		dbErr := errors.New("update failed")
		db.On("QueryRow", ctx, updateQuery, mock.Anything).
			Return(&mockDb.MockRow{Err: dbErr})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "update failed")
	})

	t.Run("update with other pg error", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		client := createClient()
		// Retornar versão 1 (igual)
		db.On("QueryRow", ctx, selectQuery, []interface{}{client.ID}).
			Return(&mockDb.MockRow{Values: []interface{}{1}})

		// Testar um código de erro PostgreSQL diferente
		db.On("QueryRow", ctx, updateQuery, mock.Anything).
			Return(&mockDb.MockRow{Err: &pgconn.PgError{Code: "22000", Message: "data exception"}})

		err := repo.Update(ctx, client)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), "data exception")
	})
}

/* ===================== DELETE ===================== */

func TestClientRepo_Delete_AllPaths(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		db.On("Exec", ctx, mock.Anything, []interface{}{int64(1)}).
			Return(pgconn.NewCommandTag("DELETE 1"), nil)

		err := repo.Delete(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		db.On("Exec", ctx, mock.Anything, []interface{}{int64(1)}).
			Return(pgconn.NewCommandTag("DELETE 0"), nil)

		err := repo.Delete(ctx, 1)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
	})

	t.Run("db error", func(t *testing.T) {
		db := new(mockDb.MockDatabase)
		repo := &clientCpfRepo{db: db}

		dbErr := errors.New("delete failed")
		db.On("Exec", ctx, mock.Anything, []interface{}{int64(1)}).
			Return(pgconn.CommandTag{}, dbErr)

		err := repo.Delete(ctx, 1)

		assert.ErrorIs(t, err, errMsg.ErrDelete)
		assert.Contains(t, err.Error(), dbErr.Error())
	})
}
