package model

import (
	"errors"
	"fmt"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	modelsSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
)

type SupplierFull struct {
	Supplier   *modelsSupplier.Supplier
	Categories []modelsSupplierCategories.SupplierCategory
	Address    *modelsAddress.Address
	Contact    *modelsContact.Contact
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
