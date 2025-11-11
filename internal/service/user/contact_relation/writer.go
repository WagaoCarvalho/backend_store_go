package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userContactRelationService) Create(ctx context.Context, relation *models.UserContactRelation) (*models.UserContactRelation, error) {
	if relation == nil {
		return nil, errMsg.ErrNilModel
	}

	if relation.UserID <= 0 || relation.ContactID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	createdRelation, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrRelationExists):
			relations, getErr := s.relationRepo.GetAllRelationsByUserID(ctx, relation.UserID)
			if getErr != nil {
				return nil, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, getErr)
			}

			for _, rel := range relations {
				if rel.ContactID == relation.ContactID {
					return rel, nil
				}
			}

			return nil, errMsg.ErrRelationExists

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			return nil, errMsg.ErrDBInvalidForeignKey

		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return createdRelation, nil
}

func (s *userContactRelationService) Delete(ctx context.Context, userID, contactID int64) error {
	if userID <= 0 {
		return errMsg.ErrZeroID
	}
	if contactID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, userID, contactID)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *userContactRelationService) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
