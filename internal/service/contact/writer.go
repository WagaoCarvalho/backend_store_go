package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *contactService) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	if contact == nil {
		return nil, errMsg.ErrInvalidData
	}

	if err := contact.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", errMsg.ErrInvalidData, err)
	}

	createdContact, err := s.repo.Create(ctx, contact)
	if err != nil {
		return nil, err
	}

	return createdContact, nil
}

func (s *contactService) Update(ctx context.Context, contact *models.Contact) error {
	if contact.ID <= 0 {
		return errMsg.ErrZeroID
	}

	if err := contact.Validate(); err != nil {
		return fmt.Errorf("%w: %w", errMsg.ErrInvalidData, err)
	}

	if err := s.repo.Update(ctx, contact); err != nil {
		switch {
		case errors.Is(err, errMsg.ErrNotFound):
			return errMsg.ErrNotFound
		case errors.Is(err, errMsg.ErrDuplicate):
			return errMsg.ErrDuplicate
		default:
			return fmt.Errorf("%w: %w", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %w", errMsg.ErrDelete, err)
	}

	return nil
}
