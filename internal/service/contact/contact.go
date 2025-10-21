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
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdContact, err := s.contactRepo.Create(ctx, contact)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return nil, errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrDuplicate):
			return nil, errMsg.ErrDuplicate
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return createdContact, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
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

func (s *contactService) Update(ctx context.Context, contact *models.Contact) error {

	if contact.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if err := contact.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrDuplicate):
			return errMsg.ErrDuplicate
		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
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
