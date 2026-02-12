package services

import (
	"context"
	"errors"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productCategoryRelationService) GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error) {
	if productID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	relations, err := s.repo.GetAllRelationsByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, err
		}

		return nil, err
	}

	if relations == nil {
		return []*models.ProductCategoryRelation{}, nil
	}

	return relations, nil
}
