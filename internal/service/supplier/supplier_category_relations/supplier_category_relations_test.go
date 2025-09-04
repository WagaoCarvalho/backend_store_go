package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mock_supplier_cat_rel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_category_relations"
)

func TestSupplierCategoryRelationService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

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
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

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
				assert.ErrorIs(t, err, err_msg.ErrInvalidData)
			})
		}
	})

	t.Run("relation exists", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, err_msg.ErrRelationExists)
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to check relation", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, errors.New("db error"))

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, err_msg.ErrRelationCheck)
		assert.ErrorContains(t, err, "db error")
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to create", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
		mockRepo.On("Create", ctx, mock.Anything).Return(nil, errors.New("db error"))

		_, created, err := svc.Create(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, err_msg.ErrCreate)
		assert.False(t, created)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_GetBySupplierID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		expected := []*models.SupplierCategoryRelations{
			{SupplierID: supplierID, CategoryID: 2},
		}

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(expected, nil)

		result, err := svc.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		_, err := svc.GetBySupplierID(ctx, 0)
		assert.ErrorIs(t, err, err_msg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("db error"))

		_, err := svc.GetBySupplierID(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_GetByCategoryID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		categoryID := int64(2)
		expected := []*models.SupplierCategoryRelations{
			{SupplierID: 1, CategoryID: categoryID},
		}

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(expected, nil)

		result, err := svc.GetByCategoryID(ctx, categoryID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		_, err := svc.GetByCategoryID(ctx, -1)
		assert.ErrorIs(t, err, err_msg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		categoryID := int64(2)

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(([]*models.SupplierCategoryRelations)(nil), errors.New("db error"))

		_, err := svc.GetByCategoryID(ctx, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_HasRelation(t *testing.T) {
	ctx := context.Background()

	t.Run("success - relation exists", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

		result, err := svc.HasRelation(ctx, supplierID, categoryID)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

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
				assert.ErrorIs(t, err, err_msg.ErrInvalidData)
			})
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, errors.New("db error"))

		_, err := svc.HasRelation(ctx, supplierID, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_DeleteByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(nil)

		err := svc.DeleteByID(ctx, supplierID, categoryID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

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
				err := svc.DeleteByID(ctx, tt.supplierID, tt.categoryID)
				assert.ErrorIs(t, err, err_msg.ErrInvalidData)
			})
		}
	})

	t.Run("relation not found", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(err_msg.ErrNotFound)

		err := svc.DeleteByID(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, err_msg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(errors.New("db error"))

		err := svc.DeleteByID(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, err_msg.ErrDelete)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_DeleteAllBySupplierID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)

		mockRepo.On("DeleteAllBySupplierID", ctx, supplierID).Return(nil)

		err := svc.DeleteAllBySupplierID(ctx, supplierID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		err := svc.DeleteAllBySupplierID(ctx, 0)
		assert.ErrorIs(t, err, err_msg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_supplier_cat_rel.MockSupplierCategoryRelationRepo)
		svc := service.NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(99)

		mockRepo.On("DeleteAllBySupplierID", ctx, supplierID).Return(errors.New("db error"))

		err := svc.DeleteAllBySupplierID(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
