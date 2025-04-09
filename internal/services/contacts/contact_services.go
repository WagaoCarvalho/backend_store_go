package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	"github.com/WagaoCarvalho/backend_store_go/utils"
)

type ContactService interface {
	CreateContact(ctx context.Context, contact *models.Contact) error
	GetContactByID(ctx context.Context, id int64) (*models.Contact, error)
	GetContactsByUser(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetContactsByClient(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetContactsBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	UpdateContact(ctx context.Context, contact *models.Contact) error
	DeleteContact(ctx context.Context, id int64) error
}

type contactService struct {
	contactRepo repositories.ContactRepository
}

func NewContactService(contactRepo repositories.ContactRepository) ContactService {
	return &contactService{
		contactRepo: contactRepo,
	}
}

func (s *contactService) CreateContact(ctx context.Context, c *models.Contact) error {
	// Validações básicas
	if c.ContactName == "" {
		return fmt.Errorf("nome do contato é obrigatório")
	}

	if c.UserID == nil && c.ClientID == nil && c.SupplierID == nil {
		return fmt.Errorf("o contato deve estar associado a um usuário, cliente ou fornecedor")
	}

	// Verifica se o email é válido (se fornecido)
	if c.Email != "" && !utils.IsValidEmail(c.Email) {
		return fmt.Errorf("email inválido")
	}

	err := s.contactRepo.CreateContact(ctx, c)
	if err != nil {
		return fmt.Errorf("erro ao criar contato: %w", err)
	}

	return nil
}

func (s *contactService) GetContactByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	contact, err := s.contactRepo.GetContactByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return nil, fmt.Errorf("contato não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar contato: %w", err)
	}

	return contact, nil
}

func (s *contactService) GetContactsByUser(ctx context.Context, userID int64) ([]*models.Contact, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID de usuário inválido")
	}

	contacts, err := s.contactRepo.GetContactByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do usuário: %w", err)
	}

	return contacts, nil
}

func (s *contactService) GetContactsByClient(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	if clientID <= 0 {
		return nil, fmt.Errorf("ID de cliente inválido")
	}

	contacts, err := s.contactRepo.GetContactByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do cliente: %w", err)
	}

	return contacts, nil
}

func (s *contactService) GetContactsBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	if supplierID <= 0 {
		return nil, fmt.Errorf("ID de fornecedor inválido")
	}

	contacts, err := s.contactRepo.GetContactBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do fornecedor: %w", err)
	}

	return contacts, nil
}

func (s *contactService) UpdateContact(ctx context.Context, c *models.Contact) error {
	// Validações básicas
	if c.ID <= 0 {
		return fmt.Errorf("ID inválido")
	}

	if c.ContactName == "" {
		return fmt.Errorf("nome do contato é obrigatório")
	}

	// Verifica se o contato existe antes de atualizar
	_, err := s.contactRepo.GetContactByID(ctx, c.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return fmt.Errorf("contato não encontrado")
		}
		return fmt.Errorf("erro ao verificar contato: %w", err)
	}

	err = s.contactRepo.Updatecontac(ctx, c)
	if err != nil {
		return fmt.Errorf("erro ao atualizar contato: %w", err)
	}

	return nil
}

func (s *contactService) DeleteContact(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID inválido")
	}

	// Verifica se o contato existe antes de deletar
	_, err := s.contactRepo.GetContactByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return fmt.Errorf("contato não encontrado")
		}
		return fmt.Errorf("erro ao verificar contato: %w", err)
	}

	err = s.contactRepo.Deletecontact(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar contato: %w", err)
	}

	return nil
}
