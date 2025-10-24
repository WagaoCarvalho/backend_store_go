package model

import (
	"errors"
	"fmt"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category"
	modelsUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
)

type UserFull struct {
	User       *modelsUser.User
	Categories []modelsUserCategories.UserCategory
	Address    *modelsAddress.Address
	Contact    *modelsContact.Contact
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
