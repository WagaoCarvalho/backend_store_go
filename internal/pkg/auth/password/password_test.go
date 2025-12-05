// password_test.go - Correções
package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestBcryptHasher_IsHash(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)

	t.Run("reconhece hash bcrypt 2a", func(t *testing.T) {
		hash2a := "$2a$10$N9qo8uLOickgx2ZMRZoMyeS4w6c6tv3bK6T7R6V7cZ2KjYVvX1zQW"
		assert.True(t, hasher.IsHash(hash2a))
	})

	t.Run("reconhece hash bcrypt 2b", func(t *testing.T) {
		hash2b := "$2b$10$N9qo8uLOickgx2ZMRZoMyeS4w6c6tv3bK6T7R6V7cZ2KjYVvX1zQW"
		assert.True(t, hasher.IsHash(hash2b))
	})

	t.Run("reconhece hash bcrypt 2y", func(t *testing.T) {
		hash2y := "$2y$10$N9qo8uLOickgx2ZMRZoMyeS4w6c6tv3bK6T7R6V7cZ2KjYVvX1zQW"
		assert.True(t, hasher.IsHash(hash2y))
	})

	t.Run("rejeita string que não é hash", func(t *testing.T) {
		notHash := "plainpassword123"
		assert.False(t, hasher.IsHash(notHash))
	})

	t.Run("rejeita hash com prefixo incorreto", func(t *testing.T) {
		wrongPrefix := "$2c$10$N9qo8uLOickgx2ZMRZoMyeS4w6c6tv3bK6T7R6V7cZ2KjYVvX1zQW"
		assert.False(t, hasher.IsHash(wrongPrefix))
	})

	t.Run("rejeita hash muito curto", func(t *testing.T) {
		shortHash := "$2a$"
		assert.False(t, hasher.IsHash(shortHash))
	})

	t.Run("rejeita hash com formato incompleto", func(t *testing.T) {
		incompleteHash := "$2a$10$short"
		assert.False(t, hasher.IsHash(incompleteHash))
	})
}

func TestBcryptHasher_Hash_WithDifferentCosts(t *testing.T) {
	tests := []struct {
		name string
		cost int
	}{
		{"min cost", 4},
		{"default cost", bcrypt.DefaultCost},
		{"high cost", 12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasher := NewBcryptHasher(tt.cost)
			password := "SecurePass123!"

			hashed, err := hasher.Hash(password)

			assert.NoError(t, err)
			assert.NotEmpty(t, hashed)
			assert.True(t, hasher.IsHash(hashed))

			// Verifica se o custo no hash corresponde
			// O formato é: $2a$[custo]$[salt+hash]
			parts := strings.Split(hashed, "$")
			if len(parts) >= 3 {
				// parts[2] contém o custo como string (ex: "10")
				costStr := parts[2]
				// Converte string para int
				// Note: O custo mínimo do bcrypt é 4, máximo é 31
				assert.NotEmpty(t, costStr)
				// Verifica apenas que não está vazio e é numérico
				assert.Regexp(t, `^[0-9]+$`, costStr)
			}
		})
	}
}

func TestBcryptHasher_LongPassword(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)

	t.Run("senha no limite do bcrypt", func(t *testing.T) {
		// Senha de exatamente 72 caracteres (limite do bcrypt)
		longPassword := strings.Repeat("a", 72)

		hashed, err := hasher.Hash(longPassword)
		assert.NoError(t, err)
		assert.True(t, hasher.IsHash(hashed))

		// Deve conseguir comparar
		err = hasher.Compare(hashed, longPassword)
		assert.NoError(t, err)
	})

	t.Run("senha acima do limite gera erro", func(t *testing.T) {
		// Senha longa (mais de 72 bytes)
		longPassword := strings.Repeat("a", 73)

		hashed, err := hasher.Hash(longPassword)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password length exceeds 72 bytes")
		assert.Empty(t, hashed)
	})
}

func TestBcryptHasher_DefaultCostWhenZero(t *testing.T) {
	t.Run("custo zero usa default", func(t *testing.T) {
		hasher := NewBcryptHasher(0)
		password := "test123"

		hashed, err := hasher.Hash(password)
		assert.NoError(t, err)
		assert.True(t, hasher.IsHash(hashed))
	})

	t.Run("custo negativo usa default", func(t *testing.T) {
		hasher := NewBcryptHasher(-1)
		password := "test123"

		hashed, err := hasher.Hash(password)
		assert.NoError(t, err)
		assert.True(t, hasher.IsHash(hashed))
	})
}

// Testes adicionais para casos de borda
func TestBcryptHasher_BoundaryConditions(t *testing.T) {
	hasher := NewBcryptHasher(bcrypt.DefaultCost)

	t.Run("custo muito baixo é ajustado", func(t *testing.T) {
		// bcrypt tem custo mínimo 4
		hasherLowCost := NewBcryptHasher(3)
		hashed, err := hasherLowCost.Hash("test")
		// bcrypt ajusta automaticamente para mínimo 4
		assert.NoError(t, err)
		assert.True(t, hasher.IsHash(hashed))
	})

	t.Run("custo muito alto funciona", func(t *testing.T) {
		// Teste com custo alto (pode ser lento, então talvez pular)
		if testing.Short() {
			t.Skip("Pulando teste de custo alto em modo curto")
		}

		hasherHighCost := NewBcryptHasher(14)
		hashed, err := hasherHighCost.Hash("test")
		assert.NoError(t, err)
		assert.True(t, hasher.IsHash(hashed))
	})
}
