package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

// Erros personalizados

type ContactService interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}

type contactService struct {
	contactRepo repositories.ContactRepository
	logger      *logger.LoggerAdapter
}

func NewContactService(contactRepo repositories.ContactRepository, logger *logger.LoggerAdapter) ContactService {
	return &contactService{
		contactRepo: contactRepo,
		logger:      logger,
	}
}

func (s *contactService) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	s.logger.Info(ctx, "[contactService] - Iniciando criação de contato", map[string]interface{}{
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

	if contact.ContactName == "" {
		s.logger.Warn(ctx, "[ContactService] - Nome do contato é obrigatório", nil)
		return &models.Contact{}, ErrContactNameRequired
	}

	if contact.UserID == nil && contact.ClientID == nil && contact.SupplierID == nil {
		s.logger.Warn(ctx, "[ContactService] - Associação inválida: precisa de UserID, ClientID ou SupplierID", nil)
		return &models.Contact{}, ErrContactAssociationRequired
	}

	if contact.Email != "" && !validators.IsValidEmail(contact.Email) {
		s.logger.Warn(ctx, "[ContactService] - Email inválido", map[string]interface{}{
			"email": contact.Email,
		})
		return &models.Contact{}, ErrInvalidEmail
	}

	createdContact, err := s.contactRepo.Create(ctx, contact)
	if err != nil {
		s.logger.Error(ctx, err, "[ContactService] - Erro ao criar contato", map[string]interface{}{
			"contact_name": contact.ContactName,
			"email":        contact.Email,
		})
		return &models.Contact{}, fmt.Errorf("%w: %v", ErrCreateContact, err)
	}

	s.logger.Info(ctx, "[ContactService] - Contato criado com sucesso", map[string]interface{}{
		"contact_id":   createdContact.ID,
		"contact_name": createdContact.ContactName,
	})
	return createdContact, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	s.logger.Info(ctx, "[contactService] - Iniciando busca de contato", map[string]interface{}{
		"contact_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, "[contactService] - ID inválido para busca", map[string]interface{}{
			"contact_id": id,
		})
		return nil, ErrInvalidID
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Info(ctx, "[contactService] - Contato não encontrado", map[string]interface{}{
				"contact_id": id,
			})
			return nil, ErrContactNotFound
		}

		s.logger.Error(ctx, err, "[contactService] - Erro ao buscar contato por ID", map[string]interface{}{
			"contact_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	s.logger.Info(ctx, "[contactService] - Contato encontrado com sucesso", map[string]interface{}{
		"contact_id": id,
	})
	return contact, nil
}

func (s *contactService) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	s.logger.Info(ctx, "[contactService] - Iniciando busca de contatos por usuário", map[string]interface{}{
		"user_id": userID,
	})

	if userID <= 0 {
		s.logger.Error(ctx, ErrUserIDInvalid, "[contactService] - ID de usuário inválido", map[string]interface{}{
			"user_id": userID,
		})
		return nil, ErrUserIDInvalid
	}

	contacts, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, "[contactService] - Erro ao listar contatos por usuário", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListUserContacts, err)
	}

	s.logger.Info(ctx, "[contactService] - Contatos do usuário listados com sucesso", map[string]interface{}{
		"user_id":   userID,
		"qtd_total": len(contacts),
	})

	return contacts, nil
}

func (s *contactService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	s.logger.Info(ctx, "[contactService] - Iniciando busca de contatos por cliente", map[string]interface{}{
		"client_id": clientID,
	})

	if clientID <= 0 {
		s.logger.Error(ctx, ErrClientIDInvalid, "[contactService] - ID do cliente inválido", map[string]interface{}{
			"client_id": clientID,
		})
		return nil, ErrClientIDInvalid
	}

	contacts, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		s.logger.Error(ctx, err, "[contactService] - Erro ao listar contatos por cliente", map[string]interface{}{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListClientContacts, err)
	}

	s.logger.Info(ctx, "[contactService] - Contatos do cliente listados com sucesso", map[string]interface{}{
		"client_id": clientID,
		"qtd_total": len(contacts),
	})

	return contacts, nil
}

func (s *contactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	s.logger.Info(ctx, "[contactService] - Iniciando busca de contatos por fornecedor", map[string]interface{}{
		"supplier_id": supplierID,
	})

	if supplierID <= 0 {
		s.logger.Warn(ctx, "[contactService] - ID de fornecedor inválido", map[string]interface{}{
			"supplier_id": supplierID,
		})
		return nil, ErrSupplierIDInvalid
	}

	contacts, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, "[contactService] - Erro ao listar contatos do fornecedor", map[string]interface{}{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListSupplierContacts, err)
	}

	s.logger.Info(ctx, "[contactService] - Contatos do fornecedor listados com sucesso", map[string]interface{}{
		"supplier_id": supplierID,
		"count":       len(contacts),
	})

	return contacts, nil
}

func (s *contactService) Update(ctx context.Context, contact *models.Contact) error {
	s.logger.Info(ctx, "[contactService] - Iniciando atualização de contato", map[string]interface{}{
		"contact_id":   contact.ID,
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

	if contact.ID <= 0 {
		s.logger.Warn(ctx, "[contactService] - ID inválido ao tentar atualizar contato", map[string]interface{}{
			"contact_id": contact.ID,
		})
		return ErrInvalidID
	}

	if contact.ContactName == "" {
		s.logger.Warn(ctx, "[contactService] - Nome do contato não fornecido na atualização", map[string]interface{}{
			"contact_id": contact.ID,
		})
		return ErrContactNameRequired
	}

	_, err := s.contactRepo.GetByID(ctx, contact.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Warn(ctx, "[contactService] - Contato não encontrado para atualização", map[string]interface{}{
				"contact_id": contact.ID,
			})
			return ErrContactNotFound
		}
		s.logger.Error(ctx, err, "[contactService] - Erro ao buscar contato antes da atualização", map[string]interface{}{
			"contact_id": contact.ID,
		})
		return fmt.Errorf("%w: %v", ErrCheckBeforeUpdate, err)
	}

	err = s.contactRepo.Update(ctx, contact)
	if err != nil {
		s.logger.Error(ctx, err, "[contactService] - Falha ao atualizar contato", map[string]interface{}{
			"contact_id": contact.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	s.logger.Info(ctx, "[contactService] - Contato atualizado com sucesso", map[string]interface{}{
		"contact_id": contact.ID,
	})
	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	s.logger.Info(ctx, "[contactService] - Iniciando exclusão de contato", map[string]interface{}{
		"contact_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, "[contactService] - ID inválido ao tentar deletar contato", map[string]interface{}{
			"contact_id": id,
		})
		return ErrInvalidID
	}

	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Warn(ctx, "[contactService] - Contato não encontrado para deleção", map[string]interface{}{
				"contact_id": id,
			})
			return ErrContactNotFound
		}
		s.logger.Error(ctx, err, "[contactService] - Erro ao verificar existência do contato antes da deleção", map[string]interface{}{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	err = s.contactRepo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, "[contactService] - Erro ao deletar contato", map[string]interface{}{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	s.logger.Info(ctx, "[contactService] - Contato deletado com sucesso", map[string]interface{}{
		"contact_id": id,
	})
	return nil
}
