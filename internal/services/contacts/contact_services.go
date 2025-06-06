package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

// Erros personalizados
var (
	ErrContactNameRequired        = errors.New("nome do contato é obrigatório")
	ErrContactAssociationRequired = errors.New("o contato deve estar associado a um usuário, cliente ou fornecedor")
	ErrInvalidEmail               = errors.New("email inválido")
	ErrInvalidID                  = errors.New("ID inválido")
	ErrUserIDInvalid              = errors.New("ID de usuário inválido")
	ErrClientIDInvalid            = errors.New("ID de cliente inválido")
	ErrSupplierIDInvalid          = errors.New("ID de fornecedor inválido")
	ErrContactNotFound            = errors.New("contato não encontrado")
	ErrCreateContact              = errors.New("erro ao criar contato")
	ErrListUserContacts           = errors.New("erro ao listar contatos do usuário")
	ErrListClientContacts         = errors.New("erro ao listar contatos do cliente")
	ErrListSupplierContacts       = errors.New("erro ao listar contatos do fornecedor")
	ErrUpdateContact              = errors.New("erro ao atualizar contato")
	ErrDeleteContact              = errors.New("erro ao deletar contato")
	ErrCheckContact               = errors.New("erro ao verificar contato")
	ErrVersionRequired            = errors.New("versão do contato é obrigatória")
	ErrVersionMismatch            = errors.New("conflito de versão ao atualizar contato")
	ErrUpdateFailed               = errors.New("erro ao atualizar contato")
	ErrCheckBeforeUpdate          = errors.New("erro ao verificar existência do contato antes da atualização")
)

type ContactService interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetByClient(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}

type contactService struct {
	contactRepo repositories.ContactRepository
}

func NewContactService(contactRepo repositories.ContactRepository) ContactService {
	return &contactService{
		contactRepo: contactRepo,
	}
}

func (s *contactService) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	if contact.ContactName == "" {
		return &models.Contact{}, ErrContactNameRequired
	}

	if contact.UserID == nil && contact.ClientID == nil && contact.SupplierID == nil {
		return &models.Contact{}, ErrContactAssociationRequired
	}

	if contact.Email != "" && !utils_validators.IsValidEmail(contact.Email) {
		return &models.Contact{}, ErrInvalidEmail
	}

	createdContact, err := s.contactRepo.Create(ctx, contact)
	if err != nil {
		return &models.Contact{}, fmt.Errorf("%w: %v", ErrCreateContact, err)
	}

	return createdContact, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return nil, ErrContactNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	return contact, nil
}

func (s *contactService) GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error) {
	if userID <= 0 {
		return nil, ErrUserIDInvalid
	}

	contacts, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrListUserContacts, err)
	}

	return contacts, nil
}

func (s *contactService) GetByClient(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	if clientID <= 0 {
		return nil, ErrClientIDInvalid
	}

	contacts, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrListClientContacts, err)
	}

	return contacts, nil
}

func (s *contactService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	if supplierID <= 0 {
		return nil, ErrSupplierIDInvalid
	}

	contacts, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrListSupplierContacts, err)
	}

	return contacts, nil
}

func (s *contactService) Update(ctx context.Context, c *models.Contact) error {
	if c.ID == 0 || c.ID <= 0 {
		return ErrInvalidID
	}

	if c.ContactName == "" {
		return ErrContactNameRequired
	}

	if c.Version <= 0 {
		return ErrVersionRequired
	}

	_, err := s.contactRepo.GetByID(ctx, c.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return ErrContactNotFound
		}
		return fmt.Errorf("%w: %v", ErrCheckBeforeUpdate, err)
	}

	err = s.contactRepo.Update(ctx, c)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			return ErrVersionMismatch
		}
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidID
	}

	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return ErrContactNotFound
		}
		return fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	err = s.contactRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	return nil
}
