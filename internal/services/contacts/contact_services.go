package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contact"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
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
	ref := "[contactService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

	if err := contact.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	createdContact, err := s.contactRepo.Create(ctx, contact)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"contact_name": contact.ContactName,
			"email":        contact.Email,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateContact, err)
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id":   createdContact.ID,
		"contact_name": createdContact.ContactName,
	})

	return createdContact, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	ref := "[contactService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"contact_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": id,
		})
		return nil, ErrInvalidID
	}

	contact, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"contact_id": id,
			})
			return nil, ErrContactNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"contact_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"contact_id": id,
	})
	return contact, nil
}

func (s *contactService) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	ref := "[contactService - GetByUserID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	if userID <= 0 {
		s.logger.Error(ctx, ErrUserIDInvalid, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return nil, ErrUserIDInvalid
	}

	contacts, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListUserContacts, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":     userID,
		"total_items": len(contacts),
	})

	return contacts, nil
}

func (s *contactService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	ref := "[contactService - GetByClientID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
	})

	if clientID <= 0 {
		s.logger.Error(ctx, ErrClientIDInvalid, ref+logger.LogValidateError, map[string]any{
			"client_id": clientID,
		})
		return nil, ErrClientIDInvalid
	}

	contacts, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListClientContacts, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":   clientID,
		"total_items": len(contacts),
	})

	return contacts, nil
}

func (s *contactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	ref := "[contactService - GetBySupplierID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	if supplierID <= 0 {
		s.logger.Error(ctx, ErrSupplierIDInvalid, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, ErrSupplierIDInvalid
	}

	contacts, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrListSupplierContacts, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierID,
		"total_items": len(contacts),
	})

	return contacts, nil
}

func (s *contactService) Update(ctx context.Context, contact *models.Contact) error {
	ref := "[contactService - Update] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"contact_id":   contact.ID,
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

	if contact.ID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": contact.ID,
		})
		return ErrInvalidID
	}

	if err := contact.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": contact.ID,
			"erro":       err.Error(),
		})
		return err
	}

	_, err := s.contactRepo.GetByID(ctx, contact.ID)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"contact_id": contact.ID,
			})
			return ErrContactNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"contact_id": contact.ID,
		})
		return fmt.Errorf("%w: %v", ErrCheckBeforeUpdate, err)
	}

	if err := s.contactRepo.Update(ctx, contact); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"contact_id": contact.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": contact.ID,
	})

	return nil
}

func (s *contactService) Delete(ctx context.Context, id int64) error {
	ref := "[contactService - Delete] - "
	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"contact_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": id,
		})
		return ErrInvalidID
	}

	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrContactNotFound) {
			s.logger.Warn(ctx, ref+"Contato não encontrado para exclusão", map[string]any{
				"contact_id": id,
			})
			return ErrContactNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrCheckContact, err)
	}

	if err := s.contactRepo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})
	return nil
}
