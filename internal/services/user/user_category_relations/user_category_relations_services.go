package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
)

type UserCategoryRelationServices interface {
	Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}
type userCategoryRelationServices struct {
	relationRepo repositories.UserCategoryRelationRepository
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepository) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
	}
}

func (s *userCategoryRelationServices) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error) {
	if userID <= 0 {
		return nil, false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return nil, false, ErrInvalidCategoryID
	}

	relation := models.UserCategoryRelations{
		UserID:     userID,
		CategoryID: categoryID,
	}

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationExists) {
			relations, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if err != nil {
				return nil, false, fmt.Errorf("%w: %v", ErrCheckExistingRelation, err)
			}
			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return rel, false, nil
				}
			}
			return nil, false, repositories.ErrRelationExists
		}
		return nil, false, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return createdRelation, true, nil
}

func (s *userCategoryRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchUserRelations, err)
	}

	return relationsPtr, nil
}

func (s *userCategoryRelationServices) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return false, ErrInvalidCategoryID
	}

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrCheckRelationExists, err)
	}

	return exists, nil
}

func (s *userCategoryRelationServices) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if categoryID <= 0 {
		return ErrInvalidCategoryID
	}

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	return nil
}

func (s *userCategoryRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	return nil
}
