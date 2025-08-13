package models

import (
	"errors"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	models_supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
)

type SupplierFull struct {
	Supplier   *models_supplier.Supplier                     `json:"supplier"`
	Categories []models_supplier_categories.SupplierCategory `json:"categories"`
	Address    *models_address.Address                       `json:"address"`
	Contact    *models_contact.Contact                       `json:"contact"`
}

func (uf *SupplierFull) Validate() error {
	if uf.Supplier == nil {
		return errors.New("fornecedor é obrigatório")
	}

	if uf.Address == nil {
		return errors.New("endereço é obrigatório")
	}

	if uf.Contact == nil {
		return errors.New("contato é obrigatório")
	}

	if len(uf.Categories) == 0 {
		return errors.New("pelo menos uma categoria é obrigatória")
	}

	if err := uf.Supplier.Validate(); err != nil {
		return fmt.Errorf("fornecedor inválido: %w", err)
	}

	return nil
}
