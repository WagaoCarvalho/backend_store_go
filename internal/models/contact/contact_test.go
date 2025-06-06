package models

import (
	"errors"
	"strings"
	"testing"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

func TestContact_Validate(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name    string
		contact Contact
		wantErr bool
		errType interface{}
		errMsg  string
	}{
		{
			name: "Valid contact",
			contact: Contact{
				UserID:      &userID,
				ContactName: "João Silva",
				Email:       "joao@email.com",
				Phone:       "(11) 1234-5678",
				Cell:        "(11) 91234-5678",
				ContactType: "financeiro",
			},
			wantErr: false,
		},
		{
			name: "Missing all IDs",
			contact: Contact{
				ContactName: "Contato",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "pelo menos um",
		},
		{
			name: "Blank ContactName",
			contact: Contact{
				UserID:      &userID,
				ContactName: " ",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "ContactName",
		},
		{
			name: "Short ContactName",
			contact: Contact{
				UserID:      &userID,
				ContactName: "AB",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "mínimo de 3",
		},
		{
			name: "Long ContactName",
			contact: Contact{
				UserID:      &userID,
				ContactName: strings.Repeat("A", 101),
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "máximo de 100",
		},
		{
			name: "Long ContactPosition",
			contact: Contact{
				UserID:          &userID,
				ContactName:     "Fulano",
				ContactPosition: strings.Repeat("X", 101),
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "ContactPosition",
		},
		{
			name: "Invalid email format",
			contact: Contact{
				UserID:      &userID,
				ContactName: "Fulano",
				Email:       "email@invalido",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "Email",
		},
		{
			name: "Email exceeds max length",
			contact: Contact{
				UserID:      &userID,
				ContactName: "Fulano",
				Email:       strings.Repeat("a", 95) + "@x.com",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "máximo de 100 caracteres",
		},
		{
			name: "Invalid phone format",
			contact: Contact{
				UserID:      &userID,
				ContactName: "Fulano",
				Phone:       "11987654321",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "Phone",
		},
		{
			name: "Invalid cell format",
			contact: Contact{
				UserID:      &userID,
				ContactName: "Fulano",
				Cell:        "(11) 1234-5678", // fixo no lugar do celular
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "Cell",
		},
		{
			name: "Invalid contact type",
			contact: Contact{
				UserID:      &userID,
				ContactName: "Fulano",
				ContactType: "RH",
			},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "tipo inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contact.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("esperado erro: %v, recebido: %v", tt.wantErr, err)
				return
			}
			if tt.wantErr {
				if tt.errType != nil {
					if !errors.As(err, &tt.errType) {
						t.Errorf("tipo de erro esperado: %T, recebido: %T", tt.errType, err)
					}
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("esperava mensagem contendo %q, recebeu %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}
