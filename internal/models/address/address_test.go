package models

import (
	"errors"
	"strings"
	"testing"

	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

func TestAddress_Validate(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name    string
		address Address
		wantErr bool
		errType any
		errMsg  string
	}{
		{
			name: "Valid address with UserID",
			address: Address{
				UserID:     &userID,
				Street:     "Rua 1",
				City:       "Cidade",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345-678",
			},
			wantErr: false,
		},
		{
			name:    "Missing all IDs",
			address: Address{Street: "Rua 1", City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "pelo menos um deve ser informado",
		},
		{
			name:    "Blank Street",
			address: Address{UserID: &userID, Street: "   ", City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "Street",
		},
		{
			name:    "Street too short",
			address: Address{UserID: &userID, Street: "Ru", City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "mínimo de 3 caracteres",
		},
		{
			name:    "Street too long",
			address: Address{UserID: &userID, Street: strings.Repeat("a", 101), City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "máximo de 100 caracteres",
		},
		{
			name:    "Missing City",
			address: Address{UserID: &userID, Street: "Rua 1", City: "", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "City",
		},
		{
			name:    "City too short",
			address: Address{UserID: &userID, Street: "Rua 1", City: "A", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "mínimo de 2 caracteres",
		},
		{
			name:    "Blank State",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "   ", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "State",
		},
		{
			name:    "Invalid State",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "XX", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "estado inválido",
		},
		{
			name:    "Missing Country",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "Country",
		},
		{
			name:    "Unsupported Country",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", Country: "Argentina", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "país não suportado",
		},
		{
			name:    "Missing PostalCode",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", Country: "Brasil"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "PostalCode",
		},
		{
			name:    "Invalid PostalCode",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "ABC"},
			wantErr: true,
			errType: &utils_errors.ValidationError{},
			errMsg:  "formato inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				return
			}
			if tt.wantErr {
				if tt.errType != nil {
					if !errors.As(err, &tt.errType) {
						t.Errorf("expected error type %T, got %T", tt.errType, err)
					}
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected message to contain %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}
