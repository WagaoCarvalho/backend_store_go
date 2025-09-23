package repo

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/config"
	mockRepo "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func init() {
	// Carrega .env local apenas se existir
	_ = godotenv.Load("../../../.env.test")
}

func TestConnect_Success(t *testing.T) {
	mockPool := new(mockRepo.MockPgxPool)
	mockConfig := &pgxpool.Config{}

	connURL := os.Getenv("DB_CONN_URL")
	if connURL == "" {
		t.Skip("DB_CONN_URL não definido no .env ou no CI")
	}

	mockPool.On("ParseConfig", connURL).Return(mockConfig, nil)
	mockPool.On("NewWithConfig", mock.Anything, mockConfig).Return(&pgxpool.Pool{}, nil)

	pool, err := Connect(mockPool)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	mockPool.AssertExpectations(t)
}

func TestConnect_ParseConfigError(t *testing.T) {
	mockPool := new(mockRepo.MockPgxPool)

	mockPool.On("ParseConfig", mock.Anything).
		Return(nil, errors.New("erro parse config")).Once()

	pool, err := Connect(mockPool)
	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Contains(t, err.Error(), "erro parse config")

	mockPool.AssertExpectations(t)
}

func TestConnect_NewWithConfigError(t *testing.T) {
	mockPool := new(mockRepo.MockPgxPool)

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
		return config.Database{ConnURL: ""}
	}

	pool, err := Connect(&RealPgxPool{})
	assert.Nil(t, pool)
	assert.ErrorIs(t, err, errMsg.ErrDBConnURLNotDefined)
}

func TestRealPgxPool_ParseConfig(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipando teste de integração no CI")
	}
	realPool := &RealPgxPool{}

	connStr := os.Getenv("DB_CONN_URL")
	if connStr == "" {
		t.Skip("DB_CONN_URL não definido no .env ou no CI")
	}

	cfg, err := realPool.ParseConfig(connStr)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
}

func TestRealPgxPool_NewWithConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("skipando teste de integração em modo short")
	}

	connStr := os.Getenv("DB_CONN_URL")
	if connStr == "" {
		t.Skip("DB_CONN_URL não definido no .env ou no CI")
	}

	realPool := &RealPgxPool{}

	cfg, err := realPool.ParseConfig(connStr)
	require.NoError(t, err)

	pool, err := realPool.NewWithConfig(context.Background(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, pool)

	err = pool.Ping(context.Background())
	assert.NoError(t, err)

	pool.Close()
}
