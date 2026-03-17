package dto

import (
	"testing"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserFilterDTO_ToModel(t *testing.T) {
	t.Run("Converte todos os campos preenchidos corretamente", func(t *testing.T) {
		createdFrom := "2024-01-01"
		createdTo := "2024-12-31"
		updatedFrom := "2024-06-01"
		updatedTo := "2024-06-30"
		status := true

		dto := UserFilterDTO{
			Username:    "john_doe",
			Email:       "john@example.com",
			Status:      &status,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
			Limit:       20,
			Offset:      10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "john_doe", model.Username)
		assert.Equal(t, "john@example.com", model.Email)
		assert.Equal(t, &status, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 20, Offset: 10}, model.BaseFilter)

		require.NotNil(t, model.CreatedFrom)
		require.NotNil(t, model.CreatedTo)
		require.NotNil(t, model.UpdatedFrom)
		require.NotNil(t, model.UpdatedTo)

		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-12-31", model.CreatedTo.Format("2006-01-02"))
		assert.Equal(t, "2024-06-01", model.UpdatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-06-30", model.UpdatedTo.Format("2006-01-02"))
	})

	t.Run("Converte apenas campos básicos", func(t *testing.T) {
		dto := UserFilterDTO{
			Username: "john_doe",
			Email:    "john@example.com",
			Limit:    10,
			Offset:   0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "john_doe", model.Username)
		assert.Equal(t, "john@example.com", model.Email)
		assert.Nil(t, model.Status)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)
	})

	t.Run("Converte apenas status", func(t *testing.T) {
		status := false

		dto := UserFilterDTO{
			Status: &status,
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Empty(t, model.Username)
		assert.Empty(t, model.Email)
		assert.Equal(t, &status, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)
	})

	t.Run("Converte apenas username", func(t *testing.T) {
		dto := UserFilterDTO{
			Username: "john_doe",
			Limit:    10,
			Offset:   0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "john_doe", model.Username)
		assert.Empty(t, model.Email)
		assert.Nil(t, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)
	})

	t.Run("Converte apenas email", func(t *testing.T) {
		dto := UserFilterDTO{
			Email:  "john@example.com",
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Empty(t, model.Username)
		assert.Equal(t, "john@example.com", model.Email)
		assert.Nil(t, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)
	})

	t.Run("Erro com created_from após created_to", func(t *testing.T) {
		createdFrom := "2024-12-31"
		createdTo := "2024-01-01"

		dto := UserFilterDTO{
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "created_from")
	})

	t.Run("Erro com updated_from após updated_to", func(t *testing.T) {
		updatedFrom := "2024-12-31"
		updatedTo := "2024-01-01"

		dto := UserFilterDTO{
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "updated_from")
	})

	t.Run("Erro com limit inválido (menor que 1)", func(t *testing.T) {
		dto := UserFilterDTO{
			Limit: 0,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "limit")
	})

	t.Run("Erro com offset negativo", func(t *testing.T) {
		dto := UserFilterDTO{
			Limit:  10,
			Offset: -1,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "offset")
	})

	t.Run("Erro quando limit é maior que 100", func(t *testing.T) {
		dto := UserFilterDTO{
			Limit: 101,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "limit")
		assert.Contains(t, err.Error(), "máximo é 100")
	})

	t.Run("Erro ao converter created_from com formato inválido", func(t *testing.T) {
		invalidDate := "31-12-2024"

		dto := UserFilterDTO{
			CreatedFrom: &invalidDate,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'created_from'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro ao converter created_to com formato inválido", func(t *testing.T) {
		invalidDate := "2024/12/31"

		dto := UserFilterDTO{
			CreatedTo: &invalidDate,
			Limit:     10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'created_to'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro ao converter updated_from com formato inválido", func(t *testing.T) {
		invalidDate := "01-06-2024"

		dto := UserFilterDTO{
			UpdatedFrom: &invalidDate,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'updated_from'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro ao converter updated_to com formato inválido", func(t *testing.T) {
		invalidDate := "30/06/2024"

		dto := UserFilterDTO{
			UpdatedTo: &invalidDate,
			Limit:     10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'updated_to'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro com email inválido", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "invalid-email",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com username muito curto", func(t *testing.T) {
		dto := UserFilterDTO{
			Username: "ab",
			Limit:    10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "username")
		assert.Contains(t, err.Error(), "mínimo 3 caracteres")
	})

	t.Run("Email com subdomínio é válido", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@sub.example.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)
		assert.Equal(t, "user@sub.example.com", model.Email)
	})

	t.Run("Email com plus sign é válido", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user+label@example.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)
		assert.Equal(t, "user+label@example.com", model.Email)
	})

	t.Run("Username com tamanho mínimo é válido", func(t *testing.T) {
		dto := UserFilterDTO{
			Username: "abc",
			Limit:    10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)
		assert.Equal(t, "abc", model.Username)
	})

	t.Run("Username com caracteres especiais é válido", func(t *testing.T) {
		dto := UserFilterDTO{
			Username: "john.doe_123",
			Limit:    10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)
		assert.Equal(t, "john.doe_123", model.Username)
	})

	t.Run("Campos de data opcionais podem ser nil", func(t *testing.T) {
		dto := UserFilterDTO{
			Username:    "john_doe",
			Email:       "john@example.com",
			CreatedFrom: nil,
			CreatedTo:   nil,
			UpdatedFrom: nil,
			UpdatedTo:   nil,
			Limit:       10,
			Offset:      0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "john_doe", model.Username)
		assert.Equal(t, "john@example.com", model.Email)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
	})
}

// Testes específicos para a função isValidEmailFormat
func TestIsValidEmailFormat(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		// Casos que retornam false
		{
			name:     "Email vazio",
			email:    "",
			expected: false,
		},
		{
			name:     "Email com @ no início",
			email:    "@example.com",
			expected: false,
		},
		{
			name:     "Email com @ no fim",
			email:    "user@",
			expected: false,
		},
		{
			name:     "Email com ponto imediatamente após @",
			email:    "user@.com",
			expected: false,
		},
		{
			name:     "Email com ponto no final após @",
			email:    "user@example.",
			expected: false,
		},
		{
			name:     "Email com múltiplos @",
			email:    "user@name@example.com",
			expected: false,
		},
		{
			name:     "Email sem @",
			email:    "usermail.com",
			expected: false,
		},
		{
			name:     "Email sem ponto após @",
			email:    "user@example",
			expected: false,
		},

		// Casos que retornam true
		{
			name:     "Email válido simples",
			email:    "user@example.com",
			expected: true,
		},
		{
			name:     "Email com múltiplos pontos no local part",
			email:    "user.name@example.com",
			expected: true,
		},
		{
			name:     "Email com domínio de múltiplos níveis",
			email:    "user@example.co.uk",
			expected: true,
		},
		{
			name:     "Email com números",
			email:    "user123@example123.com",
			expected: true,
		},
		{
			name:     "Email com underscore",
			email:    "user_name@example.com",
			expected: true,
		},
		{
			name:     "Email com hífen",
			email:    "user-name@example.com",
			expected: true,
		},
		{
			name:     "Email com plus sign",
			email:    "user+label@example.com",
			expected: true,
		},
		{
			name:     "Email com um caractere antes do @",
			email:    "u@example.com",
			expected: true,
		},
		{
			name:     "Email com domínio mínimo",
			email:    "user@e.com",
			expected: true,
		},
		{
			name:     "Email com underline antes do @",
			email:    "user_name@example.com",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEmailFormat(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Testes específicos para validações de email no DTO
func TestUserFilterDTO_EmailValidations(t *testing.T) {
	t.Run("Erro com email começando com @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "@example.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email terminando com @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email contendo ponto imediatamente após @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email contendo ponto no final", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@example.",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email contendo múltiplos @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@name@example.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email sem @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "usermail.com",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Erro com email sem ponto após @", func(t *testing.T) {
		dto := UserFilterDTO{
			Email: "user@example",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "formato inválido")
	})

	t.Run("Email vazio é permitido", func(t *testing.T) {
		dto := UserFilterDTO{
			Email:  "",
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)
		assert.Empty(t, model.Email)
	})
}

// Testes de casos de borda para isValidEmailFormat
func TestIsValidEmailFormat_EdgeCases(t *testing.T) {
	t.Run("Email com apenas um caractere antes do @", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("a@b.com"))
	})

	t.Run("Email com domínio de um caractere", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user@x.com"))
	})

	t.Run("Email com muitos pontos no domínio", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user@example.co.uk.br"))
	})

	t.Run("Email com caracteres especiais no local part", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user.name+label_test@example.com"))
	})

	t.Run("Email com ponto antes do @", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user.name@example.com"))
	})

	t.Run("Email com underline antes do @", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user_name@example.com"))
	})

	t.Run("Email com hífen no domínio", func(t *testing.T) {
		assert.True(t, isValidEmailFormat("user@my-example.com"))
	})
}
