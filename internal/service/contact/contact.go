package services

import (
	"context"
	"errors"
	"fmt"

	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/contact"
)

type ContactService interface {
	Create(ctx context.Context, contactDTO *dtoContact.ContactDTO) (*dtoContact.ContactDTO, error)
	GetByID(ctx context.Context, id int64) (*dtoContact.ContactDTO, error)
	GetByUserID(ctx context.Context, userID int64) ([]*dtoContact.ContactDTO, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*dtoContact.ContactDTO, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*dtoContact.ContactDTO, error)
	Update(ctx context.Context, contactDTO *dtoContact.ContactDTO) error
	Delete(ctx context.Context, id int64) error
}

type contactService struct {
	contactRepo repo.ContactRepository
	logger      *logger.LogAdapter
}

func NewContactService(contactRepo repo.ContactRepository, logger *logger.LogAdapter) ContactService {
	return &contactService{
		contactRepo: contactRepo,
		logger:      logger,
	}
}

func (s *contactService) Create(ctx context.Context, contactDTO *dtoContact.ContactDTO) (*dtoContact.ContactDTO, error) {
	ref := "[contactService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"contact_name": contactDTO.ContactName,
		"user_id":      utils.Int64OrNil(contactDTO.UserID),
		"client_id":    utils.Int64OrNil(contactDTO.ClientID),
		"supplier_id":  utils.Int64OrNil(contactDTO.SupplierID),
	})

	contactModel := dtoContact.ToContactModel(*contactDTO)

	if err := contactModel.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"error": err.Error(),
		})
		return nil, err
	}

	createdContact, err := s.contactRepo.Create(ctx, contactModel)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"contact_name": contactModel.ContactName,
			"email":        contactModel.Email,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id":   createdContact.ID,
		"contact_name": createdContact.ContactName,
	})

	result := dtoContact.ToContactDTO(createdContact)
	return &result, nil
}

func (s *contactService) GetByID(ctx context.Context, id int64) (*dtoContact.ContactDTO, error) {
	ref := "[contactService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"contact_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": id,
		})
		return nil, errMsg.ErrID
	}

	contactModel, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			s.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"contact_id": id,
			})
			return nil, errMsg.ErrNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"contact_id": id,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	contactDTO := dtoContact.ToContactDTO(contactModel)

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"contact_id": contactDTO.ID,
	})

	return &contactDTO, nil
}

func (s *contactService) GetByUserID(ctx context.Context, userID int64) ([]*dtoContact.ContactDTO, error) {
	ref := "[contactService - GetByUserID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	if userID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return nil, errMsg.ErrID
	}

	contactModels, err := s.contactRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, err
	}

	contactDTOs := make([]*dtoContact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		dto := dtoContact.ToContactDTO(c)
		contactDTOs[i] = &dto
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":     userID,
		"total_items": len(contactDTOs),
	})

	return contactDTOs, nil
}

func (s *contactService) GetByClientID(ctx context.Context, clientID int64) ([]*dtoContact.ContactDTO, error) {
	ref := "[contactService - GetByClientID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
	})

	if clientID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"client_id": clientID,
		})
		return nil, errMsg.ErrID
	}

	contactModels, err := s.contactRepo.GetByClientID(ctx, clientID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		return nil, err
	}

	contactDTOs := make([]*dtoContact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		dto := dtoContact.ToContactDTO(c)
		contactDTOs[i] = &dto
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":   clientID,
		"total_items": len(contactDTOs),
	})

	return contactDTOs, nil
}

func (s *contactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*dtoContact.ContactDTO, error) {
	ref := "[contactService - GetBySupplierID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	if supplierID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, errMsg.ErrID
	}

	contactModels, err := s.contactRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, err
	}

	contactDTOs := make([]*dtoContact.ContactDTO, len(contactModels))
	for i, c := range contactModels {
		dto := dtoContact.ToContactDTO(c)
		contactDTOs[i] = &dto
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplierID,
		"total_items": len(contactDTOs),
	})

	return contactDTOs, nil
}

func (s *contactService) Update(ctx context.Context, contactDTO *dtoContact.ContactDTO) error {
	ref := "[contactService - Update] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"contact_id":  utils.Int64OrNil(contactDTO.ID),
		"user_id":     utils.Int64OrNil(contactDTO.UserID),
		"client_id":   utils.Int64OrNil(contactDTO.ClientID),
		"supplier_id": utils.Int64OrNil(contactDTO.SupplierID),
	})

	contactModel := dtoContact.ToContactModel(*contactDTO)

	if err := contactModel.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return err
	}

	if contactModel.ID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"contact_id": contactModel.ID,
		})
		return errMsg.ErrID
	}

	if err := s.contactRepo.Update(ctx, contactModel); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"contact_id": contactModel.ID,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": contactModel.ID,
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
		return errMsg.ErrID
	}

	_, err := s.contactRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"contact_id": id,
			})
			return errMsg.ErrNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if err := s.contactRepo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})
	return nil
}
