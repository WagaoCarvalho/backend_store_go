package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
	service "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations"
	"github.com/WagaoCarvalho/backend_store_go/logger"
)

func TestSupplierCategoryRelationService_Create(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)
		expected := &models.SupplierCategoryRelations{SupplierID: supplierID, CategoryID: categoryID}

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.SupplierCategoryRelations")).Return(expected, nil)

		result, created, err := svc.Create(ctx, supplierID, categoryID)

		assert.NoError(t, err)
		assert.True(t, created)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		tests := []struct {
			name       string
			supplierID int64
			categoryID int64
		}{
			{"zero supplierID", 0, 1},
			{"zero categoryID", 1, 0},
			{"negative supplierID", -1, 2},
			{"negative categoryID", 3, -5},
			{"both zero", 0, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, created, err := svc.Create(ctx, tt.supplierID, tt.categoryID)
				assert.Nil(t, result)
				assert.False(t, created)
				assert.ErrorIs(t, err, service.ErrInvalidRelationData)
			})
		}
	})

	t.Run("relation exists", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, service.ErrRelationExists)
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to check relation", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, errors.New("db error"))

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, service.ErrCheckRelationExists)
		assert.ErrorContains(t, err, "db error")
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to create", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
		mockRepo.On("Create", ctx, mock.Anything).Return(nil, errors.New("db error"))

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, service.ErrCreateRelation)
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_GetBySupplierId(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		expected := []*models.SupplierCategoryRelations{
			{SupplierID: supplierID, CategoryID: 2},
		}

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(expected, nil)

		result, err := svc.GetBySupplierId(ctx, supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		_, err := svc.GetBySupplierId(ctx, 0)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("db error"))

		_, err := svc.GetBySupplierId(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_GetByCategoryId(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		categoryID := int64(2)
		expected := []*models.SupplierCategoryRelations{
			{SupplierID: 1, CategoryID: categoryID},
		}

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(expected, nil)

		result, err := svc.GetByCategoryId(ctx, categoryID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		_, err := svc.GetByCategoryId(ctx, -1)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		categoryID := int64(2)

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("db error"))

		_, err := svc.GetByCategoryId(ctx, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_HasRelation(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relation exists", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

		result, err := svc.HasRelation(ctx, supplierID, categoryID)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		tests := []struct {
			name       string
			supplierID int64
			categoryID int64
		}{
			{"zero supplierID", 0, 1},
			{"zero categoryID", 1, 0},
			{"negative supplierID", -1, 2},
			{"negative categoryID", 3, -5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := svc.HasRelation(ctx, tt.supplierID, tt.categoryID)
				assert.ErrorIs(t, err, service.ErrInvalidRelationData)
			})
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, errors.New("db error"))

		_, err := svc.HasRelation(ctx, supplierID, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_DeleteById(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(nil)

		err := svc.DeleteById(ctx, supplierID, categoryID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		tests := []struct {
			name       string
			supplierID int64
			categoryID int64
		}{
			{"zero supplierID", 0, 1},
			{"zero categoryID", 1, 0},
			{"negative supplierID", -1, 2},
			{"negative categoryID", 3, -5},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := svc.DeleteById(ctx, tt.supplierID, tt.categoryID)
				assert.ErrorIs(t, err, service.ErrInvalidRelationData)
			})
		}
	})

	t.Run("relation not found", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(repository.ErrRelationNotFound)

		err := svc.DeleteById(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, service.ErrRelationNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(errors.New("db error"))

		err := svc.DeleteById(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, service.ErrDeleteRelation)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_DeleteAllBySupplierId(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(1)

		mockRepo.On("DeleteAllBySupplierId", ctx, supplierID).Return(nil)

		err := svc.DeleteAllBySupplierId(ctx, supplierID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		err := svc.DeleteAllBySupplierId(ctx, 0)
		assert.ErrorIs(t, err, service.ErrInvalidRelationData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repository.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo, log)

		supplierID := int64(99)

		mockRepo.On("DeleteAllBySupplierId", ctx, supplierID).Return(errors.New("db error"))

		err := svc.DeleteAllBySupplierId(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
