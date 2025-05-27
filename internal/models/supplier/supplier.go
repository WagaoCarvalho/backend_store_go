package models

import (
	"time"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
)

type Supplier struct {
	ID          int64                                         `json:"id"`
	Name        string                                        `json:"name"`
	CNPJ        *string                                       `json:"cnpj,omitempty"`
	CPF         *string                                       `json:"cpf,omitempty"`
	ContactInfo string                                        `json:"contact_info"`
	Version     int                                           `json:"version"`
	CreatedAt   time.Time                                     `json:"created_at"`
	UpdatedAt   time.Time                                     `json:"updated_at"`
	Categories  []models_supplier_categories.SupplierCategory `json:"categories,omitempty"`
	Address     *models_address.Address                       `json:"address,omitempty"`
	Contact     *models_contact.Contact                       `json:"contact,omitempty"`
}
