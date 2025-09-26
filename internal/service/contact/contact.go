package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
)

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
	contactRepo repo.ContactRepository
}

func NewContactService(contactRepo repo.ContactRepository) ContactService {
	return &contactService{
		contactRepo: contactRepo,
	}
}

func (s *contactService) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	if err := contact.Validate(); err != nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	createdContact, err := s.contactRepo.Create(ctx, contact)
	if err != nil {
		return nil, err
	}

	return createdContact, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, errMsg.ErrIDZero
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return contact, nil
}

func (s *contactService) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	if userID <= 0 {
		return nil, errMsg.ErrIDZero
	}

	contacts, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (s *contactService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	if clientID <= 0 {
		return nil, errMsg.ErrIDZero
	}

	contacts, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (s *contactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	if supplierID <= 0 {
		return nil, errMsg.ErrIDZero
	}

	contacts, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

func (s *contactService) Update(ctx context.Context, contact *models.Contact) error {
	if contact.ID <= 0 {
		return errMsg.ErrIDZero
	}
	if err := contact.Validate(); err != nil {
		return fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrIDZero
	}

	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if err := s.contactRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
