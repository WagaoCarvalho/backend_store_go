package models

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupplierCategory_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   SupplierCategory
		wantErr string
	}{
		{
			name:    "nome em branco",
			input:   SupplierCategory{Name: ""},
			wantErr: "Name",
		},
		{
			name:    "nome muito longo",
			input:   SupplierCategory{Name: generateLongString(101)},
			wantErr: "Name",
		},
		{
			name:    "descrição muito longa",
			input:   SupplierCategory{Name: "Categoria", Description: generateLongString(256)},
			wantErr: "Description",
		},
		{
			name:    "dados válidos",
			input:   SupplierCategory{Name: "Tecnologia", Description: "Componentes eletrônicos"},
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

func generateLongString(n int) string {
	return strings.Repeat("a", n)
}
