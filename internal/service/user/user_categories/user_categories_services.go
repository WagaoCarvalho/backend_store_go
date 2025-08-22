package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user_categories"
)

type UserCategoryService interface {
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Delete(ctx context.Context, id int64) error
}

type userCategoryService struct {
	repo   repo.UserCategoryRepository
	logger *logger.LoggerAdapter
}

func NewUserCategoryService(repo repo.UserCategoryRepository, logger *logger.LoggerAdapter) UserCategoryService {
	return &userCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userCategoryService) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	ref := "[userCategoryService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name": category.Name,
	})

	if strings.TrimSpace(category.Name) == "" {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"name": category.Name,
		})
		return nil, ErrInvalidCategoryName
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": category.Name,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"category_id": createdCategory.ID,
		"name":        createdCategory.Name,
	})

	return createdCategory, nil
}

func (s *userCategoryService) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	ref := "[userCategoryService - GetAll] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, nil)

	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrFetchCategories, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"count": len(categories),
	})

	return categories, nil
}

func (s *userCategoryService) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	ref := "[userCategoryService - GetByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"category_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"category_id": id,
		})
		return nil, ErrCategoryIDRequired
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"category_id": id,
			})
			return nil, ErrCategoryNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"category_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"category_id": category.ID,
	})

	return category, nil
}

func (s *userCategoryService) Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	ref := "[userCategoryService - Update] - "

	if category.ID == 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"id": category.ID,
		})
		return nil, ErrCategoryIDRequired
	}

	if err := category.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"id":   category.ID,
			"erro": err.Error(),
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"id":   category.ID,
		"name": category.Name,
	})

	if _, err := s.repo.GetByID(ctx, int64(category.ID)); err != nil {
		if errors.Is(err, err_msg.ErrCategoryNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": category.ID,
			})
			return nil, ErrCategoryNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"id": category.ID,
		})
		return nil, fmt.Errorf("%w: %v", ErrCheckBeforeUpdate, err)
	}

	if err := s.repo.Update(ctx, category); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"id": category.ID,
		})
		return nil, fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"id": category.ID,
	})

	return category, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	ref := "[userCategoryService - Delete] - "

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"id": id,
		})
		return ErrCategoryIDRequired
	}

	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"id": id,
	})

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, err_msg.ErrCategoryNotFound) {
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"id": id,
			})
			return ErrCategoryNotFound
		}
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"id": id,
		})
		return fmt.Errorf("%w: %v", ErrFetchCategory, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"id": id,
	})

	return nil
}
