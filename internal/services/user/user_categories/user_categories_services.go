package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
)

var (
	ErrCreateCategory   = errors.New("erro ao criar categoria")
	ErrFetchCategories  = errors.New("erro ao buscar categorias")
	ErrFetchCategory    = errors.New("erro ao buscar categoria")
	ErrUpdateCategory   = errors.New("erro ao atualizar categoria")
	ErrDeleteCategory   = errors.New("erro ao deletar categoria")
	ErrCategoryNotFound = errors.New("categoria não encontrada")
)

type UserCategoryService interface {
	GetAll(ctx context.Context) ([]models.UserCategory, error)
	GetById(ctx context.Context, id int64) (models.UserCategory, error)
	Create(ctx context.Context, category models.UserCategory) (models.UserCategory, error)
	Update(ctx context.Context, category models.UserCategory) (models.UserCategory, error)
	Delete(ctx context.Context, id int64) error
}

type userCategoryService struct {
	repo repositories.UserCategoryRepository
}

func NewUserCategoryService(repo repositories.UserCategoryRepository) UserCategoryService {
	return &userCategoryService{repo: repo}
}

func (s *userCategoryService) Create(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return models.UserCategory{}, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}
	return createdCategory, nil
}

func (s *userCategoryService) GetAll(ctx context.Context) ([]models.UserCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchCategories, err)
	}
	return categories, nil
}

func (s *userCategoryService) GetById(ctx context.Context, id int64) (models.UserCategory, error) {
	category, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return models.UserCategory{}, ErrCategoryNotFound
		}
		return models.UserCategory{}, fmt.Errorf("%w: %w", ErrFetchCategory, err)
	}
	return category, nil
}

func (s *userCategoryService) Update(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	updatedCategory, err := s.repo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, repositories.ErrVersionConflict) {
			return models.UserCategory{}, repositories.ErrVersionConflict
		}
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			return models.UserCategory{}, repositories.ErrCategoryNotFound
		}
		return models.UserCategory{}, fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}
	return updatedCategory, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}
	return nil
}
