package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users/user_category_relations"
)

var (
	ErrInvalidUserID              = errors.New("ID do usuário inválido")
	ErrInvalidCategoryID          = errors.New("ID da categoria inválido")
	ErrCreateRelation             = errors.New("erro ao criar relação")
	ErrCheckExistingRelation      = errors.New("erro ao verificar relação existente")
	ErrFetchUserRelations         = errors.New("erro ao buscar relações do usuário")
	ErrFetchCategoryRelations     = errors.New("erro ao buscar relações da categoria")
	ErrDeleteRelation             = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations     = errors.New("erro ao deletar todas as relações do usuário")
	ErrCheckUserCategoryRelations = errors.New("erro ao verificar relações do usuário")
)

type UserCategoryRelationServices interface {
	Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, error)
	GetAll(ctx context.Context, userID int64) ([]models.UserCategoryRelations, error)
	GetRelations(ctx context.Context, categoryID int64) ([]models.UserCategoryRelations, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
	GetByCategoryID(ctx context.Context, userID, categoryID int64) (bool, error)
}

type userCategoryRelationServices struct {
	relationRepo repositories.UserCategoryRelationRepositories
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepositories) UserCategoryRelationServices {
	return &userCategoryRelationServices{
		relationRepo: repo,
	}
}

func (s *userCategoryRelationServices) Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	relation := models.UserCategoryRelations{
		UserID:     userID,
		CategoryID: categoryID,
	}

	createdRelation, err := s.relationRepo.Create(ctx, &relation)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationExists) {
			relations, err := s.relationRepo.GetByUserID(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrCheckExistingRelation, err)
			}
			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return &rel, nil
				}
			}
			return nil, repositories.ErrRelationExists
		}
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return createdRelation, nil
}

func (s *userCategoryRelationServices) GetAll(ctx context.Context, userID int64) ([]models.UserCategoryRelations, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	relations, err := s.relationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchUserRelations, err)
	}

	return relations, nil
}

func (s *userCategoryRelationServices) GetRelations(ctx context.Context, categoryID int64) ([]models.UserCategoryRelations, error) {
	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	relations, err := s.relationRepo.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchCategoryRelations, err)
	}

	return relations, nil
}

func (s *userCategoryRelationServices) Delete(ctx context.Context, userID, categoryID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if categoryID <= 0 {
		return ErrInvalidCategoryID
	}

	err := s.relationRepo.Delete(ctx, userID, categoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrRelationNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	return nil
}

func (s *userCategoryRelationServices) DeleteAll(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}

	err := s.relationRepo.DeleteAll(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	return nil
}

func (s *userCategoryRelationServices) GetByCategoryID(ctx context.Context, userID, categoryID int64) (bool, error) {
	if userID <= 0 {
		return false, ErrInvalidUserID
	}
	if categoryID <= 0 {
		return false, ErrInvalidCategoryID
	}

	relations, err := s.GetAll(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", ErrCheckUserCategoryRelations, err)
	}

	for _, rel := range relations {
		if rel.CategoryID == categoryID {
			return true, nil
		}
	}

	return false, nil
}
