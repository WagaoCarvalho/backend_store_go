package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_contact_relations"
)

type UserContactRelationServices interface {
	Create(ctx context.Context, userID, contactID int64) (*models.UserContactRelations, bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelations, error)
	HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error)
	Delete(ctx context.Context, userID, contactID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userContactRelationServices struct {
	relationRepo repo.UserContactRelationRepository
}

func NewUserContactRelationServices(repo repo.UserContactRelationRepository) UserContactRelationServices {
	return &userContactRelationServices{
		relationRepo: repo,
	}
}

func (s *userContactRelationServices) Create(ctx context.Context, userID, contactID int64) (*models.UserContactRelations, bool, error) {
	if userID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}
	if contactID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}

	relation := models.UserContactRelations{
		UserID:    userID,
		ContactID: contactID,
	}

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		switch {
		case errors.Is(err, err_msg.ErrRelationExists):
			relations, getErr := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if getErr != nil {
				return nil, false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, getErr)
			}

			for _, rel := range relations {
				if rel.ContactID == contactID {
					return rel, false, nil
				}
			}

			return nil, false, err_msg.ErrRelationExists

		case errors.Is(err, err_msg.ErrDBInvalidForeignKey):
			return nil, false, err_msg.ErrDBInvalidForeignKey

		default:
			return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	return createdRelation, true, nil
}

func (s *userContactRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelations, error) {
	if userID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *userContactRelationServices) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	if userID <= 0 {
		return false, err_msg.ErrZeroID
	}
	if contactID <= 0 {
		return false, err_msg.ErrZeroID
	}

	exists, err := s.relationRepo.HasUserContactRelation(ctx, userID, contactID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	return exists, nil
}

func (s *userContactRelationServices) Delete(ctx context.Context, userID, contactID int64) error {
	if userID <= 0 {
		return err_msg.ErrZeroID
	}
	if contactID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, userID, contactID)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *userContactRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
