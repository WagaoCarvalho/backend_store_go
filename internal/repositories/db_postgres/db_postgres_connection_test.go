package repositories

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Teste para conexão bem-sucedida
func TestConnect_Success(t *testing.T) {
	// Carrega o .env localmente, mas ignora erros no GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") == "" {
		_ = godotenv.Load()
	}

	// Salva o valor original da variável de ambiente
	originalDBConnURL := os.Getenv("DB_CONN_URL")
	defer os.Setenv("DB_CONN_URL", originalDBConnURL) // Restaura no final

	// Define a variável de ambiente para testes
	os.Setenv("DB_CONN_URL", os.Getenv("DB_CONN_URL_TESTE"))

	// Simula uma conexão ao banco de dados
	mockPool := &MockPgxPool{}
	pool, err := Connect(mockPool)

	// Verifica se a conexão foi bem-sucedida
	assert.NoError(t, err)
	assert.NotNil(t, pool)
}

// Teste para verificar erro quando a variável de ambiente não está definida
func TestConnect_NoEnvVariable(t *testing.T) {
	os.Unsetenv("DB_CONN_URL")

	mockPool := &MockPgxPool{}
	pool, err := Connect(mockPool)

	assert.Error(t, err)
	assert.Nil(t, pool)
	assert.Equal(t, "variável de ambiente DB_CONN_URL não definida", err.Error())
}
