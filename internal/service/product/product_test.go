package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mock_product "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func TestProductService_Create(t *testing.T) {
	ctx := context.Background()

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
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := &models.Product{
			Status: true, // faltando campos obrigatórios
		}

		created, err := service.Create(ctx, input)

		assert.Nil(t, created)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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

	t.Run("falha: limit inválido", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		limit := 0   // inválido
		offset := 10 // qualquer valor válido

		result, err := service.GetAll(ctx, limit, offset)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
		mockRepo.AssertNotCalled(t, "GetAll")
	})

	t.Run("falha: offset inválido", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		result, err := service.GetAll(ctx, 10, -1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOffset)
		mockRepo.AssertNotCalled(t, "GetAll")
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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

func TestProductService_GetByID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mock_product.ProductRepositoryMock)

	service := NewProductService(mockRepo)

	t.Run("GetByID com ID inválido", func(t *testing.T) {
		result, err := service.GetByID(ctx, 0)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIDZero)
	})

	t.Run("sucesso", func(t *testing.T) {

		id := int64(1)
		expectedProduct := &models.Product{
			ID:          id,
			ProductName: "Produto 1",
		}

		mockRepo.
			On("GetByID", ctx, id).
			Return(expectedProduct, nil)

		result, err := service.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expectedProduct, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("id inválido", func(t *testing.T) {

		id := int64(0)

		result, err := service.GetByID(ctx, id)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrIDZero) // ✅ usa ErrorIs
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("erro do repositório", func(t *testing.T) {

		id := int64(99)
		mockErr := errors.New("erro ao buscar produto")

		mockRepo.
			On("GetByID", ctx, id).
			Return(nil, mockErr)

		result, err := service.GetByID(ctx, id)

		assert.Nil(t, result)
		assert.EqualError(t, err, mockErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_GetByName(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

		name := "   "

		result, err := service.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "nome inválido")
		mockRepo.AssertNotCalled(t, "GetByName")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

		manufacturer := "   "

		result, err := service.GetByManufacturer(ctx, manufacturer)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "fabricante inválido")
		mockRepo.AssertNotCalled(t, "GetByManufacturer")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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
			Version:       1,
		}
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()
		input.ID = 0

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: validação inválida", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()
		input.ProductName = "" // invalida a validação

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: versão inválida", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()
		input.Version = 0

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: produto não encontrado", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()

		mockRepo.
			On("Update", ctx, input).
			Return(nil, errMsg.ErrNotFound).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: conflito de versão", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()

		mockRepo.
			On("Update", ctx, input).
			Return(nil, errMsg.ErrVersionConflict).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro genérico do repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

		input := validProduct()
		mockErr := errors.New("erro no repositório")

		mockRepo.
			On("Update", ctx, input).
			Return(nil, mockErr).Once()

		updated, err := service.Update(ctx, input)

		assert.Nil(t, updated)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro no repositório")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)

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
}

func TestProductService_DisableProduct(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		invalidID := int64(0)
		err := service.DisableProduct(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "DisableProduct")
	})

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

		assert.ErrorContains(t, err, "erro ao desabilitar")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_EnableProduct(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		invalidID := int64(0)
		err := service.EnableProduct(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "EnableProduct")
	})

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

		assert.ErrorContains(t, err, "erro ao ativar")
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

		id := int64(1)

		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

		invalidID := int64(0)
		err := service.Delete(context.Background(), invalidID)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mock_product.ProductRepositoryMock)

		service := NewProductService(mockRepo)

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

	newService := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mr := new(mock_product.ProductRepositoryMock)

		return mr, NewProductService(mr)
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := newService()

		invalidID := int64(0) // inválido
		version, err := service.GetVersionByID(context.Background(), invalidID)

		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "GetVersionByID")
	})

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
		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errMsg.ErrNotFound)

		version, err := service.GetVersionByID(context.Background(), 2)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
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
		assert.True(t, errors.Is(err, errMsg.ErrVersionConflict))
		assert.Equal(t, int64(0), version)

		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_UpdateStock(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve atualizar o estoque com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("UpdateStock", mock.Anything, int64(1), 25).Return(nil).Once()

		err := service.UpdateStock(context.Background(), 1, 25)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.UpdateStock(context.Background(), 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		mockRepo, service := setup()

		err := service.UpdateStock(context.Background(), 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrInvalidQuantity)
		mockRepo.AssertNotCalled(t, "UpdateStock")
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
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("IncreaseStock", ctx, int64(1), 10).Return(errMsg.ErrNotFound)

		err := service.IncreaseStock(ctx, 1, 10)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		err := service.IncreaseStock(ctx, 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		err := service.IncreaseStock(ctx, 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("Deve aumentar estoque com sucesso", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("IncreaseStock", ctx, int64(1), 5).Return(nil)

		err := service.IncreaseStock(ctx, 1, 5)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestProductService_DecreaseStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("DecreaseStock", ctx, int64(1), 10).Return(errMsg.ErrNotFound)

		err := service.DecreaseStock(ctx, 1, 10)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		err := service.DecreaseStock(ctx, 0, 10)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("falha: quantidade inválida", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		err := service.DecreaseStock(ctx, 1, 0)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("Deve diminuir estoque com sucesso", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("DecreaseStock", ctx, int64(1), 10).Return(nil)

		err := service.DecreaseStock(ctx, 1, 10)

		assert.NoError(t, err)
		repoMock.AssertExpectations(t)
	})
}

func TestProductService_GetStock(t *testing.T) {
	ctx := context.Background()

	t.Run("Deve retornar erro quando repo retornar erro", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("GetStock", ctx, int64(1)).Return(0, fmt.Errorf("erro inesperado"))

		stock, err := service.GetStock(ctx, 1)

		assert.Error(t, err)
		assert.Equal(t, 0, stock)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		repoMock.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		stock, err := service.GetStock(ctx, 0)

		assert.Equal(t, 0, stock)
		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "GetStock")
	})

	t.Run("Deve retornar estoque com sucesso", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)

		service := productService{repo: repoMock}

		repoMock.On("GetStock", ctx, int64(1)).Return(25, nil)

		stock, err := service.GetStock(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, 25, stock)
		repoMock.AssertExpectations(t)
	})
}

func TestProductService_EnableDiscount(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve habilitar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(nil).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		ctx := context.Background() // ⚠️ necessário criar o contexto
		err := service.EnableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "EnableDiscount")
	})

	t.Run("Deve retornar erro quando repo falhar", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("EnableDiscount", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		err := service.EnableDiscount(context.Background(), 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductService_DisableDiscount(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("Deve desabilitar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("DisableDiscount", mock.Anything, int64(1)).Return(nil).Once()

		err := service.DisableDiscount(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: ID inválido", func(t *testing.T) {
		repoMock := new(mock_product.ProductRepositoryMock)
		service := productService{repo: repoMock}

		ctx := context.Background() // ⚠️ necessário criar o contexto
		err := service.DisableDiscount(ctx, 0)

		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		repoMock.AssertNotCalled(t, "DisableDiscount")
	})

	t.Run("Erro: produto não encontrado", func(t *testing.T) {
		mockRepo, service := setup()
		mockRepo.On("EnableDiscount", mock.Anything, int64(2)).Return(errMsg.ErrNotFound).Once()

		err := service.EnableDiscount(context.Background(), 2)

		assert.ErrorIs(t, err, errMsg.ErrProductEnableDiscount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Erro inesperado do repositório", func(t *testing.T) {
		mockRepo, service := setup()
		unexpectedErr := fmt.Errorf("erro de conexão")
		mockRepo.On("DisableDiscount", mock.Anything, int64(3)).Return(unexpectedErr).Once()

		err := service.DisableDiscount(context.Background(), 3)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrProductDisableDiscount)
		assert.Contains(t, err.Error(), "erro de conexão")
		mockRepo.AssertExpectations(t)
	})

}

func TestProductService_ApplyDiscount(t *testing.T) {

	setup := func() (*mock_product.ProductRepositoryMock, ProductService) {
		mockRepo := new(mock_product.ProductRepositoryMock)
		service := NewProductService(mockRepo)
		return mockRepo, service
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		product, err := service.ApplyDiscount(ctx, 0, 10.0) // ID inválido

		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrIDZero)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("falha: percent inválido", func(t *testing.T) {
		mockRepo, service := setup()

		ctx := context.Background()
		product, err := service.ApplyDiscount(ctx, 1, 0) // percent inválido

		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrPercentInvalid)
		mockRepo.AssertNotCalled(t, "ApplyDiscount")
	})

	t.Run("Deve aplicar desconto com sucesso", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(1)
		percent := 10.0

		expectedProduct := &models.Product{
			ID:            productID,
			ProductName:   "Produto Teste",
			SalePrice:     90.0,
			AllowDiscount: true,
		}

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(expectedProduct, nil).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, expectedProduct, product)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro se produto não encontrado", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(99)
		percent := 15.0

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, errMsg.ErrNotFound).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico ao aplicar desconto", func(t *testing.T) {
		mockRepo, service := setup()
		productID := int64(2)
		percent := 5.0
		expectedErr := errors.New("erro inesperado no banco")

		mockRepo.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, expectedErr).
			Once()

		product, err := service.ApplyDiscount(context.Background(), productID, percent)

		assert.Error(t, err)
		assert.Nil(t, product)
		assert.ErrorIs(t, err, errMsg.ErrProductApplyDiscount)
		mockRepo.AssertExpectations(t)
	})
}
