package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/product"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/products"
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
		return nil, fmt.Errorf("erro ao obter produtos: %w", err)
	}
	return products, nil
}

func (s *productService) GetById(ctx context.Context, id int64) (models.Product, error) {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return models.Product{}, fmt.Errorf("erro ao obter produto: %w", err)
	}
	return product, nil
}

func (s *productService) GetByName(ctx context.Context, name string) ([]models.Product, error) {
	products, err := s.productRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por nome: %w", err)
	}
	return products, nil
}

func (s *productService) GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error) {
	products, err := s.productRepo.GetByManufacturer(ctx, manufacturer)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por fabricante: %w", err)
	}
	return products, nil
}

// product_service.go

func (s *productService) Create(ctx context.Context, product models.Product) (models.Product, error) {
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
	return s.productRepo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, product models.Product) (models.Product, error) {
	// Aqui você pode adicionar validações ou lógicas adicionais antes de atualizar o produto
	updatedProduct, err := s.productRepo.Update(ctx, product)
	if err != nil {
		return models.Product{}, fmt.Errorf("erro ao atualizar produto: %w", err)
	}
	return updatedProduct, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	err := s.productRepo.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}
	return nil
}

func (s *productService) GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetByCostPriceRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por faixa de preço de custo: %w", err)
	}
	return products, nil
}

func (s *productService) GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	products, err := s.productRepo.GetBySalePriceRange(ctx, min, max)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter produtos por faixa de preço de venda: %w", err)
	}
	return products, nil
}

func (s *productService) GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error) {
	products, err := s.productRepo.GetLowInStock(ctx, threshold)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos com estoque baixo: %w", err)
	}
	return products, nil
}
