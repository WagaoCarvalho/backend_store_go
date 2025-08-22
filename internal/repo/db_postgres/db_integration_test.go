package repositories

import (
	"context"
	"errors"
	"os"
	"testing"

	config "github.com/WagaoCarvalho/backend_store_go/config"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db_postgres/db_mock"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestConnect_Success(t *testing.T) {
	mockPool := new(repo.MockPgxPool)
	mockConfig := &pgxpool.Config{}

	// Setar variável de ambiente obrigatória para o config.LoadDatabaseConfig()
	os.Setenv("DB_CONN_URL", "postgres://user:pass@localhost:5432/dbname?sslmode=disable")
	defer os.Unsetenv("DB_CONN_URL")

	mockPool.On("ParseConfig", mock.Anything).Return(mockConfig, nil)
	mockPool.On("NewWithConfig", mock.Anything, mockConfig).Return(&pgxpool.Pool{}, nil)

	pool, err := Connect(mockPool)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	mockPool.AssertExpectations(t)
}

func TestConnect_ParseConfigError(t *testing.T) {
	mockPool := new(repo.MockPgxPool)

	os.Setenv("DB_CONN_URL", "fake_conn_string")
	defer os.Unsetenv("DB_CONN_URL")

	mockPool.On("ParseConfig", mock.Anything).
		Return(nil, errors.New("erro parse config")).Once()

	pool, err := Connect(mockPool)
	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "erro parse config")

	mockPool.AssertExpectations(t)
}

func TestConnect_NewWithConfigError(t *testing.T) {
	mockPool := new(repo.MockPgxPool)

	os.Setenv("DB_CONN_URL", "fake_conn_string")
	defer os.Unsetenv("DB_CONN_URL")

	mockPool.On("ParseConfig", mock.Anything).
		Return(&pgxpool.Config{}, nil).Once()

	mockPool.On("NewWithConfig", mock.Anything, mock.Anything).
		Return(nil, errors.New("erro new pool")).Once()

	pool, err := Connect(mockPool)

	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "erro new pool")

	mockPool.AssertExpectations(t)
}

func TestConnect_DBConnURLNotDefined(t *testing.T) {
	originalLoader := config.LoadDatabaseConfig
	defer func() { config.LoadDatabaseConfig = originalLoader }()

	config.LoadDatabaseConfig = func() config.Database {
		return config.Database{
			ConnURL: "",
		}
	}

	pool, err := Connect(&RealPgxPool{})
	assert.Nil(t, pool)
	assert.ErrorIs(t, err, err_msg.ErrDBConnURLNotDefined)
}

func TestRealPgxPool_ParseConfig(t *testing.T) {
	realPool := &RealPgxPool{}

	// Usar uma connection string válida (pode ser fake para teste)
	connStr := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

	cfg, err := realPool.ParseConfig(connStr)

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestRealPgxPool_NewWithConfig(t *testing.T) {
	realPool := &RealPgxPool{}

	connStr := "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	cfg, err := realPool.ParseConfig(connStr)
	require.NoError(t, err)

	pool, err := realPool.NewWithConfig(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	// Limpar a conexão após o teste
	pool.Close()
}
