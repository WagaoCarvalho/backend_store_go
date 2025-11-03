package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product/category_relation"
)

type ProductCategoryRelation interface {
	Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error)
	GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error)
	HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error)
	Delete(ctx context.Context, productID, categoryID int64) error
	DeleteAll(ctx context.Context, productID int64) error
}

type productCategoryRelation struct {
	relationRepo repo.ProductCategoryRelation
}

func NewProductCategoryRelation(repo repo.ProductCategoryRelation) ProductCategoryRelation {
	return &productCategoryRelation{
		relationRepo: repo,
	}
}

func (s *productCategoryRelation) Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	if relation == nil {
		return nil, err_msg.ErrNilModel
	}
	if relation.ProductID <= 0 || relation.CategoryID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	createdRelation, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, err_msg.ErrRelationExists):
			relations, getErr := s.relationRepo.GetAllRelationsByProductID(ctx, relation.ProductID)
			if getErr != nil {
				return nil, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, getErr)
			}

			for _, rel := range relations {
				if rel.CategoryID == relation.CategoryID {
					return rel, nil
				}
			}

			return nil, err_msg.ErrRelationExists

		case errors.Is(err, err_msg.ErrDBInvalidForeignKey):
			return nil, err_msg.ErrDBInvalidForeignKey

		default:
			return nil, fmt.Errorf("%w: %v", err_msg.ErrCreate, err)
		}
	}

	return createdRelation, nil
}

func (s *productCategoryRelation) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	if productID <= 0 {
		return nil, err_msg.ErrZeroID
	}

	relationsPtr, err := s.relationRepo.GetAllRelationsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *productCategoryRelation) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	if productID <= 0 {
		return false, err_msg.ErrZeroID
	}
	if categoryID <= 0 {
		return false, err_msg.ErrZeroID
	}

	exists, err := s.relationRepo.HasProductCategoryRelation(ctx, productID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", err_msg.ErrRelationCheck, err)
	}

	return exists, nil
}

func (s *productCategoryRelation) Delete(ctx context.Context, productID, categoryID int64) error {
	if productID <= 0 {
		return err_msg.ErrZeroID
	}
	if categoryID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, productID, categoryID)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}

func (s *productCategoryRelation) DeleteAll(ctx context.Context, productID int64) error {
	if productID <= 0 {
		return err_msg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, productID)
	if err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
