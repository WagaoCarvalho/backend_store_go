package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

type ProductService interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	GetById(ctx context.Context, id int64) (*models.Product, error)
	GetByName(ctx context.Context, name string) ([]*models.Product, error)
	GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error)
	GetVersionByID(ctx context.Context, uid int64) (int64, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id int64) error

	DisableProduct(ctx context.Context, uid int64) error
	EnableProduct(ctx context.Context, uid int64) error

	UpdateStock(ctx context.Context, id int64, quantity int) error
	IncreaseStock(ctx context.Context, id int64, amount int) error
	DecreaseStock(ctx context.Context, id int64, amount int) error
	//GetStock(ctx context.Context, id int64) (int, error)

	//EnableDiscount(ctx context.Context, id int64) error
	//DisableDiscount(ctx context.Context, id int64) error
	//ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error)
}

type productService struct {
	repo   repo.ProductRepository
	logger logger.LoggerAdapterInterface
}

func NewProductService(repo repo.ProductRepository, logger logger.LoggerAdapterInterface) ProductService {
	return &productService{
		repo:   repo,
		logger: logger,
	}
}

func (s *productService) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	ref := "[productService - Create] - "
	s.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":         product.ProductName,
		"manufacturer": product.Manufacturer,
		"supplier_id":  utils.Int64OrNil(product.SupplierID),
	})

	if err := product.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return nil, err
	}

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name": product.ProductName,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"product_id": createdProduct.ID,
	})

	return createdProduct, nil
}

func (s *productService) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	ref := "[productService - GetAll] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"limit":  limit,
		"offset": offset,
	})

	products, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_encontrados": len(products),
	})

	return products, nil
}

func (s *productService) GetById(ctx context.Context, id int64) (*models.Product, error) {
	ref := "[productService - GetById] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"product_id": id,
	})

	if id <= 0 {
		s.logger.Warn(ctx, ref+logger.LogInvalidID, map[string]any{
			"product_id": id,
		})
		return nil, errors.New("ID inválido")
	}

	product, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": product.ID,
	})

	return product, nil
}

func (s *productService) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	ref := "[productService - GetByName] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"name": name,
	})

	if strings.TrimSpace(name) == "" {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"name": name,
		})
		return nil, errors.New("nome inválido")
	}

	products, err := s.repo.GetByName(ctx, name)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"name": name,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"results": len(products),
	})

	return products, nil
}

func (s *productService) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
	ref := "[productService - GetByManufacturer] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"manufacturer": manufacturer,
	})

	if strings.TrimSpace(manufacturer) == "" {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"manufacturer": manufacturer,
		})
		return nil, errors.New("fabricante inválido")
	}

	products, err := s.repo.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"manufacturer": manufacturer,
		})
		return nil, err
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"results": len(products),
	})

	return products, nil
}

func (s *productService) GetVersionByID(ctx context.Context, pid int64) (int64, error) {
	ref := "[productService - GetVersionByID] - "
	s.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"product_id": pid,
	})

	version, err := s.repo.GetVersionByID(ctx, pid)
	if err != nil {
		if errors.Is(err, repo.ErrProductNotFound) {
			s.logger.Error(ctx, err, ref+logger.LogNotFound, map[string]any{
				"product_id": pid,
			})
			return 0, repo.ErrProductNotFound
		}

		s.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": pid,
		})
		return 0, fmt.Errorf("%w: %v", ErrInvalidVersion, err)
	}

	s.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": pid,
		"version":    version,
	})

	return version, nil
}

func (s *productService) DisableProduct(ctx context.Context, uid int64) error {
	ref := "[productService - Disable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": uid,
	})

	err := s.repo.DisableProduct(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDisableProduct, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
		"status":     false,
	})

	return nil
}

func (s *productService) EnableProduct(ctx context.Context, uid int64) error {
	ref := "[productService - Enable] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": uid,
	})

	err := s.repo.EnableProduct(ctx, uid)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrEnableProduct, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": uid,
		"status":     true,
	})

	return nil
}

func (s *productService) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	ref := "[productService - Update] - "

	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id":   product.ID,
		"product_name": product.ProductName,
		"version":      product.Version,
	})

	if err := product.Validate(); err != nil {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"erro": err.Error(),
		})
		return nil, ErrInvalidProduct
	}

	if product.Version <= 0 {
		s.logger.Warn(ctx, ref+logger.LogValidateError, map[string]any{
			"product_id": product.ID,
			"version":    product.Version,
		})
		return nil, ErrInvalidVersion
	}

	updatedProduct, err := s.repo.Update(ctx, product)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrProductNotFound):
			s.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": product.ID,
			})
			return nil, repo.ErrProductNotFound

		case errors.Is(err, repo.ErrVersionConflict):
			s.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"product_id": product.ID,
				"version":    product.Version,
			})
			return nil, repo.ErrVersionConflict

		default:
			s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
				"product_id": product.ID,
			})
			return nil, fmt.Errorf("%w: %v", ErrProductUpdate, err)
		}
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":   updatedProduct.ID,
		"product_name": updatedProduct.ProductName,
		"version":      updatedProduct.Version,
	})

	return updatedProduct, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	ref := "[productService - Delete] - "
	s.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"product_id": id,
	})

	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"product_id": id,
		})
		return err
	}

	s.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"product_id": id,
	})

	return nil
}

func (s *productService) UpdateStock(ctx context.Context, id int64, quantity int) error {
	ref := "[productService - UpdateStock] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": id,
		"quantity":   quantity,
	})

	err := s.repo.UpdateStock(ctx, id, quantity)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": id,
			"quantity":   quantity,
		})
		return fmt.Errorf("%w: %v", ErrUpdateStock, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"quantity":   quantity,
	})

	return nil
}

func (s *productService) IncreaseStock(ctx context.Context, id int64, amount int) error {
	ref := "[productService - IncreaseStock] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	err := s.repo.IncreaseStock(ctx, id, amount)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id":     id,
			"stock_quantity": amount,
		})
		return fmt.Errorf("%w: %v", ErrUpdateStock, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	return nil

}

func (s *productService) DecreaseStock(ctx context.Context, id int64, amount int) error {
	ref := "[productService - DecreaseStock] - "
	s.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	err := s.repo.DecreaseStock(ctx, id, amount)
	if err != nil {
		s.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id":     id,
			"stock_quantity": amount,
		})
		return fmt.Errorf("%w: %v", ErrUpdateStock, err)
	}

	s.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	return nil
}
