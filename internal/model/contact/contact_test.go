package model

import (
	"errors"
	"strings"
	"testing"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func TestContact_Validate(t *testing.T) {
	tests := []struct {
		name     string
		contact  Contact
		wantErr  bool
		errField string
		errMsg   string
	}{
		{
			name: "Valid contact",
			contact: Contact{
				ContactName:        "João Silva",
				ContactDescription: "Gerente de contas",
				Email:              "joao@example.com",
				Phone:              "1134567890",
				Cell:               "11998765432",
				ContactType:        "principal",
			},
			wantErr: false,
		},
		{
			name: "Blank ContactName",
			contact: Contact{
				Email: "teste@example.com",
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgRequiredField,
		},
		{
			name: "ContactName too short",
			contact: Contact{
				ContactName: "Jo",
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgMin3,
		},
		{
			name: "ContactName too long",
			contact: Contact{
				ContactName: strings.Repeat("a", 101),
			},
			wantErr:  true,
			errField: "contact_name",
			errMsg:   validators.MsgMax100,
		},
		{
			name: "ContactDescription too long",
			contact: Contact{
				ContactName:        "Pedro",
				ContactDescription: strings.Repeat("x", 101),
			},
			wantErr:  true,
			errField: "contact_description",
			errMsg:   validators.MsgMax100,
		},
		{
			name: "Invalid Email format",
			contact: Contact{
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
