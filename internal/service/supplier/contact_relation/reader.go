package services

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierContactRelationService) GetAllRelationsBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierContactRelation, error) {
	if supplierID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	relations, err := s.relationRepo.GetAllRelationsBySupplierID(ctx, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return relations, nil
}

func (s *supplierContactRelationService) HasSupplierContactRelation(ctx context.Context, supplierID, contactID int64) (bool, error) {
	if supplierID <= 0 || contactID <= 0 {
		return false, errMsg.ErrZeroID
	}

	exists, err := s.relationRepo.HasSupplierContactRelation(ctx, supplierID, contactID)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return exists, nil
}
