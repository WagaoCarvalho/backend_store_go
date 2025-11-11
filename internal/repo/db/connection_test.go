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
	// Salva o carregador original e restaura no final
	originalLoader := config.LoadDatabaseConfig
	defer func() { config.LoadDatabaseConfig = originalLoader }()

	// Força retorno de configuração sem URL
	config.LoadDatabaseConfig = func() config.Database {
		return config.Database{ConnURL: ""}
	}

	// Executa a função que deve falhar
	pool, err := Connect(nil)

	// Valida o comportamento esperado
	assert.Nil(t, pool, "pool deve ser nulo quando DB_CONN_URL não está definido")
	assert.ErrorIs(t, err, errMsg.ErrDBConnURLNotDefined, "erro esperado deve ser ErrDBConnURLNotDefined")
}

func TestRealPgxPool_ParseConfig(t *testing.T) {

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
	t.Run("successfully create new pool with config", func(t *testing.T) {
		realPool := &RealPgxPool{}
		ctx := context.Background()

		cfg, err := pgxpool.ParseConfig("postgres://user:pass@localhost:5432/testdb")
		if err != nil {
			t.Skip("Configuração de teste inválida, pulando teste")
		}

		pool, err := realPool.NewWithConfig(ctx, cfg)
		assert.NoError(t, err)
		assert.NotNil(t, pool)

		if pool != nil {
			pool.Close()
		}
	})

	t.Run("return error when config is nil", func(t *testing.T) {
		realPool := &RealPgxPool{}
		ctx := context.Background()

		pool, err := realPool.NewWithConfig(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, pool)
	})
}
