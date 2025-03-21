package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func TestConnect_Success(t *testing.T) {
	mockPgx := new(MockPgxPool)

	// Configura o mock para simular uma conexão bem-sucedida
	config := &pgxpool.Config{}
	mockPgx.On("ParseConfig", "postgres://user:pass@localhost:5432/db_postgres?sslmode=disable").Return(config, nil)
	mockPgx.On("NewWithConfig", context.Background(), config).Return(&pgxpool.Pool{}, nil)

	// Chama a função Connect com o mock
	pool, err := Connect(mockPgx)

	// Verifica se o pool foi retornado corretamente
	assert.NoError(t, err)
	assert.NotNil(t, pool)
	mockPgx.AssertExpectations(t)
}

func TestConnect_ParseConfigError(t *testing.T) {
	mockPgx := new(MockPgxPool)

	// Configura o mock para simular um erro ao parsear a URL
	mockPgx.On("ParseConfig", "postgres://user:pass@localhost:5432/db_postgres?sslmode=disable").Return(&pgxpool.Config{}, errors.New("erro ao parsear URL"))

	// Chama a função Connect com o mock
	pool, err := Connect(mockPgx)

	// Verifica se o erro foi retornado corretamente
	assert.Error(t, err)
	assert.Nil(t, pool)
	mockPgx.AssertExpectations(t)
}

func TestConnect_NewWithConfigError(t *testing.T) {
	mockPgx := new(MockPgxPool)

	// Configura o mock para simular um erro ao criar a conexão
	config := &pgxpool.Config{}
	mockPgx.On("ParseConfig", "postgres://user:pass@localhost:5432/db_postgres?sslmode=disable").Return(config, nil)
	mockPgx.On("NewWithConfig", context.Background(), config).Return(&pgxpool.Pool{}, errors.New("erro ao conectar"))

	// Chama a função Connect com o mock
	pool, err := Connect(mockPgx)

	// Verifica se o erro foi retornado corretamente
	assert.Error(t, err)
	assert.Nil(t, pool)
	mockPgx.AssertExpectations(t)
}
