package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginCredentials_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   LoginCredentials
		wantErr string
	}{
		{
			name:    "email em branco",
			input:   LoginCredentials{Email: "", Password: "Senha123!"},
			wantErr: "Email",
		},
		{
			name:    "email inválido",
			input:   LoginCredentials{Email: "invalido.com", Password: "Senha123!"},
			wantErr: "Email",
		},
		{
			name:    "email muito longo",
			input:   LoginCredentials{Email: generateLongString(101), Password: "Senha123!"},
			wantErr: "Email",
		},
		{
			name:    "senha em branco",
			input:   LoginCredentials{Email: "user@example.com", Password: ""},
			wantErr: "Password",
		},
		{
			name:    "senha curta",
			input:   LoginCredentials{Email: "user@example.com", Password: "Ab1!"},
			wantErr: "Password",
		},
		{
			name:    "senha muito longa",
			input:   LoginCredentials{Email: "user@example.com", Password: generateLongString(65)},
			wantErr: "Password",
		},
		{
			name:    "senha sem maiúscula",
			input:   LoginCredentials{Email: "user@example.com", Password: "senha123!"},
			wantErr: "Password",
		},
		{
			name:    "senha sem minúscula",
			input:   LoginCredentials{Email: "user@example.com", Password: "SENHA123!"},
			wantErr: "Password",
		},
		{
			name:    "senha sem número",
			input:   LoginCredentials{Email: "user@example.com", Password: "Senha!@#"},
			wantErr: "Password",
		},
		{
			name:    "senha sem caractere especial",
			input:   LoginCredentials{Email: "user@example.com", Password: "Senha123"},
			wantErr: "Password",
		},
		{
			name:    "credenciais válidas",
			input:   LoginCredentials{Email: "user@example.com", Password: "Senha123!"},
			wantErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
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
