package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

func newTestLogger() *logger.LoggerAdapter {
	logrusLogger := logrus.New()
	logrusLogger.SetLevel(logrus.FatalLevel)
	return logger.NewLoggerAdapter(logrusLogger)
}

func TestProductService_Create(t *testing.T) {
	ctx := context.Background()

	// Produto válido base para o teste
	validProduct := func() *models.Product {
		return &models.Product{
			ProductName:   "Produto Teste",
			Manufacturer:  "Fabricante X",
			SupplierID:    utils.Int64Ptr(1),
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 5,
			Status:        true,
		}
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()

		mockRepo.On("Create", ctx, input).Return(&models.Product{
			ID:            1,
			ProductName:   input.ProductName,
			Manufacturer:  input.Manufacturer,
			SupplierID:    input.SupplierID,
			CostPrice:     input.CostPrice,
			SalePrice:     input.SalePrice,
			StockQuantity: input.StockQuantity,
			Status:        input.Status,
		}, nil).Once()

		created, err := service.Create(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, created)
		assert.Equal(t, int64(1), created.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("produto inválido", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := &models.Product{} // faltam todos os campos obrigatórios

		created, err := service.Create(ctx, input)

		assert.Nil(t, created)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, models.ErrInvalidProductName))
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()

		mockErr := errors.New("erro no repositório")
		mockRepo.On("Create", ctx, input).Return(nil, mockErr).Once()

		created, err := service.Create(ctx, input)

		assert.Nil(t, created)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		limit := 10
		offset := 0

		expectedProducts := []*models.Product{
			{ID: 1, ProductName: "Produto 1"},
			{ID: 2, ProductName: "Produto 2"},
		}

		mockRepo.
			On("GetAll", ctx, limit, offset).
			Return(expectedProducts, nil)

		result, err := service.GetAll(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, result, len(expectedProducts))
		assert.Equal(t, expectedProducts, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		limit := 10
		offset := 0

		mockErr := errors.New("erro ao buscar produtos")
		mockRepo.
			On("GetAll", ctx, limit, offset).
			Return(nil, mockErr)

		result, err := service.GetAll(ctx, limit, offset)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetById(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		id := int64(1)
		expectedProduct := &models.Product{
			ID:          id,
			ProductName: "Produto 1",
		}

		mockRepo.
			On("GetById", ctx, id).
			Return(expectedProduct, nil)

		result, err := service.GetById(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("id inválido", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		id := int64(0)

		result, err := service.GetById(ctx, id)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "ID inválido")
		mockRepo.AssertNotCalled(t, "GetById")
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		id := int64(99)
		mockErr := errors.New("erro ao buscar produto")

		mockRepo.
			On("GetById", ctx, id).
			Return(nil, mockErr)

		result, err := service.GetById(ctx, id)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByName(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		name := "Produto X"
		expectedProducts := []*models.Product{
			{ID: 1, ProductName: name},
			{ID: 2, ProductName: name},
		}

		mockRepo.
			On("GetByName", ctx, name).
			Return(expectedProducts, nil)

		result, err := service.GetByName(ctx, name)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedProducts, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("nome inválido (string vazia)", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		name := "   "

		result, err := service.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "nome inválido")
		mockRepo.AssertNotCalled(t, "GetByName")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		name := "Produto X"
		mockErr := errors.New("erro no banco")

		mockRepo.
			On("GetByName", ctx, name).
			Return(nil, mockErr)

		result, err := service.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByManufacturer(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		manufacturer := "Fabricante X"
		expectedProducts := []*models.Product{
			{ID: 1, Manufacturer: manufacturer},
			{ID: 2, Manufacturer: manufacturer},
		}

		mockRepo.
			On("GetByManufacturer", ctx, manufacturer).
			Return(expectedProducts, nil)

		result, err := service.GetByManufacturer(ctx, manufacturer)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, expectedProducts, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fabricante inválido (string vazia)", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		manufacturer := "   "

		result, err := service.GetByManufacturer(ctx, manufacturer)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "fabricante inválido")
		mockRepo.AssertNotCalled(t, "GetByManufacturer")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		manufacturer := "Fabricante X"
		mockErr := errors.New("erro no banco")

		mockRepo.
			On("GetByManufacturer", ctx, manufacturer).
			Return(nil, mockErr)

		result, err := service.GetByManufacturer(ctx, manufacturer)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Update(t *testing.T) {
	ctx := context.Background()

	validProduct := func() *models.Product {
		return &models.Product{
			ID:            1,
			ProductName:   "Produto Atualizado",
			Manufacturer:  "Fabricante X",
			SupplierID:    utils.Int64Ptr(1),
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 5,
			Status:        true,
			Version:       1, // já inicializa versão aqui
		}
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()

		mockRepo.
			On("Update", ctx, input).
			Return(input, nil).Once()

		updated, err := service.Update(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, int64(1), updated.ID)
		assert.Equal(t, "Produto Atualizado", updated.ProductName)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro de validação: nome inválido", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()
		input.ProductName = "" // invalida nome

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, ErrInvalidProduct)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("erro de validação: versão inválida", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()
		input.Version = 0 // versão inválida

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, ErrInvalidVersion)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("erro do repositório: produto não encontrado", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()

		mockRepo.
			On("Update", ctx, input).
			Return(nil, repo.ErrProductNotFound).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, repo.ErrProductNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório: conflito de versão", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()

		mockRepo.
			On("Update", ctx, input).
			Return(nil, repo.ErrVersionConflict).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, repo.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório genérico", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		input := validProduct()
		mockErr := errors.New("erro no repositório")

		mockRepo.
			On("Update", ctx, input).
			Return(nil, mockErr).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar produto")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_DisableProduct(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.ProductRepositoryMock, ProductService) {
		mockRepo := new(repo.ProductRepositoryMock)
		service := NewProductService(mockRepo, logger)
		return mockRepo, service
	}

	t.Run("Deve desabilitar produto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableProduct", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DisableProduct(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao desabilitar produto", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableProduct", mock.Anything, int64(2)).Return(fmt.Errorf("erro banco")).Once()

		err := service.DisableProduct(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao desativar produto")
		assert.ErrorContains(t, err, "erro banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_EnableProduct(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.ProductRepositoryMock, ProductService) {
		mockRepo := new(repo.ProductRepositoryMock)
		service := NewProductService(mockRepo, logger)
		return mockRepo, service
	}

	t.Run("Deve habilitar produto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableProduct", mock.Anything, int64(1)).Return(nil).Once()

		err := service.EnableProduct(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro ao habilitar produto", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableProduct", mock.Anything, int64(2)).Return(fmt.Errorf("erro banco")).Once()

		err := service.EnableProduct(context.Background(), 2)

		assert.ErrorContains(t, err, "erro ao ativar produto")
		assert.ErrorContains(t, err, "erro banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		id := int64(1)

		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := NewProductService(mockRepo, log)

		id := int64(1)
		mockErr := errors.New("erro no repositório")

		mockRepo.On("Delete", ctx, id).Return(mockErr)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetVersionByID(t *testing.T) {
	t.Parallel()

	newService := func() (*repo.ProductRepositoryMock, ProductService) {
		mr := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		return mr, NewProductService(mr, log)
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil)

		version, err := service.GetVersionByID(context.Background(), 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), repo.ErrProductNotFound)

		version, err := service.GetVersionByID(context.Background(), 2)
		assert.ErrorIs(t, err, repo.ErrProductNotFound)
		assert.Equal(t, int64(0), version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		t.Parallel()

		mockRepo, service := newService()
		mockRepo.On("GetVersionByID", mock.Anything, int64(3)).Return(int64(0), errors.New("db failure"))

		version, err := service.GetVersionByID(context.Background(), 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db failure")
		assert.True(t, errors.Is(err, ErrInvalidVersion))
		assert.Equal(t, int64(0), version)

		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateStock(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.ProductRepositoryMock, ProductService) {
		mockRepo := new(repo.ProductRepositoryMock)
		service := NewProductService(mockRepo, logger)
		return mockRepo, service
	}

	t.Run("Deve atualizar o estoque com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("UpdateStock", mock.Anything, int64(1), 25).Return(nil).Once()

		err := service.UpdateStock(context.Background(), 1, 25)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando repo falhar", func(t *testing.T) {
		mockRepo, service := setup()
		expectedErr := fmt.Errorf("erro de banco")

		mockRepo.On("UpdateStock", mock.Anything, int64(1), 25).Return(expectedErr).Once()

		err := service.UpdateStock(context.Background(), 1, 25)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_IncreaseStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := productService{repo: repoMock, logger: log}

		repoMock.On("IncreaseStock", ctx, int64(1), 10).Return(repo.ErrProductNotFound)

		err := service.IncreaseStock(ctx, 1, 10)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrUpdateStock)
		repoMock.AssertExpectations(t)
	})

	t.Run("Deve aumentar estoque com sucesso", func(t *testing.T) {
		repoMock := new(repo.ProductRepositoryMock)
		log := newTestLogger()
		service := productService{repo: repoMock, logger: log}

		repoMock.On("IncreaseStock", ctx, int64(1), 5).Return(nil)

		err := service.IncreaseStock(ctx, 1, 5)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}
