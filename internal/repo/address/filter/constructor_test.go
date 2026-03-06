package repo

import (
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
)

func TestNewFilterAddress(t *testing.T) {
	t.Run("successfully create new address filter instance", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewFilterAddress(mockDB)

		assert.NotNil(t, result)
		assert.Implements(t, (*AddressFilter)(nil), result)
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewFilterAddress(mockDB)

		assert.NotNil(t, result)

		// Verifica se o db executor foi corretamente atribuído (acessando via reflexão ou getter)
		repo, ok := result.(*addressFilterRepo)
		assert.True(t, ok)
		assert.Equal(t, mockDB, repo.db)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		instance1 := NewFilterAddress(mockDB)
		instance2 := NewFilterAddress(mockDB)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)

		// Verifica que são instâncias diferentes mas apontam para o mesmo db
		repo1, ok1 := instance1.(*addressFilterRepo)
		repo2, ok2 := instance2.(*addressFilterRepo)

		assert.True(t, ok1)
		assert.True(t, ok2)
		assert.Equal(t, repo1.db, repo2.db)
	})

	t.Run("handle nil db executor gracefully", func(t *testing.T) {
		// O comportamento esperado depende da implementação
		// Se a função aceitar nil, este teste deve passar
		// Se não aceitar, deve panic - ajuste conforme necessário
		result := NewFilterAddress(nil)

		assert.NotNil(t, result)

		repo, ok := result.(*addressFilterRepo)
		assert.True(t, ok)
		assert.Nil(t, repo.db)
	})
}

// Teste adicional para garantir que o contrato da interface seja mantido
func TestAddressFilterRepo_ImplementsInterface(t *testing.T) {
	// Verificação em tempo de compilação
	var _ AddressFilter = (*addressFilterRepo)(nil)

	// Teste em tempo de execução
	mockDB := new(mockDb.MockDatabase)
	repo := NewFilterAddress(mockDB)

	assert.Implements(t, (*AddressFilter)(nil), repo)
}

// Teste de benchmark opcional
func BenchmarkNewFilterAddress(b *testing.B) {
	mockDB := new(mockDb.MockDatabase)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewFilterAddress(mockDB)
	}
}
