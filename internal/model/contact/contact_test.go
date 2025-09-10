package model

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func TestContact_Validate(t *testing.T) {
	id := int64(1)

	tests := []struct {
		name     string
		contact  Contact
		wantErr  bool
		errField string
		errMsg   string
	}{
		{
			name: "Valid contact with UserID",
			contact: Contact{
				UserID:      &id,
				ContactName: "João Silva",
				Email:       "joao@example.com",
				Phone:       "1134567890",
				Cell:        "11998765432",
				ContactType: "principal",
			},
			wantErr: false,
		},
		{
			name: "Invalid IDs (none provided)",
			contact: Contact{
				ContactName: "Maria",
			},
			wantErr:  true,
			errField: "UserID/ClientID/SupplierID",     // atualizado
			errMsg:   validators.MsgInvalidAssociation, // atualizado
		},

		{
			name: "Invalid IDs (multiple provided)",
			contact: Contact{
				UserID:      &id,
				ClientID:    &id,
				ContactName: "Carlos",
			},
			wantErr:  true,
			errField: "UserID/ClientID/SupplierID",
			errMsg:   validators.MsgInvalidAssociation,
		},
		{
			name: "Blank ContactName",
			contact: Contact{
				UserID: &id,
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgRequiredField,
		},
		{
			name: "ContactName too short",
			contact: Contact{
				UserID:      &id,
				ContactName: "Jo",
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgMin3,
		},
		{
			name: "ContactName too long",
			contact: Contact{
				UserID:      &id,
				ContactName: strings.Repeat("a", 101),
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgMax100,
		},
		{
			name: "ContactPosition too long",
			contact: Contact{
				UserID:          &id,
				ContactName:     "Pedro",
				ContactPosition: strings.Repeat("x", 101),
			},
			wantErr:  true,
			errField: "contact_position",
			errMsg:   validators.MsgMax100,
		},
		{
			name: "Invalid Email format",
			contact: Contact{
				UserID:      &id,
				ContactName: "José",
				Email:       "invalid-email",
			},
			wantErr:  true,
			errField: "email",
			errMsg:   validators.MsgInvalidFormat,
		},
		{
			name: "Email too long",
			contact: Contact{
				UserID:      &id,
				ContactName: "José",
				Email:       strings.Repeat("a", 101) + "@example.com",
			},
			wantErr:  true,
			errField: "email",
			errMsg:   validators.MsgMax100,
		},
		{
			name: "Invalid Phone format",
			contact: Contact{
				UserID:      &id,
				ContactName: "Ana",
				Phone:       "abc123",
			},
			wantErr:  true,
			errField: "phone",
			errMsg:   validators.MsgInvalidPhone,
		},
		{
			name: "Invalid Cell format",
			contact: Contact{
				UserID:      &id,
				ContactName: "Ana",
				Cell:        "12345",
			},
			wantErr:  true,
			errField: "cell",
			errMsg:   validators.MsgInvalidCell,
		},
		{
			name: "Invalid ContactType",
			contact: Contact{
				UserID:      &id,
				ContactName: "Roberto",
				ContactType: "gerente",
			},
			wantErr:  true,
			errField: "contact_type",
			errMsg:   validators.MsgInvalidType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.contact.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err)
				return
			}
			if tt.wantErr {
				var vErrs validators.ValidationErrors
				if !errors.As(err, &vErrs) {
					t.Errorf("expected ValidationErrors, got %T", err)
				}
				if !strings.Contains(err.Error(), tt.errField) {
					t.Errorf("expected error field %q, got %q", tt.errField, err.Error())
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected message %q, got %q", tt.errMsg, err.Error())
				}
			}
		})
	}
}
