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
	Create(ctx context.Context, contact models.Contact) (models.Contact, error)
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

func (s *contactService) Create(ctx context.Context, c models.Contact) (models.Contact, error) {
	// Validações básicas
	if c.ContactName == "" {
		return models.Contact{}, fmt.Errorf("nome do contato é obrigatório")
	}

	if c.UserID == nil && c.ClientID == nil && c.SupplierID == nil {
		return models.Contact{}, fmt.Errorf("o contato deve estar associado a um usuário, cliente ou fornecedor")
	}

	if c.Email != "" && !utils.IsValidEmail(c.Email) {
		return models.Contact{}, fmt.Errorf("email inválido")
	}

	// Chama o repositório — adaptando para trabalhar com struct por valor
	createdContact, err := s.contactRepo.Create(ctx, c)
	if err != nil {
		return models.Contact{}, fmt.Errorf("erro ao criar contato: %w", err)
	}

	return createdContact, nil
}
func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID inválido")
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return nil, fmt.Errorf("contato não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar contato: %w", err)
	}

	return contact, nil
}

func (s *contactService) GetByUser(ctx context.Context, userID int64) ([]*models.Contact, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("ID de usuário inválido")
	}

	contacts, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do usuário: %w", err)
	}

	return contacts, nil
}

func (s *contactService) GetByClient(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	if clientID <= 0 {
		return nil, fmt.Errorf("ID de cliente inválido")
	}

	contacts, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do cliente: %w", err)
	}

	return contacts, nil
}

func (s *contactService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	if supplierID <= 0 {
		return nil, fmt.Errorf("ID de fornecedor inválido")
	}

	contacts, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar contatos do fornecedor: %w", err)
	}

	return contacts, nil
}

func (s *contactService) Update(ctx context.Context, c *models.Contact) error {
	// Validações básicas
	if c.ID == nil || *c.ID <= 0 {
		return fmt.Errorf("ID inválido")
	}

	if c.ContactName == "" {
		return fmt.Errorf("nome do contato é obrigatório")
	}

	// Verifica se o contato existe antes de atualizar
	_, err := s.contactRepo.GetByID(ctx, *c.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return fmt.Errorf("contato não encontrado")
		}
		return fmt.Errorf("erro ao verificar contato: %w", err)
	}

	err = s.contactRepo.Update(ctx, c)
	if err != nil {
		return fmt.Errorf("erro ao atualizar contato: %w", err)
	}

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("ID inválido")
	}

	// Verifica se o contato existe antes de deletar
	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			return fmt.Errorf("contato não encontrado")
		}
		return fmt.Errorf("erro ao verificar contato: %w", err)
	}

	err = s.contactRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar contato: %w", err)
	}

	return nil
}
