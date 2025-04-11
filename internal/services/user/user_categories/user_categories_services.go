package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_categories"
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

func (s *userCategoryService) GetAll(ctx context.Context) ([]models.UserCategory, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar categorias: %w", err)
	}
	return categories, nil
}

func (s *userCategoryService) GetById(ctx context.Context, id int64) (models.UserCategory, error) {
	category, err := s.repo.GetById(ctx, id)
	if err != nil {
		if err.Error() == "categoria não encontrada" {
			return models.UserCategory{}, err // Mantém a mensagem original do erro
		}
		return models.UserCategory{}, fmt.Errorf("erro ao buscar categoria: %w", err) // Adiciona apenas para outros tipos de erro
	}
	return category, nil
}

func (s *userCategoryService) Create(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		return models.UserCategory{}, fmt.Errorf("erro ao criar categoria: %w", err)
	}
	return createdCategory, nil
}

func (s *userCategoryService) Update(ctx context.Context, category models.UserCategory) (models.UserCategory, error) {
	updatedCategory, err := s.repo.Update(ctx, category)
	if err != nil {
		return models.UserCategory{}, fmt.Errorf("erro ao atualizar categoria: %w", err)
	}
	return updatedCategory, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("erro ao deletar categoria: %w", err)
	}
	return nil
}
