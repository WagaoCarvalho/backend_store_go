package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	validate "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func (s *saleItemService) GetByID(ctx context.Context, id int64) (*models.SaleItem, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	return s.repo.GetByID(ctx, id)
}

func (s *saleItemService) GetBySaleID(ctx context.Context, saleID int64, limit, offset int) ([]*models.SaleItem, error) {
	if saleID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	return s.repo.GetBySaleID(ctx, saleID, limit, offset)
}

func (s *saleItemService) GetByProductID(ctx context.Context, productID int64, limit, offset int) ([]*models.SaleItem, error) {
	if productID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	return s.repo.GetByProductID(ctx, productID, limit, offset)
}
