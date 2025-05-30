package models

import (
	"errors"
	"strings"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/utils"
)

func TestAddress_Validate(t *testing.T) {
	userID := int64(1)

	tests := []struct {
		name    string
		address Address
		wantErr bool
		errType interface{}
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
			errType: &utils.ValidationError{},
			errMsg:  "pelo menos um deve ser informado",
		},
		{
			name:    "Missing Street",
			address: Address{UserID: &userID, City: "Cidade", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils.ValidationError{},
			errMsg:  "Street",
		},
		{
			name:    "Missing City",
			address: Address{UserID: &userID, Street: "Rua 1", State: "SP", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils.ValidationError{},
			errMsg:  "City",
		},
		{
			name:    "Missing State",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", Country: "Brasil", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils.ValidationError{},
			errMsg:  "State",
		},
		{
			name:    "Missing Country",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", PostalCode: "12345-678"},
			wantErr: true,
			errType: &utils.ValidationError{},
			errMsg:  "Country",
		},
		{
			name:    "Missing PostalCode",
			address: Address{UserID: &userID, Street: "Rua 1", City: "Cidade", State: "SP", Country: "Brasil"},
			wantErr: true,
			errType: &utils.ValidationError{},
			errMsg:  "PostalCode",
		},
		{
			name: "Generic error on street",
			address: Address{
				UserID:     &userID,
				Street:     "cause_generic_error",
				City:       "Cidade",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345-678",
			},
			wantErr: true,
			errType: nil,
			errMsg:  "erro genérico na validação",
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
				// Verifica o tipo se for especificado
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
