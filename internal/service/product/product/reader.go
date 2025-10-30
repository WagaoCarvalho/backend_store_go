package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *product) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	if limit <= 0 {
		return nil, errMsg.ErrInvalidLimit
	}
	if offset < 0 {
		return nil, errMsg.ErrInvalidOffset
	}

	products, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *product) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *product) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("nome inválido")
	}

	products, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *product) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {

	if strings.TrimSpace(manufacturer) == "" {
		return nil, errors.New("fabricante inválido")
	}

	products, err := s.repo.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *product) GetVersionByID(ctx context.Context, pid int64) (int64, error) {

	if pid <= 0 {
		return 0, errMsg.ErrZeroID
	}

	version, err := s.repo.GetVersionByID(ctx, pid)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return 0, errMsg.ErrNotFound
		}

		return 0, fmt.Errorf("%w: %v", errMsg.ErrVersionConflict, err)
	}

	return version, nil
}

func (s *product) ProductExists(ctx context.Context, productID int64) (bool, error) {
	if productID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.repo.ProductExists(ctx, productID)
	if err != nil {
		return false, err
	}

	return exists, nil
}
