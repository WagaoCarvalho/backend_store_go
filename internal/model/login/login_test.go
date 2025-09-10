package model

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/stretchr/testify/assert"
)

func TestLoginCredentials_Validate(t *testing.T) {
	tests := []struct {
		name     string
		input    LoginCredentials
		wantErr  bool
		errField string
	}{
		{
			name:     "email em branco",
			input:    LoginCredentials{Email: "", Password: "Senha123!"},
			wantErr:  true,
			errField: "email",
		},
		{
			name:     "email inválido",
			input:    LoginCredentials{Email: "invalido.com", Password: "Senha123!"},
			wantErr:  true,
			errField: "email",
		},
		{
			name:     "email muito longo",
			input:    LoginCredentials{Email: generateLongString(101) + "@example.com", Password: "Senha123!"},
			wantErr:  true,
			errField: "email",
		},
		{
			name:     "senha em branco",
			input:    LoginCredentials{Email: "user@example.com", Password: ""},
			wantErr:  true,
			errField: "password",
		},
		{
			name:     "senha fraca",
			input:    LoginCredentials{Email: "user@example.com", Password: "abc"},
			wantErr:  true,
			errField: "password",
		},
		{
			name:    "credenciais válidas",
			input:   LoginCredentials{Email: "user@example.com", Password: "Senha123!"},
			wantErr: false,
		},
		{
			name: "senha fraca retorna ValidationError",
			input: LoginCredentials{
				Email:    "user@example.com",
				Password: "abc", // fraca o suficiente
			},
			wantErr:  true,
			errField: "password",
		},
		{
			name: "senha retorna erro genérico",
			input: LoginCredentials{
				Email:    "user@example.com",
				Password: "generic-error",
			},
			wantErr:  true,
			errField: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if !tt.wantErr {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)

			// Verifica se o erro é ValidationErrors
			var vErrs validators.ValidationErrors
			if ok := errors.As(err, &vErrs); ok {
				found := false
				for _, e := range vErrs {
					if strings.EqualFold(e.Field, tt.errField) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected error field %q in %v", tt.errField, vErrs)
			} else {
				// fallback: verificar string de erro
				assert.Contains(t, err.Error(), tt.errField)
			}
		})
	}
}

// Gera uma string com n caracteres
func generateLongString(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "a"
	}
	return s
}
