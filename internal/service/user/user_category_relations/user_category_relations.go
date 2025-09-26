package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_category_relations"
)

type UserCategoryRelationServices interface {
	Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelationServices struct {
	relationRepo repo.UserCategoryRelationRepository
}

func NewUserCategoryRelationServices(repo repo.UserCategoryRelationRepository) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
	}
}

func (s *userCategoryRelationServices) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error) {
	if userID <= 0 {
		return nil, false, err_msg.ErrIDZero
	}
	if categoryID <= 0 {
		return nil, false, err_msg.ErrIDZero
	}

	relation := models.UserCategoryRelations{
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

		case errors.Is(err, err_msg.ErrInvalidForeignKey):
			return nil, false, err_msg.ErrInvalidForeignKey

		default:
			return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	return createdRelation, true, nil
}

func (s *userCategoryRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	if userID <= 0 {
		return nil, err_msg.ErrIDZero
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *userCategoryRelationServices) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, err_msg.ErrIDZero
	}
	if categoryID <= 0 {
		return false, err_msg.ErrIDZero
	}

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	return exists, nil
}

func (s *userCategoryRelationServices) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return err_msg.ErrIDZero
	}
	if categoryID <= 0 {
		return err_msg.ErrIDZero
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

func (s *userCategoryRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return err_msg.ErrIDZero
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
