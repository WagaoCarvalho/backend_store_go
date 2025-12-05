package middleware

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetRequestID(t *testing.T) {
	t.Run("set e get request ID válido", func(t *testing.T) {
		ctx := context.Background()
		expectedID := "req-12345"

		ctx = SetRequestID(ctx, expectedID)
		actualID := GetRequestID(ctx)

		assert.Equal(t, expectedID, actualID)
	})

	t.Run("get request ID de contexto vazio", func(t *testing.T) {
		ctx := context.Background()
		id := GetRequestID(ctx)

		assert.Empty(t, id)
	})

	t.Run("get request ID de contexto TODO", func(t *testing.T) {
		ctx := context.TODO()
		id := GetRequestID(ctx)

		assert.Empty(t, id)
	})

	t.Run("get request ID com contexto nil", func(t *testing.T) {
		// Cria uma variável nil explicitamente
		var nilCtx context.Context // zero value é nil

		id := GetRequestID(nilCtx)
		assert.Empty(t, id)
	})

	t.Run("set request ID múltiplas vezes", func(t *testing.T) {
		ctx := context.Background()

		ctx = SetRequestID(ctx, "first-id")
		ctx = SetRequestID(ctx, "second-id")

		id := GetRequestID(ctx)
		assert.Equal(t, "second-id", id)
	})

	t.Run("request ID vazio", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetRequestID(ctx, "")

		id := GetRequestID(ctx)
		assert.Empty(t, id)
	})

	t.Run("request ID com tipos diferentes no contexto", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), requestIDKey, 123) // Int ao invés de string
		ctx = SetRequestID(ctx, "string-id")                              // Sobrescreve

		id := GetRequestID(ctx)
		assert.Equal(t, "string-id", id)
	})
}

func TestSetAndGetUserID(t *testing.T) {
	t.Run("set e get user ID válido", func(t *testing.T) {
		ctx := context.Background()
		expectedID := "user-98765"

		ctx = SetUserID(ctx, expectedID)
		actualID := GetUserID(ctx)

		assert.Equal(t, expectedID, actualID)
	})

	t.Run("get user ID de contexto vazio", func(t *testing.T) {
		ctx := context.Background()
		id := GetUserID(ctx)

		assert.Empty(t, id)
	})

	t.Run("get user ID de contexto TODO", func(t *testing.T) {
		ctx := context.TODO()
		id := GetUserID(ctx)

		assert.Empty(t, id)
	})

	t.Run("set user ID múltiplas vezes", func(t *testing.T) {
		ctx := context.Background()

		ctx = SetUserID(ctx, "user-1")
		ctx = SetUserID(ctx, "user-2")

		id := GetUserID(ctx)
		assert.Equal(t, "user-2", id)
	})

	t.Run("get user ID com contexto nil", func(t *testing.T) {
		var nilCtx context.Context

		id := GetUserID(nilCtx)
		assert.Empty(t, id)
	})

	t.Run("user ID vazio", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetUserID(ctx, "")

		id := GetUserID(ctx)
		assert.Empty(t, id)
	})

	t.Run("user ID com tipos diferentes no contexto", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), userIDKey, 456) // Int ao invés de string
		ctx = SetUserID(ctx, "string-user-id")                         // Sobrescreve

		id := GetUserID(ctx)
		assert.Equal(t, "string-user-id", id)
	})
}

func TestContextKeysAreDifferent(t *testing.T) {
	t.Run("request ID e user ID não interferem", func(t *testing.T) {
		ctx := context.Background()

		ctx = SetRequestID(ctx, "req-id-123")
		ctx = SetUserID(ctx, "user-id-456")

		requestID := GetRequestID(ctx)
		userID := GetUserID(ctx)

		assert.Equal(t, "req-id-123", requestID)
		assert.Equal(t, "user-id-456", userID)
		assert.NotEqual(t, requestID, userID)
	})

	t.Run("set apenas request ID, user ID deve estar vazio", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetRequestID(ctx, "req-only")

		assert.Equal(t, "req-only", GetRequestID(ctx))
		assert.Empty(t, GetUserID(ctx))
	})

	t.Run("set apenas user ID, request ID deve estar vazio", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetUserID(ctx, "user-only")

		assert.Equal(t, "user-only", GetUserID(ctx))
		assert.Empty(t, GetRequestID(ctx))
	})
}

func TestContextPropagation(t *testing.T) {
	t.Run("contexto propagado entre funções", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetRequestID(ctx, "propagated-id")

		// Simula uma função que recebe contexto
		checkRequestID := func(ctx context.Context) string {
			return GetRequestID(ctx)
		}

		id := checkRequestID(ctx)
		assert.Equal(t, "propagated-id", id)
	})

	t.Run("contexto em goroutines", func(t *testing.T) {
		ctx := context.Background()
		ctx = SetUserID(ctx, "goroutine-user")

		resultChan := make(chan string, 1)

		go func(ctx context.Context) {
			resultChan <- GetUserID(ctx)
		}(ctx)

		userID := <-resultChan
		assert.Equal(t, "goroutine-user", userID)
	})
}

func TestTypeSafety(t *testing.T) {
	t.Run("GetRequestID retorna string vazia para tipo errado", func(t *testing.T) {
		// Contexto com valor de tipo errado para requestIDKey
		ctx := context.WithValue(context.Background(), requestIDKey, 12345)

		id := GetRequestID(ctx)
		assert.Empty(t, id) // Type assertion falha, retorna ""
	})

	t.Run("GetUserID retorna string vazia para tipo errado", func(t *testing.T) {
		// Contexto com valor de tipo errado para userIDKey
		ctx := context.WithValue(context.Background(), userIDKey, []byte("user"))

		id := GetUserID(ctx)
		assert.Empty(t, id) // Type assertion falha, retorna ""
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("contexto com múltiplos valores", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "other_key", "other_value")
		ctx = SetRequestID(ctx, "request-123")
		ctx = SetUserID(ctx, "user-456")

		// Deve ainda conseguir acessar os valores
		assert.Equal(t, "request-123", GetRequestID(ctx))
		assert.Equal(t, "user-456", GetUserID(ctx))
	})

	t.Run("IDs com caracteres especiais", func(t *testing.T) {
		ctx := context.Background()

		specialRequestID := "req-123_abc@DEF#456"
		specialUserID := "usr-789|xyz!123"

		ctx = SetRequestID(ctx, specialRequestID)
		ctx = SetUserID(ctx, specialUserID)

		assert.Equal(t, specialRequestID, GetRequestID(ctx))
		assert.Equal(t, specialUserID, GetUserID(ctx))
	})

	t.Run("IDs muito longos", func(t *testing.T) {
		ctx := context.Background()

		longRequestID := "req-" + string(make([]byte, 1000))
		longUserID := "usr-" + string(make([]byte, 1000))

		ctx = SetRequestID(ctx, longRequestID)
		ctx = SetUserID(ctx, longUserID)

		assert.Equal(t, longRequestID, GetRequestID(ctx))
		assert.Equal(t, longUserID, GetUserID(ctx))
	})
}

// Testes adicionais para garantir robustez
func TestMiddlewareFunctionsRobustness(t *testing.T) {
	t.Run("GetRequestID com contexto cancelado", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = SetRequestID(ctx, "canceled-context-id")
		cancel()

		// Ainda deve conseguir obter o ID mesmo com contexto cancelado
		id := GetRequestID(ctx)
		assert.Equal(t, "canceled-context-id", id)
	})

	t.Run("GetUserID com contexto com timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		ctx = SetUserID(ctx, "timeout-context-user")

		// Ainda deve conseguir obter o ID mesmo com timeout expirado
		id := GetUserID(ctx)
		assert.Equal(t, "timeout-context-user", id)
	})

	t.Run("SetRequestID com contexto de background válido", func(t *testing.T) {
		ctx := context.Background()
		newCtx := SetRequestID(ctx, "test-id")

		assert.NotNil(t, newCtx)
		assert.NotEqual(t, ctx, newCtx) // Deve retornar novo contexto
		assert.Equal(t, "test-id", GetRequestID(newCtx))
	})

	t.Run("SetUserID com contexto de TODO válido", func(t *testing.T) {
		ctx := context.TODO()
		newCtx := SetUserID(ctx, "todo-user")

		assert.NotNil(t, newCtx)
		assert.NotEqual(t, ctx, newCtx) // Deve retornar novo contexto
		assert.Equal(t, "todo-user", GetUserID(newCtx))
	})
}
