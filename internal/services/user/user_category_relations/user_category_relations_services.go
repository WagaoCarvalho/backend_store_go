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
	ErrGetVersion                 = errors.New("erro ao obter a versão da relação")
)

type UserCategoryRelationServices interface {
	Create(ctx context.Context, userID, categoryID int64) (*models.UserCategoryRelations, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	GetVersionByUserID(ctx context.Context, userID int64) (int, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}
type userCategoryRelationServices struct {
	relationRepo repositories.UserCategoryRelationRepository
}

func NewUserCategoryRelationServices(repo repositories.UserCategoryRelationRepository) UserCategoryRelationServices {
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
			relations, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
			if err != nil {
				return nil, fmt.Errorf("%w: %v", ErrCheckExistingRelation, err)
			}
			for _, rel := range relations {
				if rel.CategoryID == categoryID {
					return rel, nil
				}
			}
			return nil, repositories.ErrRelationExists
		}
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return createdRelation, nil
}

func (s *userCategoryRelationServices) GetVersionByUserID(ctx context.Context, userID int64) (int, error) {
	if userID <= 0 {
		return 0, ErrInvalidUserID
	}

	version, err := s.relationRepo.GetVersionByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrGetVersion, err)
	}

	return version, nil
}

func (s *userCategoryRelationServices) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchUserRelations, err)
	}

	return relationsPtr, nil
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
