package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/contact_relation"
)

type UserContactRelation interface {
	Create(ctx context.Context, userID, contactID int64) (*models.UserContactRelation, bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error)
	HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error)
	Delete(ctx context.Context, userID, contactID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userContactRelation struct {
	relationRepo repo.UserContactRelation
}

func NewUserContactRelation(repo repo.UserContactRelation) UserContactRelation {
	return &userContactRelation{
		relationRepo: repo,
	}
}

func (s *userContactRelation) Create(ctx context.Context, userID, contactID int64) (*models.UserContactRelation, bool, error) {
	if userID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}
	if contactID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}

	relation := models.UserContactRelation{
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

func (s *userContactRelation) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelation, error) {
	if userID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *userContactRelation) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
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

func (s *userContactRelation) Delete(ctx context.Context, userID, contactID int64) error {
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

func (s *userContactRelation) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
