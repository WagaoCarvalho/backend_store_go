package repo

import (
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
)

func TestNewUserFilter(t *testing.T) {
	t.Run("successfully create new user filter instance", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserFilter(mockDB)

		assert.NotNil(t, result)
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserFilter(mockDB)

		assert.NotNil(t, result)
		// Verifica se o tipo interno está correto (opcional)
		repoImpl, ok := result.(*userFilterRepo)
		assert.True(t, ok)
		assert.Equal(t, mockDB, repoImpl.db)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		instance1 := NewUserFilter(mockDB)
		instance2 := NewUserFilter(mockDB)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})

	t.Run("return instance implementing UserFilter interface", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserFilter(mockDB)

		var _ UserFilter = result // Verificação em tempo de compilação
		assert.NotNil(t, result)
	})
}
