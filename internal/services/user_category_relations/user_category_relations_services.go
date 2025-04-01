package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/user_category_relations"
)

var (
	ErrInvalidUserID     = errors.New("ID do usuário inválido")
	ErrInvalidCategoryID = errors.New("ID da categoria inválido")
)

type UserCategoryRelationServices interface {
	CreateRelation(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelation, error)
	GetUserRelations(ctx context.Context, userID int64) ([]models.UserCategoryRelation, error)
	GetCategoryRelations(ctx context.Context, categoryID int64) ([]models.UserCategoryRelation, error)
	DeleteRelation(ctx context.Context, userID, categoryID int64) error
	DeleteAllUserRelations(ctx context.Context, userID int64) error
	UserHasCategory(ctx context.Context, userID, categoryID int64) (bool, error)
}

type userCategoryRelationServices struct {
	relationRepo repositories.UserCategoryRelationRepositories
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepositories) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
	}
}

func (s *userCategoryRelationServices) CreateRelation(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelation, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	relation := models.UserCategoryRelation{
		UserID:     userID,
		CategoryID: categoryID,
	}

	createdRelation, err := s.relationRepo.CreateRelation(ctx, relation)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationExists) {
			// Se a relação já existe, retornamos ela
			relations, err := s.relationRepo.GetRelationsByUserID(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("erro ao verificar relação existente: %w", err)
			}
			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return &rel, nil
				}
			}
			return nil, repositories.ErrRelationExists
		}
		return nil, fmt.Errorf("erro ao criar relação: %w", err)
	}

	return &createdRelation, nil
}

func (s *userCategoryRelationServices) GetUserRelations(ctx context.Context, userID int64) ([]models.UserCategoryRelation, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	relations, err := s.relationRepo.GetRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações do usuário: %w", err)
	}

	return relations, nil
}

func (s *userCategoryRelationServices) GetCategoryRelations(ctx context.Context, categoryID int64) ([]models.UserCategoryRelation, error) {
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	relations, err := s.relationRepo.GetRelationsByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar relações da categoria: %w", err)
	}

	return relations, nil
}

func (s *userCategoryRelationServices) DeleteRelation(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if categoryID <= 0 {
		return ErrInvalidCategoryID
	}

	err := s.relationRepo.DeleteRelation(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationNotFound) {
			return err
		}
		return fmt.Errorf("erro ao deletar relação: %w", err)
	}

	return nil
}

func (s *userCategoryRelationServices) DeleteAllUserRelations(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}

	err := s.relationRepo.DeleteAllUserRelations(ctx, userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar todas as relações do usuário: %w", err)
	}

	return nil
}

func (s *userCategoryRelationServices) UserHasCategory(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return false, ErrInvalidCategoryID
	}

	relations, err := s.GetUserRelations(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar relações do usuário: %w", err)
	}

	for _, rel := range relations {
		if rel.CategoryID == categoryID {
			return true, nil
		}
	}

	return false, nil
}
