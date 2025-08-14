package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/supplier/supplier_categories"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

func TestSupplierCategoryService_Create(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{Name: "Alimentos"}

		mockRepo.On("Create", mock.Anything, category).Return(category, nil)

		result, err := service.Create(ctx, category)

		assert.NoError(t, err)
		assert.Equal(t, category, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid name", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{Name: " "}

		result, err := service.Create(ctx, category)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{Name: "Tecnologia"}

		mockRepo.On("Create", mock.Anything, category).Return((*models.SupplierCategory)(nil), errors.New("erro no banco"))

		result, err := service.Create(ctx, category)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryService_GetByID(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		expected := &models.SupplierCategory{ID: 1, Name: "Eletrônicos"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		result, err := service.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		result, err := service.GetByID(ctx, -1)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		mockRepo.On("GetByID", mock.Anything, int64(999)).Return((*models.SupplierCategory)(nil), repo.ErrSupplierCategoryNotFound)

		result, err := service.GetByID(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		mockRepo.On("GetByID", mock.Anything, int64(2)).Return((*models.SupplierCategory)(nil), errors.New("erro no banco"))

		result, err := service.GetByID(ctx, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryService_Update(t *testing.T) {
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve atualizar com sucesso", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{
			ID:   1,
			Name: "Atualizada",
		}

		mockRepo.On("Update", mock.Anything, category).Return(nil)

		err := service.Update(context.Background(), category)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deve retornar erro se ID for zero", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{
			ID:   0,
			Name: "Nome válido",
		}

		err := service.Update(context.Background(), category)

		assert.Error(t, err)
		assert.Equal(t, ErrCategoryIDRequired, err)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("deve retornar erro se nome for vazio", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		baseLogger := logrus.New()
		log := logger.NewLoggerAdapter(baseLogger)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{
			ID:   1,
			Name: "   ", // inválido
		}

		err := service.Update(context.Background(), category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "campo obrigatório")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("deve propagar erro do repositório", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		category := &models.SupplierCategory{
			ID:   1,
			Name: "Nome válido",
		}

		mockRepo.On("Update", mock.Anything, category).
			Return(fmt.Errorf("erro ao atualizar no repositório"))

		err := service.Update(context.Background(), category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar no repositório")
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryService_Delete(t *testing.T) {
	baseLogger := logrus.New()
	mockRepo := new(repo.MockSupplierCategoryRepo)
	log := logger.NewLoggerAdapter(baseLogger)
	service := NewSupplierCategoryService(mockRepo, log)

	t.Run("deve deletar com sucesso", func(t *testing.T) {
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", mock.Anything, int64(1))
	})

	t.Run("deve retornar erro se id for inválido", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything, int64(0))
	})

	t.Run("deve retornar erro se ocorrer falha ao deletar", func(t *testing.T) {
		expectedErr := errors.New("erro no banco")
		mockRepo.On("Delete", mock.Anything, int64(2)).Return(expectedErr)

		err := service.Delete(context.Background(), 2)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar categorias")
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertCalled(t, "Delete", mock.Anything, int64(2))
	})

}

func TestSupplierCategoryService_GetAll(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		expectedCategories := []*models.SupplierCategory{
			{ID: 1, Name: "Categoria A"},
			{ID: 2, Name: "Categoria B"},
		}

		mockRepo.On("GetAll", ctx).Return(expectedCategories, nil)

		categories, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedCategories, categories)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repo.MockSupplierCategoryRepo)
		service := NewSupplierCategoryService(mockRepo, log)

		mockRepo.On("GetAll", ctx).Return(([]*models.SupplierCategory)(nil), errors.New("erro ao buscar categorias"))

		categories, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, categories)
		assert.Contains(t, err.Error(), "erro ao buscar categorias")
		mockRepo.AssertExpectations(t)
	})
}
