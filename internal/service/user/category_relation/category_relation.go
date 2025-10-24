package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/category_relation"
)

type UserCategoryRelation interface {
	Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelation, bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelation struct {
	relationRepo repo.UserCategoryRelation
}

func NewUserCategoryRelation(repo repo.UserCategoryRelation) UserCategoryRelation {
	return &userCategoryRelation{
		relationRepo: repo,
	}
}

func (s *userCategoryRelation) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelation, bool, error) {
	if userID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}
	if categoryID <= 0 {
		return nil, false, err_msg.ErrZeroID
	}

	relation := models.UserCategoryRelation{
		UserID:     userID,
		CategoryID: categoryID,
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
				if rel.CategoryID == categoryID {
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

func (s *userCategoryRelation) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error) {
	if userID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *userCategoryRelation) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, err_msg.ErrZeroID
	}
	if categoryID <= 0 {
		return false, err_msg.ErrZeroID
	}

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	return exists, nil
}

func (s *userCategoryRelation) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return err_msg.ErrZeroID
	}
	if categoryID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *userCategoryRelation) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
