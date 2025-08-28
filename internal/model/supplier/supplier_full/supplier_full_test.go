package models

import (
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	models_supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	"github.com/stretchr/testify/assert"
)

func TestSupplierFull_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   SupplierFull
		wantErr string
	}{
		{
			name:    "fornecedor nulo",
			input:   SupplierFull{},
			wantErr: "fornecedor é obrigatório",
		},
		{
			name: "endereço nulo",
			input: SupplierFull{
				Supplier: &models_supplier.Supplier{},
			},
			wantErr: "endereço é obrigatório",
		},
		{
			name: "contato nulo",
			input: SupplierFull{
				Supplier: &models_supplier.Supplier{},
				Address:  &models_address.Address{},
			},
			wantErr: "contato é obrigatório",
		},
		{
			name: "sem categorias",
			input: SupplierFull{
				Supplier: &models_supplier.Supplier{},
				Address:  &models_address.Address{},
				Contact:  &models_contact.Contact{},
			},
			wantErr: "pelo menos uma categoria é obrigatória",
		},
		{
			name: "fornecedor inválido",
			input: SupplierFull{
				// cria um Supplier vazio -> Validate dele deve falhar
				Supplier: &models_supplier.Supplier{},
				Address:  &models_address.Address{},
				Contact:  &models_contact.Contact{},
				Categories: []models_supplier_categories.SupplierCategory{
					{},
				},
			},
			wantErr: "fornecedor inválido",
		},
		{
			name: "tudo válido",
			input: SupplierFull{
				Supplier: &models_supplier.Supplier{
					// preencha com dados mínimos válidos
					Name: "Fornecedor Teste",
				},
				Address: &models_address.Address{
					Street: "Rua X",
					City:   "São Paulo",
				},
				Contact: &models_contact.Contact{
					ContactName: "Maria",
					Email:       "maria@example.com",
				},
				Categories: []models_supplier_categories.SupplierCategory{
					{Name: "Categoria Teste"},
				},
			},
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
