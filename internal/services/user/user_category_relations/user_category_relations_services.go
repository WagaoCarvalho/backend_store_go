package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
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
	logger       *logger.LoggerAdapter
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepository, logger *logger.LoggerAdapter) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
		logger:       logger,
	}
}

func (s *userCategoryRelationServices) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, bool, error) {
	if userID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de usuário inválido", map[string]interface{}{
			"user_id": userID,
		})
		return nil, false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de categoria inválido", map[string]interface{}{
			"category_id": categoryID,
		})
		return nil, false, ErrInvalidCategoryID
	}

	relation := models.UserCategoryRelations{
		UserID:     userID,
		CategoryID: categoryID,
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Iniciando criação da relação usuário-categoria", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrRelationExists):
			s.logger.Warn(ctx, "[userCategoryRelationServices] - Relação já existente", map[string]interface{}{
				"user_id":     userID,
				"category_id": categoryID,
			})

			// Verifica se a relação de fato existe
			relations, getErr := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if getErr != nil {
				s.logger.Error(ctx, getErr, "[userCategoryRelationServices] - Erro ao verificar relações existentes", map[string]interface{}{
					"user_id": userID,
				})
				return nil, false, fmt.Errorf("%w: %v", ErrCheckExistingRelation, getErr)
			}

			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return rel, false, nil
				}
			}
			// Retorna erro se não encontrar
			return nil, false, repositories.ErrRelationExists

		case errors.Is(err, repositories.ErrInvalidForeignKey):
			s.logger.Warn(ctx, "[userCategoryRelationServices] - Chave estrangeira inválida", map[string]interface{}{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, repositories.ErrInvalidForeignKey

		default:
			s.logger.Error(ctx, err, "[userCategoryRelationServices] - Erro ao criar relação", map[string]interface{}{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return nil, false, fmt.Errorf("%w: %v", ErrCreateRelation, err)
		}
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Relação criada com sucesso", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return createdRelation, true, nil
}

func (s *userCategoryRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	if userID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de usuário inválido para listar relações", map[string]interface{}{
			"user_id": userID,
		})
		return nil, ErrInvalidUserID
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Iniciando listagem das relações do usuário", map[string]interface{}{
		"user_id": userID,
	})

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, "[userCategoryRelationServices] - Erro ao buscar relações do usuário", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUserRelations, err)
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Relações do usuário listadas com sucesso", map[string]interface{}{
		"user_id":       userID,
		"relations_len": len(relationsPtr),
	})

	return relationsPtr, nil
}

func (s *userCategoryRelationServices) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de usuário inválido para verificação de relação", map[string]interface{}{
			"user_id": userID,
		})
		return false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de categoria inválido para verificação de relação", map[string]interface{}{
			"category_id": categoryID,
		})
		return false, ErrInvalidCategoryID
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Verificando existência de relação usuário-categoria", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	exists, err := s.relationRepo.HasUserCategoryRelation(ctx, userID, categoryID)
	if err != nil {
		s.logger.Error(ctx, err, "[userCategoryRelationServices] - Erro ao verificar existência da relação", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", ErrCheckRelationExists, err)
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Verificação concluída", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
		"exists":      exists,
	})

	return exists, nil
}

func (s *userCategoryRelationServices) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de usuário inválido para deleção de relação", map[string]interface{}{
			"user_id": userID,
		})
		return ErrInvalidUserID
	}
	if categoryID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de categoria inválido para deleção de relação", map[string]interface{}{
			"category_id": categoryID,
		})
		return ErrInvalidCategoryID
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Iniciando deleção da relação usuário-categoria", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationNotFound) {
			s.logger.Warn(ctx, "[userCategoryRelationServices] - Relação usuário-categoria não encontrada para deleção", map[string]interface{}{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return err
		}
		s.logger.Error(ctx, err, "[userCategoryRelationServices] - Erro ao deletar relação usuário-categoria", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Relação usuário-categoria deletada com sucesso", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return nil
}

func (s *userCategoryRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		s.logger.Warn(ctx, "[userCategoryRelationServices] - ID de usuário inválido para deleção de todas as relações", map[string]interface{}{
			"user_id": userID,
		})
		return ErrInvalidUserID
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Iniciando deleção de todas as relações do usuário", map[string]interface{}{
		"user_id": userID,
	})

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		s.logger.Error(ctx, err, "[userCategoryRelationServices] - Erro ao deletar todas as relações do usuário", map[string]interface{}{
			"user_id": userID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	s.logger.Info(ctx, "[userCategoryRelationServices] - Todas as relações do usuário deletadas com sucesso", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}
