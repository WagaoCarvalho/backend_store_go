package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/logger"
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
	logger       *logger.LoggerAdapter
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepository, logger *logger.LoggerAdapter) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
		logger:       logger,
	}
}

func (s *userCategoryRelationServices) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error) {
	ref := "[userCategoryRelationServices - Create] - "

	if userID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return nil, false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return nil, false, ErrInvalidCategoryID
	}

	relation := models.UserCategoryRelations{
		UserID:     userID,
		CategoryID: categoryID,
	}

	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrRelationExists):
			s.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})

			relations, getErr := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if getErr != nil {
				s.logger.Error(ctx, getErr, ref+logger.LogCheckError, map[string]any{
					"user_id": userID,
				})
				return nil, false, fmt.Errorf("%w: %v", ErrCheckExistingRelation, getErr)
			}

			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return rel, false, nil
				}
			}

			return nil, false, repositories.ErrRelationExists

		case errors.Is(err, repositories.ErrInvalidForeignKey):
			s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, repositories.ErrInvalidForeignKey

		default:
			s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, fmt.Errorf("%w: %v", ErrCreateRelation, err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return createdRelation, true, nil
}

func (s *userCategoryRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	ref := "[userCategoryRelationServices - GetAllRelationsByUserID] - "

	if userID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return nil, ErrInvalidUserID
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUserRelations, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":       userID,
		"relations_len": len(relationsPtr),
	})

	return relationsPtr, nil
}

func (s *userCategoryRelationServices) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	ref := "[userCategoryRelationServices - HasUserCategoryRelation] - "

	if userID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return false, ErrInvalidCategoryID
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", ErrCheckRelationExists, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
		"exists":      exists,
	})

	return exists, nil
}

func (s *userCategoryRelationServices) Delete(ctx context.Context, userID, categoryID int64) error {
	ref := "[userCategoryRelationServices - Delete] - "

	if userID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return ErrInvalidCategoryID
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return err
		}
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return nil
}

func (s *userCategoryRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	ref := "[userCategoryRelationServices - DeleteAll] - "

	if userID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"user_id": userID,
		})
		return ErrInvalidUserID
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id": userID,
	})

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": userID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": userID,
	})

	return nil
}
