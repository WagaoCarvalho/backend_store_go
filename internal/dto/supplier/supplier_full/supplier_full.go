package dto

import (
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	dtoSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_category" // usar DTO
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	modelsSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category"
	modelsSupplierFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
)

type SupplierFullDTO struct {
	Supplier   *dtoSupplier.SupplierDTO                    `json:"supplier,omitempty"`
	Categories []dtoSupplierCategories.SupplierCategoryDTO `json:"categories,omitempty"`
	Address    *dtoAddress.AddressDTO                      `json:"address,omitempty"`
	Contact    *dtoContact.ContactDTO                      `json:"contact,omitempty"`
}

// Converte DTO para Model
func ToSupplierFullModel(dto SupplierFullDTO) *modelsSupplierFull.SupplierFull {
	var categories []modelsSupplierCategories.SupplierCategory
	for _, c := range dto.Categories {
		if model := dtoSupplierCategories.ToSupplierCategoryModel(c); model != nil {
			categories = append(categories, *model)
		}
	}

	var supplier *modelsSupplier.Supplier
	if dto.Supplier != nil {
		supplier = dtoSupplier.ToSupplierModel(*dto.Supplier)
	}

	var address *modelAddress.Address
	if dto.Address != nil {
		address = dtoAddress.ToAddressModel(*dto.Address)
	}

	var contact *modelContact.Contact
	if dto.Contact != nil {
		contact = dtoContact.ToContactModel(*dto.Contact)
	}

	return &modelsSupplierFull.SupplierFull{
		Supplier:   supplier,
		Categories: categories,
		Address:    address,
		Contact:    contact,
	}
}

// Converte Model para DTO
func ToSupplierFullDTO(model *modelsSupplierFull.SupplierFull) SupplierFullDTO {
	if model == nil {
		return SupplierFullDTO{}
	}

	var categoriesDTO []dtoSupplierCategories.SupplierCategoryDTO
	if model.Categories != nil {
		categoriesDTO = make([]dtoSupplierCategories.SupplierCategoryDTO, len(model.Categories))
		for i, c := range model.Categories {
			categoriesDTO[i] = dtoSupplierCategories.ToSupplierCategoryDTO(&c)
		}
	}

	var supplierDTO *dtoSupplier.SupplierDTO
	if model.Supplier != nil {
		tmp := dtoSupplier.ToSupplierDTO(model.Supplier)
		supplierDTO = &tmp
	}

	var addressDTO *dtoAddress.AddressDTO
	if model.Address != nil {
		tmp := dtoAddress.ToAddressDTO(model.Address)
		addressDTO = &tmp
	}

	var contactDTO *dtoContact.ContactDTO
	if model.Contact != nil {
		tmp := dtoContact.ToContactDTO(model.Contact)
		contactDTO = &tmp
	}

	return SupplierFullDTO{
		Supplier:   supplierDTO,
		Categories: categoriesDTO,
		Address:    addressDTO,
		Contact:    contactDTO,
	}
}
