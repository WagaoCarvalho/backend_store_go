package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierCategoryRelationService) Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error) {
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

func (s *supplierCategoryRelationService) Delete(ctx context.Context, supplierID, categoryID int64) error {
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

func (s *supplierCategoryRelationService) DeleteAllBySupplierID(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return errMsg.ErrInvalidData
	}

	err := s.relationRepo.DeleteAllBySupplierID(ctx, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
