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

	t.Run("ErrInvalidCategoryName", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		service := NewProductCategoryService(mockRepo)

		invalidCategory := &models.ProductCategory{Name: "   "} // nome só com espaços

		category, err := service.Create(context.Background(), invalidCategory)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
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

	t.Run("ErrAlreadyExists", func(t *testing.T) {
		mockRepo := new(mockProductCat.MockProductCategory)
		service := NewProductCategoryService(mockRepo)

		inputCategory := &models.ProductCategory{
			Name:        "Duplicated",
			Description: "Desc",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(cat *models.ProductCategory) bool {
			return cat.Name == inputCategory.Name &&
				cat.Description == inputCategory.Description
		})).Return(nil, errMsg.ErrAlreadyExists)

		category, err := service.Create(context.Background(), inputCategory)

		assert.Nil(t, category)
		assert.ErrorIs(t, err, errMsg.ErrAlreadyExists)

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
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
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
		assert.Contains(t, err.Error(), "campo obrigatório")
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
		category := &models.ProductCategory{ID: 1, Name: "Categoria", Description: "Desc"}

		mockRepo.On("GetByID", ctx, int64(1)).Return(category, nil).Once()
		mockRepo.On("Delete", ctx, int64(1)).Return(nil).Once()

		err := service.Delete(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID repository error", func(t *testing.T) {
		ctx := context.Background()
		id := int64(10)
		dbErr := errors.New("erro inesperado no banco de dados")

		mockRepo.On("GetByID", ctx, id).Return(nil, dbErr).Once()

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.ErrorContains(t, err, errMsg.ErrGet.Error())
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("CategoryNotFound", func(t *testing.T) {
		ctx := context.Background()

		mockRepo.On("GetByID", ctx, int64(3)).Return(nil, errMsg.ErrNotFound).Once()

		err := service.Delete(ctx, 3)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		ctx := context.Background()
		category := &models.ProductCategory{ID: 2, Name: "Categoria", Description: "Desc"}
		repoErr := errors.New("db delete error")

		mockRepo.On("GetByID", ctx, int64(2)).Return(category, nil).Once()
		mockRepo.On("Delete", ctx, int64(2)).Return(repoErr).Once()

		err := service.Delete(ctx, 2)

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errMsg.ErrDelete))
		assert.ErrorContains(t, err, "db delete error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		err := service.Delete(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})
}
