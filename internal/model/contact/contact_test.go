package models

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
)

func TestContact_Validate(t *testing.T) {
	id := int64(1)

	tests := []struct {
		name    string
		contact Contact
		wantErr bool
		errType any
		errMsg  string
	}{
		{
			name: "Valid contact with UserID",
			contact: Contact{
				UserID:      &id,
				ContactName: "João Silva",
				Email:       "joao@email.com",
				Phone:       "(11) 1234-5678",
				Cell:        "(11) 91234-5678",
				ContactType: "financeiro",
			},
			wantErr: false,
		},
		{
			name: "Valid contact with ClientID",
			contact: Contact{
				ClientID:    &id,
				ContactName: "João Silva",
				Email:       "joao@email.com",
				Phone:       "(11) 1234-5678",
				Cell:        "(11) 91234-5678",
				ContactType: "suporte",
			},
			wantErr: false,
		},
		{
			name: "Valid contact with SupplierID",
			contact: Contact{
				SupplierID:  &id,
				ContactName: "João Silva",
				Email:       "joao@email.com",
				Phone:       "(11) 1234-5678",
				Cell:        "(11) 91234-5678",
				ContactType: "comercial",
			},
			wantErr: false,
		},
		{
			name: "Missing all IDs",
			contact: Contact{
				ContactName: "Contato",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "exatamente um deve ser informado",
		},
		{
			name: "More than one ID provided (UserID + ClientID)",
			contact: Contact{
				UserID:      &id,
				ClientID:    &id,
				ContactName: "Contato",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "exatamente um deve ser informado",
		},
		{
			name: "More than one ID provided (UserID + SupplierID)",
			contact: Contact{
				UserID:      &id,
				SupplierID:  &id,
				ContactName: "Contato",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "exatamente um deve ser informado",
		},
		{
			name: "More than one ID provided (ClientID + SupplierID)",
			contact: Contact{
				ClientID:    &id,
				SupplierID:  &id,
				ContactName: "Contato",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "exatamente um deve ser informado",
		},
		{
			name: "Blank ContactName",
			contact: Contact{
				UserID:      &id,
				ContactName: " ",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "ContactName",
		},
		{
			name: "Short ContactName",
			contact: Contact{
				UserID:      &id,
				ContactName: "AB",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "mínimo de 3",
		},
		{
			name: "Long ContactName",
			contact: Contact{
				UserID:      &id,
				ContactName: strings.Repeat("A", 101),
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "máximo de 100",
		},
		{
			name: "Long ContactPosition",
			contact: Contact{
				UserID:          &id,
				ContactName:     "Fulano",
				ContactPosition: strings.Repeat("X", 101),
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "ContactPosition",
		},
		{
			name: "Invalid email format",
			contact: Contact{
				UserID:      &id,
				ContactName: "Fulano",
				Email:       "email@invalido",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "Email",
		},
		{
			name: "Email exceeds max length",
			contact: Contact{
				UserID:      &id,
				ContactName: "Fulano",
				Email:       strings.Repeat("a", 95) + "@x.com",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "máximo de 100 caracteres",
		},
		{
			name: "Invalid phone format",
			contact: Contact{
				UserID:      &id,
				ContactName: "Fulano",
				Phone:       "11987654321",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "Phone",
		},
		{
			name: "Invalid cell format",
			contact: Contact{
				UserID:      &id,
				ContactName: "Fulano",
				Cell:        "(11) 1234-5678", // fixo no lugar de celular
			},
			wantErr: true,
			errType: &validators.ValidationError{},
			errMsg:  "Cell",
		},
		{
			name: "Invalid contact type",
			contact: Contact{
				UserID:      &id,
				ContactName: "Fulano",
				ContactType: "RH",
			},
			wantErr: true,
			errType: &validators.ValidationError{},
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
