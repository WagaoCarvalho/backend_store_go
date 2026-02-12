package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productCategoryRelationService) Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error) {
	if relation == nil {
		return nil, errMsg.ErrNilModel
	}
	if relation.ProductID <= 0 || relation.CategoryID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	createdRelation, err := s.repo.Create(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrRelationExists):
			return nil, errMsg.ErrRelationExists

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			return nil, errMsg.ErrDBInvalidForeignKey

		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return createdRelation, nil
}

func (s *productCategoryRelationService) Delete(ctx context.Context, productID, categoryID int64) error {
	if productID <= 0 {
		return errMsg.ErrZeroID
	}
	if categoryID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.Delete(ctx, productID, categoryID)
	if err != nil {
		// Propaga erros específicos sem encapsular
		if errors.Is(err, errMsg.ErrNotFound) {
			return err
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *productCategoryRelationService) DeleteAll(ctx context.Context, productID int64) error {
	if productID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DeleteAll(ctx, productID)
	if err != nil {
		// DeleteAll não retorna ErrNotFound (sucesso com 0 linhas)
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
