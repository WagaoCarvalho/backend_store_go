package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
)

type UserCategoryService interface {
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Delete(ctx context.Context, id int64) error
}

type userCategoryService struct {
	repo repositories.UserCategoryRepository
}

func NewUserCategoryService(repo repositories.UserCategoryRepository) UserCategoryService {
	return &userCategoryService{repo: repo}
}

func (s *userCategoryService) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	if strings.TrimSpace(category.Name) == "" {
		return nil, ErrInvalidCategoryName
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	return createdCategory, nil
}

func (s *userCategoryService) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchCategories, err)
	}
	return categories, nil
}

func (s *userCategoryService) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	if id == 0 {
		return nil, ErrCategoryIDRequired
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("%w: %w", ErrFetchCategory, err)
	}

	return category, nil
}

func (s *userCategoryService) Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	if category == nil {
		return nil, ErrInvalidCategory
	}
	if category.ID == 0 {
		return nil, ErrCategoryIDRequired
	}

	err := s.repo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return nil, repositories.ErrCategoryNotFound
		}
		return nil, fmt.Errorf("%w", err)
	}

	return category, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrCategoryIDRequired
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	return nil
}
