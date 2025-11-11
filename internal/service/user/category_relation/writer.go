package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *userCategoryRelationService) Create(ctx context.Context, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error) {
	if relation == nil {
		return nil, errMsg.ErrNilModel
	}

	if relation.UserID <= 0 || relation.CategoryID <= 0 {
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
				if rel.CategoryID == relation.CategoryID {
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

func (s *userCategoryRelationService) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return errMsg.ErrZeroID
	}
	if categoryID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *userCategoryRelationService) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
