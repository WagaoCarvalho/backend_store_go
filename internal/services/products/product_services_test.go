package services

import (
	"context"
	"errors"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetById(ctx context.Context, id int64) (models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByName(ctx context.Context, name string) ([]models.Product, error) {
	args := m.Called(ctx, name)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error) {
	args := m.Called(ctx, manufacturer)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product models.Product) (models.Product, error) {
	args := m.Called(ctx, product)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockProductRepository) DeleteById(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	args := m.Called(ctx, min, max)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	args := m.Called(ctx, min, max)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockProductRepository) GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error) {
	args := m.Called(ctx, threshold)
	return args.Get(0).([]models.Product), args.Error(1)
}

func TestGetProductById_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProduct := models.Product{ID: 1, ProductName: "Produto A"}
	mockRepo.On("GetById", ctx, int64(1)).Return(expectedProduct, nil)

	product, err := svc.GetById(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByName_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{ID: 1, ProductName: "Produto A"},
		{ID: 2, ProductName: "Produto AA"},
	}
	mockRepo.On("GetByName", ctx, "Produto").Return(expectedProducts, nil)

	products, err := svc.GetByName(ctx, "Produto")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByName_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByName", ctx, "Produto").Return([]models.Product{}, errors.New("erro no banco de dados"))

	_, err := svc.GetByName(ctx, "Produto")

	assert.Error(t, err)
	// Atualizado para "erro ao obter" para corresponder à implementação
	assert.Equal(t, "erro ao obter produtos por nome", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByManufacturer_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{ID: 1, Manufacturer: "Fabricante A"},
		{ID: 2, Manufacturer: "Fabricante A"},
	}
	mockRepo.On("GetByManufacturer", ctx, "Fabricante A").Return(expectedProducts, nil)

	products, err := svc.GetByManufacturer(ctx, "Fabricante A")

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	updatedProduct := models.Product{
		ID:           1,
		ProductName:  "Produto Atualizado",
		Manufacturer: "Fabricante A",
		CostPrice:    100.0,
		SalePrice:    150.0,
	}

	mockRepo.On("Update", ctx, updatedProduct).Return(updatedProduct, nil)

	result, err := svc.Update(ctx, updatedProduct)

	assert.NoError(t, err)
	assert.Equal(t, updatedProduct, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	updatedProduct := models.Product{ID: 1}
	mockRepo.On("Update", ctx, updatedProduct).Return(models.Product{}, errors.New("erro ao atualizar"))

	_, err := svc.Update(ctx, updatedProduct)

	assert.Error(t, err)
	assert.Equal(t, "erro ao atualizar produto", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteProductById_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("DeleteById", ctx, int64(1)).Return(errors.New("erro ao deletar"))

	err := svc.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, "erro ao deletar produto", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByCostPriceRange_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{ID: 1, CostPrice: 50.0},
		{ID: 2, CostPrice: 80.0},
	}
	mockRepo.On("GetByCostPriceRange", ctx, 40.0, 100.0).Return(expectedProducts, nil)

	products, err := svc.GetByCostPriceRange(ctx, 40.0, 100.0)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsByCostPriceRange_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetByCostPriceRange", ctx, 40.0, 100.0).Return([]models.Product{}, errors.New("erro no banco"))

	_, err := svc.GetByCostPriceRange(ctx, 40.0, 100.0)

	assert.Error(t, err)
	assert.Equal(t, "erro ao obter produtos por faixa de preço de custo", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetProductsBySalePriceRange_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{ID: 3, SalePrice: 300.0},
		{ID: 4, SalePrice: 400.0},
	}
	mockRepo.On("GetBySalePriceRange", ctx, 250.0, 450.0).Return(expectedProducts, nil)

	products, err := svc.GetBySalePriceRange(ctx, 250.0, 450.0)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsBySalePriceRange_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetBySalePriceRange", ctx, 250.0, 450.0).Return([]models.Product{}, errors.New("erro no banco"))

	_, err := svc.GetBySalePriceRange(ctx, 250.0, 450.0)

	assert.Error(t, err)
	assert.Equal(t, "erro ao obter produtos por faixa de preço de venda", err.Error())
	mockRepo.AssertExpectations(t)
}
func TestGetProductsLowInStock_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	expectedProducts := []models.Product{
		{ID: 1, StockQuantity: 5},
		{ID: 2, StockQuantity: 10},
	}
	mockRepo.On("GetLowInStock", ctx, 10).Return(expectedProducts, nil)

	products, err := svc.GetLowInStock(ctx, 10)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	mockRepo.AssertExpectations(t)
}

func TestGetProductsLowInStock_Error(t *testing.T) {
	mockRepo := new(MockProductRepository)
	svc := NewProductService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetLowInStock", ctx, 10).Return([]models.Product{}, errors.New("erro no banco"))

	_, err := svc.GetLowInStock(ctx, 10)

	assert.Error(t, err)
	assert.Equal(t, "erro ao buscar produtos com estoque baixo", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_ValidationError(t *testing.T) {
	tests := []struct {
		name          string
		product       models.Product
		expectedError string
	}{
		{
			name: "Nome vazio",
			product: models.Product{
				ProductName:  "",
				Manufacturer: "Fabricante",
				CostPrice:    10.0,
				SalePrice:    15.0,
			},
			expectedError: "nome do produto é obrigatório",
		},
		{
			name: "Preço de custo negativo",
			product: models.Product{
				ProductName:  "Produto",
				Manufacturer: "Fabricante",
				CostPrice:    -10.0,
				SalePrice:    15.0,
			},
			expectedError: "preço de custo deve ser positivo",
		},
		{
			name: "Fabricante vazio",
			product: models.Product{
				ProductName:  "Produto",
				Manufacturer: "",
				CostPrice:    10.0,
				SalePrice:    15.0,
			},
			expectedError: "fabricante é obrigatório",
		},
		{
			name: "Preço de venda inválido",
			product: models.Product{
				ProductName:  "Produto",
				Manufacturer: "Fabricante",
				CostPrice:    10.0,
				SalePrice:    5.0,
			},
			expectedError: "preço de venda deve ser maior que o preço de custo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			svc := NewProductService(mockRepo)
			ctx := context.Background()

			_, err := svc.Create(ctx, tt.product)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
			mockRepo.AssertNotCalled(t, "Create")
		})
	}
}
