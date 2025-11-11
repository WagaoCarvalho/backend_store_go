package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productCategoryRelationService) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	if productID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	relationsPtr, err := s.repo.GetAllRelationsByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return relationsPtr, nil
}

func (s *productCategoryRelationService) HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error) {
	if productID <= 0 {
		return false, errMsg.ErrZeroID
	}
	if categoryID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.repo.HasProductCategoryRelation(ctx, productID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return exists, nil
}
