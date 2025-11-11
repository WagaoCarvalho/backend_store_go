package services

import (
	"context"
	"errors"
	"testing"

	mockSupplierCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
)

func TestSupplierCategoryRelationService_GetBySupplierID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		expected := []*models.SupplierCategoryRelation{
			{SupplierID: supplierID, CategoryID: 2},
		}

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(expected, nil)

		result, err := svc.GetBySupplierID(ctx, supplierID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		_, err := svc.GetBySupplierID(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)

		mockRepo.On("GetBySupplierID", ctx, supplierID).Return(([]*models.SupplierCategoryRelation)(nil), errors.New("db error"))

		_, err := svc.GetBySupplierID(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_GetByCategoryID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		categoryID := int64(2)
		expected := []*models.SupplierCategoryRelation{
			{SupplierID: 1, CategoryID: categoryID},
		}

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(expected, nil)

		result, err := svc.GetByCategoryID(ctx, categoryID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		_, err := svc.GetByCategoryID(ctx, -1)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		categoryID := int64(2)

		mockRepo.On("GetByCategoryID", ctx, categoryID).Return(([]*models.SupplierCategoryRelation)(nil), errors.New("db error"))

		_, err := svc.GetByCategoryID(ctx, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_HasRelation(t *testing.T) {
	ctx := context.Background()

	t.Run("success - relation exists", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

		result, err := svc.HasRelation(ctx, supplierID, categoryID)

		assert.NoError(t, err)
		assert.True(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

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
				assert.ErrorIs(t, err, errMsg.ErrInvalidData)
			})
		}
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("HasRelation", ctx, supplierID, categoryID).Return(false, errors.New("db error"))

		_, err := svc.HasRelation(ctx, supplierID, categoryID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
