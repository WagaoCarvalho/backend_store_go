package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
)

type UserCategoryService interface {
	GetCategories(ctx context.Context) ([]models.UserCategory, error)
	GetCategoryById(ctx context.Context, id int64) (models.UserCategory, error)
	CreateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error)
	UpdateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error)
	DeleteCategoryById(ctx context.Context, id int64) error
}

type userCategoryService struct {
	repo repositories.UserCategoryRepository
}

func NewUserCategoryService(repo repositories.UserCategoryRepository) UserCategoryService {
	return &userCategoryService{repo: repo}
}

func (s *userCategoryService) GetCategories(ctx context.Context) ([]models.UserCategory, error) {
	categories, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	return categories, nil
}

func (s *userCategoryService) GetCategoryById(ctx context.Context, id int64) (models.UserCategory, error) {
	category, err := s.repo.GetCategoryById(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			return models.UserCategory{}, err // Mantém a mensagem original do erro
		}
		return models.UserCategory{}, fmt.Errorf("erro ao buscar categoria: %w", err) // Adiciona apenas para outros tipos de erro
	}
	return category, nil
}

func (s *userCategoryService) CreateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	createdCategory, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return models.UserCategory{}, fmt.Errorf("erro ao criar categoria: %w", err)
	}
	return createdCategory, nil
}

func (s *userCategoryService) UpdateCategory(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	updatedCategory, err := s.repo.UpdateCategory(ctx, category)
	if err != nil {
		return models.UserCategory{}, fmt.Errorf("erro ao atualizar categoria: %w", err)
	}
	return updatedCategory, nil
}

func (s *userCategoryService) DeleteCategoryById(ctx context.Context, id int64) error {
	if err := s.repo.DeleteCategoryById(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}
	return nil
}
