package services

import (
	"context"
	"errors"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/product"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
)

var (
	ErrProductFetch               = errors.New("erro ao obter produtos")
	ErrProductFetchByID           = errors.New("erro ao obter produto")
	ErrProductFetchByName         = errors.New("erro ao obter produtos por nome")
	ErrProductFetchByManufacturer = errors.New("erro ao obter produtos por fabricante")
	ErrProductCreateNameRequired  = errors.New("validação falhou: nome do produto é obrigatório")
	ErrProductCreateCostPrice     = errors.New("validação falhou: preço de custo deve ser positivo")
	ErrProductCreateManufacturer  = errors.New("validação falhou: fabricante é obrigatório")
	ErrProductCreatePriceLogic    = errors.New("validação falhou: preço de venda deve ser maior que o preço de custo")
	ErrProductUpdate              = errors.New("erro ao atualizar produto")
	ErrProductDelete              = errors.New("erro ao deletar produto")
	ErrProductFetchByCostPrice    = errors.New("erro ao obter produtos por faixa de preço de custo")
	ErrProductFetchBySalePrice    = errors.New("erro ao obter produtos por faixa de preço de venda")
	ErrProductLowStock            = errors.New("erro ao buscar produtos com estoque baixo")
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetById(ctx context.Context, id int64) (models.Product, error)
	GetByName(ctx context.Context, name string) ([]models.Product, error)
	GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error)
	Create(ctx context.Context, product models.Product) (models.Product, error)
	Update(ctx context.Context, product models.Product) (models.Product, error)
	Delete(ctx context.Context, id int64) error
	GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error)
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	products, err := s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, ErrProductFetch
	}
	return products, nil
}

func (s *productService) GetById(ctx context.Context, id int64) (models.Product, error) {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return models.Product{}, ErrProductFetchByID
	}
	return product, nil
}

func (s *productService) GetByName(ctx context.Context, name string) ([]models.Product, error) {
	products, err := s.productRepo.GetByName(ctx, name)
	if err != nil {
		return nil, ErrProductFetchByName
	}
	return products, nil
}

func (s *productService) GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error) {
	products, err := s.productRepo.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		return nil, ErrProductFetchByManufacturer
	}
	return products, nil
}

func (s *productService) Create(ctx context.Context, product models.Product) (models.Product, error) {
	if product.ProductName == "" {
		return models.Product{}, ErrProductCreateNameRequired
	}
	if product.CostPrice <= 0 {
		return models.Product{}, ErrProductCreateCostPrice
	}
	if product.Manufacturer == "" {
		return models.Product{}, ErrProductCreateManufacturer
	}
	if product.SalePrice <= product.CostPrice {
		return models.Product{}, ErrProductCreatePriceLogic
	}
	return s.productRepo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, product models.Product) (models.Product, error) {
	updatedProduct, err := s.productRepo.Update(ctx, product)
	if err != nil {
		return models.Product{}, ErrProductUpdate
	}
	return updatedProduct, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	if err := s.productRepo.DeleteById(ctx, id); err != nil {
		return ErrProductDelete
	}
	return nil
}

func (s *productService) GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetByCostPriceRange(ctx, min, max)
	if err != nil {
		return nil, ErrProductFetchByCostPrice
	}
	return products, nil
}

func (s *productService) GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetBySalePriceRange(ctx, min, max)
	if err != nil {
		return nil, ErrProductFetchBySalePrice
	}
	return products, nil
}

func (s *productService) GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error) {
	products, err := s.productRepo.GetLowInStock(ctx, threshold)
	if err != nil {
		return nil, ErrProductLowStock
	}
	return products, nil
}
