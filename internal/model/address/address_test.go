package models

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
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
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: false,
		},
		{
			name: "Missing all IDs",
			address: Address{
				Street:       "Rua 1",
				StreetNumber: "123",
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgInvalidAssociation,
		},
		{
			name: "Blank Street",
			address: Address{
				UserID:       &userID,
				Street:       "   ",
				StreetNumber: "123",
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  "street",
		},
		{
			name: "Street too short",
			address: Address{
				UserID:       &userID,
				Street:       "Ru",
				StreetNumber: "123",
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgMin3,
		},
		{
			name: "Street too long",
			address: Address{
				UserID:       &userID,
				Street:       strings.Repeat("a", 101),
				StreetNumber: "123",
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgMax100,
		},
		{
			name: "StreetNumber too long",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: strings.Repeat("9", 21),
				Complement:   "Apto 1",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  "street_number",
		},
		{
			name: "Complement too long",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				Complement:   strings.Repeat("a", 256),
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  "complement",
		},
		{
			name: "Missing City",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  "city",
		},
		{
			name: "City too short",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "A",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgMin2,
		},
		{
			name: "Invalid State",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "Cidade",
				State:        "XX",
				Country:      "Brasil",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgInvalidState,
		},
		{
			name: "Unsupported Country",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "Cidade",
				State:        "SP",
				Country:      "Argentina",
				PostalCode:   "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgInvalidCountry,
		},
		{
			name: "Missing PostalCode",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  "postal_code",
		},
		{
			name: "Invalid PostalCode",
			address: Address{
				UserID:       &userID,
				Street:       "Rua 1",
				StreetNumber: "123",
				City:         "Cidade",
				State:        "SP",
				Country:      "Brasil",
				PostalCode:   "ABC",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgInvalidPostalCode,
		},
		{
			name: "Blank State",
			address: Address{
				UserID:     &userID,
				Street:     "Rua 1",
				City:       "Cidade",
				State:      "   ", // só espaços
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgRequiredField, // validar mensagem de campo obrigatório
		},
		{
			name: "Blank Country",
			address: Address{
				UserID:     &userID,
				Street:     "Rua 1",
				City:       "Cidade",
				State:      "SP",
				Country:    "   ", // só espaços
				PostalCode: "12345678",
			},
			wantErr: true,
			errType: validators.ValidationErrors{},
			errMsg:  validators.MsgRequiredField,
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
