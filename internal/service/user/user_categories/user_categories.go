package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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
	repo repo.UserCategoryRepository
}

func NewUserCategoryService(repo repo.UserCategoryRepository) UserCategoryService {
	return &userCategoryService{
		repo: repo,
	}
}

func (s *userCategoryService) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	if strings.TrimSpace(category.Name) == "" {
		return nil, err_msg.ErrInvalidData
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
	}

	return createdCategory, nil
}

func (s *userCategoryService) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return categories, nil
}

func (s *userCategoryService) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	if id <= 0 {
		return nil, err_msg.ErrID
	}

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return nil, err_msg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return category, nil
}

func (s *userCategoryService) Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	if category.ID == 0 {
		return nil, err_msg.ErrID
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	if _, err := s.repo.GetByID(ctx, int64(category.ID)); err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return nil, err_msg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
	}

	return category, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return err_msg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
