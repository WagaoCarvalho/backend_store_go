package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
)

type ProductService interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetByName(ctx context.Context, name string) ([]*models.Product, error)
	GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error)
	GetVersionByID(ctx context.Context, uid int64) (int64, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error

	DisableProduct(ctx context.Context, uid int64) error
	EnableProduct(ctx context.Context, uid int64) error

	UpdateStock(ctx context.Context, id int64, quantity int) error
	IncreaseStock(ctx context.Context, id int64, amount int) error
	DecreaseStock(ctx context.Context, id int64, amount int) error
	GetStock(ctx context.Context, id int64) (int, error)

	EnableDiscount(ctx context.Context, id int64) error
	DisableDiscount(ctx context.Context, id int64) error
	ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error)
}

type productService struct {
	repo repo.ProductRepository
}

func NewProductService(repo repo.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	if product == nil {
		return nil, fmt.Errorf("%w", errMsg.ErrInvalidData)
	}

	if err := product.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return createdProduct, nil
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
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

func (s *productService) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("nome inválido")
	}

	products, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {

	if strings.TrimSpace(manufacturer) == "" {
		return nil, errors.New("fabricante inválido")
	}

	products, err := s.repo.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *productService) GetVersionByID(ctx context.Context, pid int64) (int64, error) {

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

func (s *productService) DisableProduct(ctx context.Context, uid int64) error {

	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DisableProduct(ctx, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (s *productService) EnableProduct(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.EnableProduct(ctx, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}

func (s *productService) Update(ctx context.Context, product *models.Product) error {

	if product.ID <= 0 {
		return errMsg.ErrZeroID
	}
	if err := product.Validate(); err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrInvalidData, err)
	}

	if product.Version <= 0 {
		return errMsg.ErrVersionConflict
	}

	err := s.repo.Update(ctx, product)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			exists, errCheck := s.repo.ProductExists(ctx, product.ID)
			if errCheck != nil {
				return fmt.Errorf("%w: %v", errMsg.ErrGet, errCheck)
			}

			if !exists {
				return errMsg.ErrNotFound
			}
			return errMsg.ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *productService) UpdateStock(ctx context.Context, id int64, quantity int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if quantity <= 0 {
		return errMsg.ErrInvalidQuantity
	}

	err := s.repo.UpdateStock(ctx, id, quantity)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productService) IncreaseStock(ctx context.Context, id int64, amount int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if amount <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.IncreaseStock(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil

}

func (s *productService) DecreaseStock(ctx context.Context, id int64, amount int) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	if amount <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DecreaseStock(ctx, id, amount)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (s *productService) GetStock(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return 0, errMsg.ErrZeroID
	}

	stock, err := s.repo.GetStock(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return stock, nil
}

func (s *productService) EnableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.EnableDiscount(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	return nil
}

func (s *productService) DisableDiscount(ctx context.Context, id int64) error {
	if id <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.repo.DisableDiscount(ctx, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	return nil
}

func (s *productService) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	if id <= 0 {
		return nil, errMsg.ErrZeroID
	}

	if percent <= 0 {
		return nil, errMsg.ErrPercentInvalid
	}

	product, err := s.repo.ApplyDiscount(ctx, id, percent)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return product, nil
}
