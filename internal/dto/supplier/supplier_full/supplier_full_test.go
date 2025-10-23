package dto

import (
	"testing"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoSupplier "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier"
	dtoSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_category"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	modelSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category"
	modelSupplierFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSupplierFullDTO_ToModel(t *testing.T) {

	// DTO de teste
	supplierDTO := &dtoSupplier.SupplierDTO{
		ID:   utils.Int64Ptr(1),
		Name: "Fornecedor X",
	}
	categoryDTO := dtoSupplierCategories.SupplierCategoryDTO{
		ID:   utils.Int64Ptr(10),
		Name: "Categoria Y",
	}
	addressDTO := &dtoAddress.AddressDTO{
		ID:     utils.Int64Ptr(100),
		Street: "Rua A",
		City:   "Cidade B",
	}
	contactDTO := &dtoContact.ContactDTO{
		ID:    utils.Int64Ptr(200),
		Email: "teste@email.com",
	}

	fullDTO := SupplierFullDTO{
		Supplier:   supplierDTO,
		Categories: []dtoSupplierCategories.SupplierCategoryDTO{categoryDTO},
		Address:    addressDTO,
		Contact:    contactDTO,
	}

	model := ToSupplierFullModel(fullDTO)

	assert.Equal(t, *fullDTO.Supplier.ID, model.Supplier.ID)
	assert.Equal(t, fullDTO.Supplier.Name, model.Supplier.Name)
	assert.Equal(t, *fullDTO.Categories[0].ID, model.Categories[0].ID)
	assert.Equal(t, fullDTO.Categories[0].Name, model.Categories[0].Name)
	assert.Equal(t, *fullDTO.Address.ID, model.Address.ID)
	assert.Equal(t, fullDTO.Address.Street, model.Address.Street)
	assert.Equal(t, *fullDTO.Contact.ID, model.Contact.ID)
	assert.Equal(t, fullDTO.Contact.Email, model.Contact.Email)
}

func TestToSupplierFullDTO(t *testing.T) {
	t.Run("Retorna vazio quando model é nil", func(t *testing.T) {
		result := ToSupplierFullDTO(nil)
		assert.Equal(t, SupplierFullDTO{}, result)
	})

	t.Run("Model preenchido", func(t *testing.T) {
		// Model de teste
		supplierModel := &modelSupplier.Supplier{
			ID:   1,
			Name: "Fornecedor X",
		}
		categoryModel := modelSupplierCategories.SupplierCategory{
			ID:   10,
			Name: "Categoria Y",
		}
		addressModel := &modelAddress.Address{
			ID:     100,
			Street: "Rua A",
			City:   "Cidade B",
		}
		contactModel := &modelContact.Contact{
			ID:    200,
			Email: "teste@email.com",
		}

		fullModel := &modelSupplierFull.SupplierFull{
			Supplier:   supplierModel,
			Categories: []modelSupplierCategories.SupplierCategory{categoryModel},
			Address:    addressModel,
			Contact:    contactModel,
		}

		dto := ToSupplierFullDTO(fullModel)

		assert.Equal(t, fullModel.Supplier.ID, *dto.Supplier.ID)
		assert.Equal(t, fullModel.Supplier.Name, dto.Supplier.Name)
		assert.Equal(t, fullModel.Categories[0].ID, *dto.Categories[0].ID)
		assert.Equal(t, fullModel.Categories[0].Name, dto.Categories[0].Name)
		assert.Equal(t, fullModel.Address.ID, *dto.Address.ID)
		assert.Equal(t, fullModel.Address.Street, dto.Address.Street)
		assert.Equal(t, fullModel.Contact.ID, *dto.Contact.ID)
		assert.Equal(t, fullModel.Contact.Email, dto.Contact.Email)
	})
}
