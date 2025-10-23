package dto

import (
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoUser "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	dtoUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_category"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	modelsUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category"
	modelsUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
)

type UserFullDTO struct {
	User       *dtoUser.UserDTO                    `json:"user,omitempty"`
	Categories []dtoUserCategories.UserCategoryDTO `json:"categories,omitempty"`
	Address    *dtoAddress.AddressDTO              `json:"address,omitempty"`
	Contact    *dtoContact.ContactDTO              `json:"contact,omitempty"`
}

// Converte DTO para Model
func ToUserFullModel(dto UserFullDTO) *modelsUserFull.UserFull {
	var categories []modelsUserCategories.UserCategory
	for _, c := range dto.Categories {
		if model := dtoUserCategories.ToUserCategoryModel(c); model != nil {
			categories = append(categories, *model)
		}
	}

	var user *modelsUser.User
	if dto.User != nil {
		user = dtoUser.ToUserModel(*dto.User)
	}

	var address *modelAddress.Address
	if dto.Address != nil {
		address = dtoAddress.ToAddressModel(*dto.Address)
	}

	var contact *modelContact.Contact
	if dto.Contact != nil {
		contact = dtoContact.ToContactModel(*dto.Contact)
	}

	return &modelsUserFull.UserFull{
		User:       user,
		Categories: categories,
		Address:    address,
		Contact:    contact,
	}
}

// Converte Model para DTO
func ToUserFullDTO(model *modelsUserFull.UserFull) UserFullDTO {
	if model == nil {
		return UserFullDTO{}
	}

	var categoriesDTO []dtoUserCategories.UserCategoryDTO
	if model.Categories != nil {
		categoriesDTO = make([]dtoUserCategories.UserCategoryDTO, len(model.Categories))
		for i, c := range model.Categories {
			categoriesDTO[i] = dtoUserCategories.ToUserCategoryDTO(&c)
		}
	}

	var userDTO *dtoUser.UserDTO
	if model.User != nil {
		tmp := dtoUser.ToUserDTO(model.User)
		userDTO = &tmp
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

	return UserFullDTO{
		User:       userDTO,
		Categories: categoriesDTO,
		Address:    addressDTO,
		Contact:    contactDTO,
	}
}
