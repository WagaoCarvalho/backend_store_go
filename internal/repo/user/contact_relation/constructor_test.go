package repo

import (
	"testing"

	mockDb "github.com/WagaoCarvalho/backend_store_go/infra/mock/db"
	"github.com/stretchr/testify/assert"
)

func TestNewUserContactRelation(t *testing.T) {
	t.Run("successfully create new user contact instance", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserContactRelation(mockDB)

		assert.NotNil(t, result)
	})

	t.Run("return instance with provided db executor", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		result := NewUserContactRelation(mockDB)

		assert.NotNil(t, result)
	})

	t.Run("return different instances for different calls", func(t *testing.T) {
		mockDB := new(mockDb.MockDatabase)

		instance1 := NewUserContactRelation(mockDB)
		instance2 := NewUserContactRelation(mockDB)

		assert.NotSame(t, instance1, instance2)
		assert.NotNil(t, instance1)
		assert.NotNil(t, instance2)
	})
}
