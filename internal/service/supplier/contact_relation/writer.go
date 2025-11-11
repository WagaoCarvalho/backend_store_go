package services

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

func (s *supplierContactRelationService) Create(ctx context.Context, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error) {
	if relation == nil {
		return nil, fmt.Errorf("%w: relação não fornecida", errMsg.ErrNilModel)
	}

	if relation.SupplierID <= 0 || relation.ContactID <= 0 {
		return nil, errMsg.ErrZeroID
	}

	createdRelation, err := s.relationRepo.Create(ctx, relation)
	if err != nil {
		switch {
		case errors.Is(err, errMsg.ErrRelationExists):
			// Recupera todas as relações e tenta retornar a já existente
			relations, getErr := s.relationRepo.GetAllRelationsBySupplierID(ctx, relation.SupplierID)
			if getErr != nil {
				return nil, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, getErr)
			}

			for _, rel := range relations {
				if rel.ContactID == relation.ContactID {
					return rel, nil
				}
			}

			return nil, errMsg.ErrRelationExists

		case errors.Is(err, errMsg.ErrDBInvalidForeignKey):
			return nil, errMsg.ErrDBInvalidForeignKey

		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return createdRelation, nil
}

func (s *supplierContactRelationService) Delete(ctx context.Context, supplierID, contactID int64) error {
	if supplierID <= 0 || contactID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.Delete(ctx, supplierID, contactID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (s *supplierContactRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	if supplierID <= 0 {
		return errMsg.ErrZeroID
	}

	err := s.relationRepo.DeleteAll(ctx, supplierID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
