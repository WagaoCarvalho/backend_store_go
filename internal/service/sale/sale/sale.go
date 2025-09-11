package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/sale"
)

type SaleService interface {
	Create(ctx context.Context, sale *models.Sale) (*models.Sale, error)
	GetByID(ctx context.Context, id int64) (*models.Sale, error)
	GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	Update(ctx context.Context, sale *models.Sale) error
	Delete(ctx context.Context, id int64) error
}

type saleService struct {
	repo repo.SaleRepository
}

func NewSaleService(repo repo.SaleRepository) SaleService {
	return &saleService{
		repo: repo,
	}
}

func (s *saleService) Create(ctx context.Context, sale *models.Sale) (*models.Sale, error) {
	if err := sale.Validate(); err != nil {
		return nil, err
	}

	createdSale, err := s.repo.Create(ctx, sale)
	if err != nil {
		return nil, err
	}

	return createdSale, nil
}

func (s *saleService) GetByID(ctx context.Context, id int64) (*models.Sale, error) {
	if id <= 0 {
		return nil, err_msg.ErrID
	}

	saleModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, err_msg.ErrNotFound) {
			return nil, err_msg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", err_msg.ErrGet, err)
	}

	return saleModel, nil
}

func (s *saleService) GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if clientID <= 0 {
		return nil, err_msg.ErrID
	}

	sales, err := s.repo.GetByClientID(ctx, clientID, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if userID <= 0 {
		return nil, err_msg.ErrID
	}

	sales, err := s.repo.GetByUserID(ctx, userID, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if status == "" {
		return nil, err_msg.ErrInvalidData
	}

	sales, err := s.repo.GetByStatus(ctx, status, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error) {
	if start.IsZero() || end.IsZero() {
		return nil, err_msg.ErrInvalidData
	}

	sales, err := s.repo.GetByDateRange(ctx, start, end, limit, offset, orderBy, orderDir)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) Update(ctx context.Context, sale *models.Sale) error {
	if err := sale.Validate(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, sale); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrUpdate, err)
	}

	return nil
}

func (s *saleService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return err_msg.ErrID
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("%w: %v", err_msg.ErrDelete, err)
	}

	return nil
}
