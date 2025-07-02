package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
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
	repo   repositories.UserCategoryRepository
	logger *logger.LoggerAdapter
}

func NewUserCategoryService(repo repositories.UserCategoryRepository, logger *logger.LoggerAdapter) UserCategoryService {
	return &userCategoryService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userCategoryService) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	s.logger.Info(ctx, "[userCategoryService] - Iniciando criação de categoria", map[string]interface{}{
		"name": category.Name,
	})

	if strings.TrimSpace(category.Name) == "" {
		s.logger.Warn(ctx, "[userCategoryService] - Nome da categoria inválido (vazio)", nil)
		return nil, ErrInvalidCategoryName
	}

	createdCategory, err := s.repo.Create(ctx, category)
	if err != nil {
		s.logger.Error(ctx, err, "[userCategoryService] - Erro ao criar categoria", map[string]interface{}{
			"name": category.Name,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	s.logger.Info(ctx, "[userCategoryService] - Categoria criada com sucesso", map[string]interface{}{
		"id":   createdCategory.ID,
		"name": createdCategory.Name,
	})

	return createdCategory, nil
}

func (s *userCategoryService) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	s.logger.Info(ctx, "[userCategoryService] - Iniciando recuperação de todas as categorias", nil)

	categories, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error(ctx, err, "[userCategoryService] - Erro ao recuperar categorias", nil)
		return nil, fmt.Errorf("%w: %v", ErrFetchCategories, err)
	}

	s.logger.Info(ctx, "[userCategoryService] - Categorias recuperadas com sucesso", map[string]interface{}{
		"count": len(categories),
	})

	return categories, nil
}

func (s *userCategoryService) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	if id == 0 {
		s.logger.Warn(ctx, "[userCategoryService] - ID da categoria inválido (zero)", map[string]interface{}{
			"category_id": id,
		})
		return nil, ErrCategoryIDRequired
	}

	s.logger.Info(ctx, "[userCategoryService] - Iniciando recuperação da categoria por ID", map[string]interface{}{
		"category_id": id,
	})

	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			s.logger.Warn(ctx, "[userCategoryService] - Categoria não encontrada", map[string]interface{}{
				"category_id": id,
			})
			return nil, ErrCategoryNotFound
		}
		s.logger.Error(ctx, err, "[userCategoryService] - Erro ao recuperar categoria", map[string]interface{}{
			"category_id": id,
		})
		return nil, fmt.Errorf("%w: %w", ErrFetchCategory, err)
	}

	s.logger.Info(ctx, "[userCategoryService] - Categoria recuperada com sucesso", map[string]interface{}{
		"category_id": id,
	})

	return category, nil
}

func (s *userCategoryService) Update(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	if category == nil {
		s.logger.Warn(ctx, "[userCategoryService] - Categoria inválida (nil) passada para atualização", nil)
		return nil, ErrInvalidCategory
	}
	if category.ID == 0 {
		s.logger.Warn(ctx, "[userCategoryService] - ID da categoria inválido (zero) para atualização", nil)
		return nil, ErrCategoryIDRequired
	}

	s.logger.Info(ctx, "[userCategoryService] - Iniciando atualização da categoria", map[string]interface{}{
		"category_id": category.ID,
		"name":        category.Name,
	})

	err := s.repo.Update(ctx, category)
	if err != nil {
		if errors.Is(err, repositories.ErrCategoryNotFound) {
			s.logger.Warn(ctx, "[userCategoryService] - Categoria não encontrada para atualização", map[string]interface{}{
				"category_id": category.ID,
			})
			return nil, repositories.ErrCategoryNotFound
		}
		s.logger.Error(ctx, err, "[userCategoryService] - Erro ao atualizar categoria", map[string]interface{}{
			"category_id": category.ID,
		})
		return nil, fmt.Errorf("%w", err)
	}

	s.logger.Info(ctx, "[userCategoryService] - Categoria atualizada com sucesso", map[string]interface{}{
		"category_id": category.ID,
	})

	return category, nil
}

func (s *userCategoryService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		s.logger.Warn(ctx, "[userCategoryService] - ID da categoria inválido (zero) para deleção", nil)
		return ErrCategoryIDRequired
	}

	s.logger.Info(ctx, "[userCategoryService] - Iniciando deleção da categoria", map[string]interface{}{
		"category_id": id,
	})

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error(ctx, err, "[userCategoryService] - Erro ao deletar categoria", map[string]interface{}{
			"category_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	s.logger.Info(ctx, "[userCategoryService] - Categoria deletada com sucesso", map[string]interface{}{
		"category_id": id,
	})

	return nil
}
