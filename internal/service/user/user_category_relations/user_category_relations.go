package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
	logger       *logger.LoggerAdapter
}

func NewUserCategoryRelationServices(repo repo.UserCategoryRelationRepository, logger *logger.LoggerAdapter) UserCategoryRelationServices {
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
		return nil, false, err_msg.ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return nil, false, err_msg.ErrInvalidCategoryID
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
		case errors.Is(err, err_msg.ErrRelationExists):
			s.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})

			relations, getErr := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if getErr != nil {
				s.logger.Error(ctx, getErr, ref+logger.LogCheckError, map[string]any{
					"user_id": userID,
				})
				return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCheckExistingRelation, getErr)
			}

			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return rel, false, nil
				}
			}

			return nil, false, err_msg.ErrRelationExists

		case errors.Is(err, err_msg.ErrInvalidForeignKey):
			s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, err_msg.ErrInvalidForeignKey

		default:
			s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, fmt.Errorf("%w: %v", err_msg.ErrCreateRelation, err)
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
		return nil, err_msg.ErrInvalidUserID
	}

	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", err_msg.ErrFetchUserRelations, err)
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
		return false, err_msg.ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return false, err_msg.ErrInvalidCategoryID
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
		return false, fmt.Errorf("%w: %v", err_msg.ErrCheckRelationExists, err)
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
		return err_msg.ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": categoryID,
		})
		return err_msg.ErrInvalidCategoryID
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, err_msg.ErrRelationNotFound) {
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
		return fmt.Errorf("%w: %v", err_msg.ErrDeleteRelation, err)
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
		return err_msg.ErrInvalidUserID
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id": userID,
	})

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": userID,
		})
		return fmt.Errorf("%w: %v", err_msg.ErrDeleteAllUserRelations, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": userID,
	})

	return nil
}
