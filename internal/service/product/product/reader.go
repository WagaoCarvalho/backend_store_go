package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *productService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
