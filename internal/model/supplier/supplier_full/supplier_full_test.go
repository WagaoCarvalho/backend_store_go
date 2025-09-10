package model

import (
	"testing"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	modelsSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
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
				Supplier: &modelsSupplier.Supplier{},
			},
			wantErr: "endereço é obrigatório",
		},
		{
			name: "contato nulo",
			input: SupplierFull{
				Supplier: &modelsSupplier.Supplier{},
				Address:  &modelsAddress.Address{},
			},
			wantErr: "contato é obrigatório",
		},
		{
			name: "sem categorias",
			input: SupplierFull{
				Supplier: &modelsSupplier.Supplier{},
				Address:  &modelsAddress.Address{},
				Contact:  &modelsContact.Contact{},
			},
			wantErr: "pelo menos uma categoria é obrigatória",
		},
		{
			name: "fornecedor inválido",
			input: SupplierFull{
				// cria um Supplier vazio -> Validate dele deve falhar
				Supplier: &modelsSupplier.Supplier{},
				Address:  &modelsAddress.Address{},
				Contact:  &modelsContact.Contact{},
				Categories: []modelsSupplierCategories.SupplierCategory{
					{},
				},
			},
			wantErr: "fornecedor inválido",
		},
		{
			name: "tudo válido",
			input: SupplierFull{
				Supplier: &modelsSupplier.Supplier{
					// preencha com dados mínimos válidos
					Name: "Fornecedor Teste",
				},
				Address: &modelsAddress.Address{
					Street: "Rua X",
					City:   "São Paulo",
				},
				Contact: &modelsContact.Contact{
					ContactName: "Maria",
					Email:       "maria@example.com",
				},
				Categories: []modelsSupplierCategories.SupplierCategory{
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
