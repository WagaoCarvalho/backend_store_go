package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
)

type ProductService interface {
	GetProducts(ctx context.Context) ([]models.Product, error)
	GetProductById(ctx context.Context, id int64) (models.Product, error)
	GetProductsByName(ctx context.Context, name string) ([]models.Product, error)
	GetProductsByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) (models.Product, error)
	DeleteProductById(ctx context.Context, id int64) error
	GetProductsBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetProductsByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetProductsLowInStock(ctx context.Context, threshold int) ([]models.Product, error)
}

type productService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.productRepo.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos: %w", err)
	}
	return products, nil
}

func (s *productService) GetProductById(ctx context.Context, id int64) (models.Product, error) {
	product, err := s.productRepo.GetProductById(ctx, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("erro ao obter produto: %w", err)
	}
	return product, nil
}

func (s *productService) GetProductsByName(ctx context.Context, name string) ([]models.Product, error) {
	products, err := s.productRepo.GetProductsByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por nome: %w", err)
	}
	return products, nil
}

func (s *productService) GetProductsByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error) {
	products, err := s.productRepo.GetProductsByManufacturer(ctx, manufacturer)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por fabricante: %w", err)
	}
	return products, nil
}

// product_service.go

func (s *productService) CreateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	// Validação do nome do produto
	if product.ProductName == "" {
		return models.Product{}, errors.New("validação falhou: nome do produto é obrigatório")
	}

	// Validação do preço de custo
	if product.CostPrice <= 0 {
		return models.Product{}, errors.New("validação falhou: preço de custo deve ser positivo")
	}

	// Validação do fabricante
	if product.Manufacturer == "" {
		return models.Product{}, errors.New("validação falhou: fabricante é obrigatório")
	}

	// Validação do preço de venda
	if product.SalePrice <= product.CostPrice {
		return models.Product{}, errors.New("validação falhou: preço de venda deve ser maior que o preço de custo")
	}

	// Só chama o repositório se todas as validações passarem
	return s.productRepo.CreateProduct(ctx, product)
}

func (s *productService) UpdateProduct(ctx context.Context, product models.Product) (models.Product, error) {
	// Aqui você pode adicionar validações ou lógicas adicionais antes de atualizar o produto
	updatedProduct, err := s.productRepo.UpdateProduct(ctx, product)
	if err != nil {
		return models.Product{}, fmt.Errorf("erro ao atualizar produto: %w", err)
	}
	return updatedProduct, nil
}

func (s *productService) DeleteProductById(ctx context.Context, id int64) error {
	err := s.productRepo.DeleteProductById(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}
	return nil
}

func (s *productService) GetProductsByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetProductsByCostPriceRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por faixa de preço de custo: %w", err)
	}
	return products, nil
}

func (s *productService) GetProductsBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetProductsBySalePriceRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por faixa de preço de venda: %w", err)
	}
	return products, nil
}

func (s *productService) GetProductsLowInStock(ctx context.Context, threshold int) ([]models.Product, error) {
	products, err := s.productRepo.GetProductsLowInStock(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos com estoque baixo: %w", err)
	}
	return products, nil
}
