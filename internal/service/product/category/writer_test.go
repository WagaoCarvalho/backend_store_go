package services

import (
	"context"
	"errors"
	"testing"

	mockProductCat "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)

		inputCategory := &models.ProductCategory{Name: "NewCategory", Description: "NewDesc"}
		createdCategory := &models.ProductCategory{ID: 1, Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.ProductCategory) bool {
			return cat.Name == inputCategory.Name && cat.Description == inputCategory.Description
		})).Return(createdCategory, nil)

		service := NewProductCategoryService(mockRepo)
		category, err := service.Create(context.Background(), inputCategory)

		assert.NoError(t, err)
		assert.Equal(t, createdCategory, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		service := NewProductCategoryService(mockRepo)

		invalidCategory := &models.ProductCategory{Name: "   "} // nome só com espaços

		category, err := service.Create(context.Background(), invalidCategory)

		assert.Nil(t, category)
		// Service deve retornar erro de validação específico, não ErrInvalidData genérico
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro no campo 'name'") // Erro específico da validação
	})

	t.Run("ErrorOnCreate", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		inputCategory := &models.ProductCategory{Name: "NewCategory", Description: "NewDesc"}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.ProductCategory) bool {
			return cat.Name == inputCategory.Name && cat.Description == inputCategory.Description
		})).Return(nil, errors.New("erro ao criar"))

		service := NewProductCategoryService(mockRepo)
		category, err := service.Create(context.Background(), inputCategory)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar")
		assert.Nil(t, category)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ErrDuplicate", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		service := NewProductCategoryService(mockRepo)

		inputCategory := &models.ProductCategory{
			Name:        "Duplicated",
			Description: "Desc",
		}

		// Usar ErrDuplicate que é o erro do repo
		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.ProductCategory) bool {
			return cat.Name == inputCategory.Name &&
				cat.Description == inputCategory.Description
		})).Return(nil, errMsg.ErrDuplicate)

		category, err := service.Create(context.Background(), inputCategory)

		assert.Nil(t, category)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro já cadastrado") // Service deve propagar o erro
		mockRepo.AssertExpectations(t)
	})
}

func TestProductCategoryService_Update(t *testing.T) {
	mockRepo := new(mockProductCat.MockProductCategory)
	service := NewProductCategoryService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		category := &models.ProductCategory{
			ID:          1,
			Name:        "UpdatedCategory",
			Description: "UpdatedDesc",
		}

		mockRepo.On("Update", ctx, category).Return(nil).Once()

		err := service.Update(ctx, category)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		category := &models.ProductCategory{
			ID:          2,
			Name:        "Fail",
			Description: "Desc",
		}

		dbErr := errors.New("falha ao atualizar")

		mockRepo.On("Update", ctx, category).Return(dbErr).Once()

		err := service.Update(ctx, category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "falha ao atualizar")
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		ctx := context.Background()
		category := &models.ProductCategory{
			ID:          3,
			Name:        "Missing",
			Description: "Missing Desc",
		}

		mockRepo.On("Update", ctx, category).Return(errMsg.ErrNotFound).Once()

		err := service.Update(ctx, category)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		ctx := context.Background()
		category := &models.ProductCategory{
			ID:          4,
			Name:        "",
			Description: "Sem nome",
		}

		err := service.Update(ctx, category)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "campo obrigatório") // Erro específico
	})

	t.Run("InvalidCategoryID", func(t *testing.T) {
		category := &models.ProductCategory{ID: 0}

		err := service.Update(context.Background(), category)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})
}

func TestProductCategoryService_Delete(t *testing.T) {
	mockRepo := new(mockProductCat.MockProductCategory)
	service := NewProductCategoryService(mockRepo)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		id := int64(1)

		// Service corrigido não chama GetByID, apenas Delete
		mockRepo.On("Delete", ctx, id).Return(nil).Once()

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		ctx := context.Background()
		id := int64(3)

		mockRepo.On("Delete", ctx, id).Return(errMsg.ErrNotFound).Once()

		err := service.Delete(ctx, id)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		id := int64(2)
		repoErr := errors.New("db delete error")

		mockRepo.On("Delete", ctx, id).Return(repoErr).Once()

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db delete error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})
}
