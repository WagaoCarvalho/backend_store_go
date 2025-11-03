package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/category_relation"
)

type SupplierCategoryRelation interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierID(ctx context.Context, supplierID int64) error
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type supplierCategoryRelation struct {
	relationRepo repo.SupplierCategoryRelation
}

func NewSupplierCategoryRelation(repository repo.SupplierCategoryRelation) SupplierCategoryRelation {
	return &supplierCategoryRelation{relationRepo: repository}
}

func (s *supplierCategoryRelation) Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
	if relation == nil {
		return nil, errMsg.ErrNilModel
	}

	if relation.SupplierID <= 0 || relation.CategoryID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	exists, err := s.relationRepo.HasRelation(ctx, relation.SupplierID, relation.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	if exists {
		return nil, errMsg.ErrRelationExists
	}

	created, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return created, nil
}

func (s *supplierCategoryRelation) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error) {
	if supplierID <= 0 {
		return nil, errMsg.ErrInvalidData
	}

	result, err := s.relationRepo.GetBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return result, nil
}

func (s *supplierCategoryRelation) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error) {
	if categoryID <= 0 {
		return nil, errMsg.ErrInvalidData
	}

	result, err := s.relationRepo.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return result, nil
}

func (s *supplierCategoryRelation) Delete(ctx context.Context, supplierID, categoryID int64) error {
	if supplierID <= 0 || categoryID <= 0 {
		return errMsg.ErrInvalidData
	}

	err := s.relationRepo.Delete(ctx, supplierID, categoryID)
	if err != nil {
		if errors.Is(err, errMsg.ErrNotFound) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *supplierCategoryRelation) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return errMsg.ErrInvalidData
	}

	err := s.relationRepo.DeleteAllBySupplierID(ctx, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *supplierCategoryRelation) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	if supplierID <= 0 || categoryID <= 0 {
		return false, errMsg.ErrInvalidData
	}

	exists, err := s.relationRepo.HasRelation(ctx, supplierID, categoryID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return exists, nil
}
