package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

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
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductService(mockRepo)

		limit := 0   // inválido
		offset := 10 // qualquer valor válido

		result, err := service.GetAll(ctx, limit, offset)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidLimit)
		mockRepo.AssertNotCalled(t, "GetAll")
	})

	t.Run("falha: offset inválido", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)
		service := NewProductService(mockRepo)

		result, err := service.GetAll(ctx, 10, -1)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidOffset)
		mockRepo.AssertNotCalled(t, "GetAll")
	})

	t.Run("erro do repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

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
	mockRepo := new(mockProduct.ProductMock)

	service := NewProductService(mockRepo)

	t.Run("GetByID com ID inválido", func(t *testing.T) {
		result, err := service.GetByID(ctx, 0)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
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
		assert.ErrorIs(t, err, errMsg.ErrZeroID) // ✅ usa ErrorIs
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
		mockRepo := new(mockProduct.ProductMock)

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
		mockRepo := new(mockProduct.ProductMock)

		service := NewProductService(mockRepo)

		name := "   "

		result, err := service.GetByName(ctx, name)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "nome inválido")
		mockRepo.AssertNotCalled(t, "GetByName")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

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
		mockRepo := new(mockProduct.ProductMock)

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
		mockRepo := new(mockProduct.ProductMock)

		service := NewProductService(mockRepo)

		manufacturer := "   "

		result, err := service.GetByManufacturer(ctx, manufacturer)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, "fabricante inválido")
		mockRepo.AssertNotCalled(t, "GetByManufacturer")
	})

	t.Run("erro no repositório", func(t *testing.T) {
		mockRepo := new(mockProduct.ProductMock)

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

func TestProductService_GetVersionByID(t *testing.T) {
	t.Parallel()

	newService := func() (*mockProduct.ProductMock, ProductService) {
		mr := new(mockProduct.ProductMock)

		return mr, NewProductService(mr)
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo, service := newService()

		invalidID := int64(0) // inválido
		version, err := service.GetVersionByID(context.Background(), invalidID)

		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
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

func TestProductService_ProductExists(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(mockProduct.ProductMock)
	service := NewProductService(mockRepo)

	t.Run("ProductExists com ID inválido", func(t *testing.T) {
		exists, err := service.ProductExists(ctx, 0)
		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("ProductExists com erro no repositório", func(t *testing.T) {
		mockRepo.On("ProductExists", ctx, int64(1)).
			Return(false, fmt.Errorf("erro no banco")).
			Once()

		exists, err := service.ProductExists(ctx, 1)
		assert.False(t, exists)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ProductExists com sucesso", func(t *testing.T) {
		mockRepo.On("ProductExists", ctx, int64(1)).
			Return(true, nil).
			Once()

		exists, err := service.ProductExists(ctx, 1)
		assert.NoError(t, err)
		assert.True(t, exists)
		mockRepo.AssertExpectations(t)
	})
}
