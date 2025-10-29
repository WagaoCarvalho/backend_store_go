package services

import (
	"context"
	"errors"
	"strings"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

// --- Validação comum de paginação e ordenação ---
func validatePagination(limit, offset int) error {
	if limit <= 0 {
		return errMsg.ErrInvalidLimit
	}
	if offset < 0 {
		return errMsg.ErrInvalidOffset
	}
	return nil
}

func validateOrder(orderBy string, allowedFields map[string]bool, orderDir string) (string, error) {
	if !allowedFields[orderBy] {
		return "", errMsg.ErrInvalidOrderField
	}

	orderDir = strings.ToLower(orderDir)
	if orderDir != "asc" && orderDir != "desc" {
		return "", errMsg.ErrInvalidOrderDirection
	}

	return orderDir, nil
}

// --- Sale Reader Service ---
func (s *sale) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
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

func (s *sale) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if clientID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validateOrder(orderBy, map[string]bool{
		"sale_date":    true,
		"total_amount": true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)
}

func (s *sale) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if userID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if err := validatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validateOrder(orderBy, map[string]bool{
		"sale_date":    true,
		"total_amount": true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(ctx, userID, limit, offset, orderBy, orderDir)
}

func (s *sale) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if strings.TrimSpace(status) == "" {
		return nil, errMsg.ErrInvalidData
	}

	if err := validatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validateOrder(orderBy, map[string]bool{
		"id":        true,
		"sale_date": true,
		"total":     true,
		"status":    true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByStatus(ctx, status, limit, offset, orderBy, orderDir)
}

func (s *sale) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errMsg.ErrInvalidData
	}
	if start.After(end) {
		return nil, errMsg.ErrInvalidDateRange
	}

	if err := validatePagination(limit, offset); err != nil {
		return nil, err
	}

	orderDir, err := validateOrder(orderBy, map[string]bool{
		"id":        true,
		"sale_date": true,
		"total":     true,
	}, orderDir)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByDateRange(ctx, start, end, limit, offset, orderBy, orderDir)
}
