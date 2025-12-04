package services

import (
	"context"
	"errors"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	validate "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

// --- Sale Reader Service ---
func (s *saleService) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return nil, errMsg.ErrNotFound
		}
		return nil, err
	}

	return saleModel, nil
}

func (s *saleService) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if clientID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validate.ValidateOrder(orderBy, map[string]bool{
		"sale_date":    true,
		"total_amount": true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)
}

func (s *saleService) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if userID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validate.ValidateOrder(orderBy, map[string]bool{
		"sale_date":    true,
		"total_amount": true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID, limit, offset, orderBy, orderDir)
}

func (s *saleService) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errMsg.ErrInvalidData
	}
	if start.After(end) {
		return nil, errMsg.ErrInvalidDateRange
	}

	if err := validate.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validate.ValidateOrder(orderBy, map[string]bool{
		"id":        true,
		"sale_date": true,
		"total":     true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByDateRange(ctx, start, end, limit, offset, orderBy, orderDir)
}
