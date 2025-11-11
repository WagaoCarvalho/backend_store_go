package services

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockSupplierCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func TestSupplierCategoryRelationService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 2}
		expected := &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 2}

		mockRepo.On("HasRelation", ctx, relation.SupplierID, relation.CategoryID).Return(false, nil)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(r *models.SupplierCategoryRelation) bool {
			return r.SupplierID == relation.SupplierID && r.CategoryID == relation.CategoryID
		})).Return(expected, nil)

		result, err := svc.Create(ctx, relation)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid data", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		tests := []struct {
			name     string
			relation *models.SupplierCategoryRelation
		}{
			{"zero supplierID", &models.SupplierCategoryRelation{SupplierID: 0, CategoryID: 1}},
			{"zero categoryID", &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 0}},
			{"negative supplierID", &models.SupplierCategoryRelation{SupplierID: -1, CategoryID: 2}},
			{"negative categoryID", &models.SupplierCategoryRelation{SupplierID: 3, CategoryID: -5}},
			{"both zero", &models.SupplierCategoryRelation{SupplierID: 0, CategoryID: 0}},
			{"nil relation", nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := svc.Create(ctx, tt.relation)
				assert.Nil(t, result)
				if tt.relation == nil {
					assert.ErrorIs(t, err, errMsg.ErrNilModel)
				} else {
					assert.ErrorIs(t, err, errMsg.ErrZeroID)
				}
			})
		}
	})

	t.Run("relation exists", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 2}

		mockRepo.On("HasRelation", ctx, relation.SupplierID, relation.CategoryID).Return(true, nil)

		result, err := svc.Create(ctx, relation)
		assert.ErrorIs(t, err, errMsg.ErrRelationExists)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to check relation", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 2}

		mockRepo.On("HasRelation", ctx, relation.SupplierID, relation.CategoryID).Return(false, errors.New("db error"))

		result, err := svc.Create(ctx, relation)
		assert.ErrorIs(t, err, errMsg.ErrRelationCheck)
		assert.ErrorContains(t, err, "db error")
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("failed to create", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		relation := &models.SupplierCategoryRelation{SupplierID: 1, CategoryID: 2}

		mockRepo.On("HasRelation", ctx, relation.SupplierID, relation.CategoryID).Return(false, nil)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(r *models.SupplierCategoryRelation) bool {
			return r.SupplierID == relation.SupplierID && r.CategoryID == relation.CategoryID
		})).Return(nil, errors.New("db error"))

		result, err := svc.Create(ctx, relation)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(nil)

		err := svc.Delete(ctx, supplierID, categoryID)
		assert.NoError(t, err)
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
				err := svc.Delete(ctx, tt.supplierID, tt.categoryID)
				assert.ErrorIs(t, err, errMsg.ErrInvalidData)
			})
		}
	})

	t.Run("relation not found", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(errMsg.ErrNotFound)

		err := svc.Delete(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)
		categoryID := int64(2)

		mockRepo.On("Delete", ctx, supplierID, categoryID).Return(errors.New("db error"))

		err := svc.Delete(ctx, supplierID, categoryID)
		assert.ErrorIs(t, err, errMsg.ErrDelete)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationService_DeleteAllBySupplierID(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(1)

		mockRepo.On("DeleteAllBySupplierID", ctx, supplierID).Return(nil)

		err := svc.DeleteAllBySupplierID(ctx, supplierID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		err := svc.DeleteAllBySupplierID(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockSupplierCatRel.MockSupplierCategoryRelation)
		svc := NewSupplierCategoryRelationService(mockRepo)

		supplierID := int64(99)

		mockRepo.On("DeleteAllBySupplierID", ctx, supplierID).Return(errors.New("db error"))

		err := svc.DeleteAllBySupplierID(ctx, supplierID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
