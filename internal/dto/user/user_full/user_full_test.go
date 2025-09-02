package dto

import (
	"testing"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoUser "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	dtoUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_category"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	modelUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	modelUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUserFullDTO_ToModel(t *testing.T) {
	userDTO := &dtoUser.UserDTO{
		UID:      utils.Int64Ptr(1),
		Username: "usuarioX",
		Email:    "teste@email.com",
	}
	categoryDTO := dtoUserCategories.UserCategoryDTO{
		ID:   utils.UintPtr(10),
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

	fullDTO := UserFullDTO{
		User:       userDTO,
		Categories: []dtoUserCategories.UserCategoryDTO{categoryDTO},
		Address:    addressDTO,
		Contact:    contactDTO,
	}

	model := ToUserFullModel(fullDTO)

	assert.Equal(t, *fullDTO.User.UID, model.User.UID)
	assert.Equal(t, fullDTO.User.Username, model.User.Username)
	assert.Equal(t, fullDTO.User.Email, model.User.Email)
	assert.Equal(t, *fullDTO.Categories[0].ID, model.Categories[0].ID)
	assert.Equal(t, fullDTO.Categories[0].Name, model.Categories[0].Name)
	assert.Equal(t, *fullDTO.Address.ID, model.Address.ID)
	assert.Equal(t, fullDTO.Address.Street, model.Address.Street)
	assert.Equal(t, *fullDTO.Contact.ID, model.Contact.ID)
	assert.Equal(t, fullDTO.Contact.Email, model.Contact.Email)
}

func TestToUserFullDTO(t *testing.T) {
	t.Run("Retorna vazio quando model é nil", func(t *testing.T) {
		result := ToUserFullDTO(nil)
		assert.Equal(t, UserFullDTO{}, result)
	})

	t.Run("Model preenchido", func(t *testing.T) {
		userModel := &modelUser.User{
			UID:      1,
			Username: "usuarioX",
			Email:    "teste@email.com",
		}
		categoryModel := modelUserCategories.UserCategory{
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

		fullModel := &modelUserFull.UserFull{
			User:       userModel,
			Categories: []modelUserCategories.UserCategory{categoryModel},
			Address:    addressModel,
			Contact:    contactModel,
		}

		dto := ToUserFullDTO(fullModel)

		assert.Equal(t, fullModel.User.UID, *dto.User.UID)
		assert.Equal(t, fullModel.User.Username, dto.User.Username)
		assert.Equal(t, fullModel.User.Email, dto.User.Email)
		assert.Equal(t, fullModel.Categories[0].ID, *dto.Categories[0].ID)
		assert.Equal(t, fullModel.Categories[0].Name, dto.Categories[0].Name)
		assert.Equal(t, fullModel.Address.ID, *dto.Address.ID)
		assert.Equal(t, fullModel.Address.Street, dto.Address.Street)
		assert.Equal(t, fullModel.Contact.ID, *dto.Contact.ID)
		assert.Equal(t, fullModel.Contact.Email, dto.Contact.Email)
	})
}
