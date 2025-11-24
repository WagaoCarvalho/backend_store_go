package repo

import (
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
)

func TestNewUserCategory(t *testing.T) {
	t.Run("successfully create new user category instance", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserCategory(mockDB)

		assert.NotNil(t, result)
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserCategory(mockDB)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		instance1 := NewUserCategory(mockDB)
		instance2 := NewUserCategory(mockDB)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}
