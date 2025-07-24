package models

import (
	"errors"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
)

type UserFull struct {
	User       *models_user.User                     `json:"user"`
	Categories []models_user_categories.UserCategory `json:"categories"`
	Address    *models_address.Address               `json:"address"`
	Contact    *models_contact.Contact               `json:"contact"`
}

func (uf *UserFull) Validate() error {
	// Validação em ordem segura
	if uf.User == nil {
		return errors.New("usuário é obrigatório")
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

	// Validação do User só depois de verificar que não é nil
	if err := uf.User.Validate(); err != nil {
		return fmt.Errorf("usuário inválido: %w", err)
	}

	return nil
}
